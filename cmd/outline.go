package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zarar/vaultfs/internal/markdown"
	"github.com/zarar/vaultfs/internal/output"
	"github.com/zarar/vaultfs/internal/vfs"
)

func runOutline(vaultPath, path string) ([]*markdown.Heading, error) {
	data, err := os.ReadFile(filepath.Join(vaultPath, path))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, vfs.NewNotFound(path)
		}
		return nil, err
	}
	result := markdown.ExtractOutline(data)
	if result == nil {
		result = []*markdown.Heading{}
	}
	return result, nil
}

var outlineCmd = &cobra.Command{
	Use:   "outline <path>",
	Short: "Extract heading structure",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		outline, err := runOutline(vaultPath, args[0])
		if handled, e := handleNotFound(cmd, true, err); handled || e != nil {
			return e
		}
		format := output.ResolveFormat(formatFlag, true)
		if format == output.FormatJSON {
			return output.WriteJSON(os.Stdout, outline)
		}
		printOutline(outline, 0)
		return nil
	},
}

func printOutline(headings []*markdown.Heading, indent int) {
	for _, h := range headings {
		fmt.Printf("%s%s\n", strings.Repeat("  ", indent), h.Text)
		if len(h.Children) > 0 {
			printOutline(h.Children, indent+1)
		}
	}
}

func init() {
	rootCmd.AddCommand(outlineCmd)
}
