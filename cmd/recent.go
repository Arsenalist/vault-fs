package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zarar/vaultfs/internal/output"
)

func runRecent(vaultPath string, days, limit int, folder string) ([]FileInfo, error) {
	cutoff := time.Now().Add(-time.Duration(days) * 24 * time.Hour)

	searchRoot := vaultPath
	if folder != "" {
		searchRoot = filepath.Join(vaultPath, folder)
	}

	var files []FileInfo

	err := filepath.WalkDir(searchRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(vaultPath, path)

		if strings.HasPrefix(rel, ".vaultfs") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if info.ModTime().Before(cutoff) {
			return nil
		}

		files = append(files, FileInfo{
			Path:     filepath.ToSlash(rel),
			Size:     info.Size(),
			Modified: info.ModTime().Format(time.RFC3339),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort by mtime descending
	sort.Slice(files, func(i, j int) bool {
		return files[i].Modified > files[j].Modified
	})

	if limit > 0 && len(files) > limit {
		files = files[:limit]
	}

	return files, nil
}

var recentCmd = &cobra.Command{
	Use:   "recent",
	Short: "List recently modified files",
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		days, _ := cmd.Flags().GetInt("days")
		limit, _ := cmd.Flags().GetInt("limit")
		folder, _ := cmd.Flags().GetString("folder")

		files, err := runRecent(vaultPath, days, limit, folder)
		if err != nil {
			return err
		}

		format := output.ResolveFormat(formatFlag, true)
		if format == output.FormatJSON {
			return output.WriteJSON(os.Stdout, files)
		}
		for _, f := range files {
			fmt.Printf("%s  %s\n", f.Modified, f.Path)
		}
		return nil
	},
}

func init() {
	recentCmd.Flags().Int("days", 7, "Time window in days")
	recentCmd.Flags().Int("limit", 20, "Maximum results")
	recentCmd.Flags().String("folder", "", "Filter by folder")
	rootCmd.AddCommand(recentCmd)
}
