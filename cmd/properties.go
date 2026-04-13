package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/zarar/vaultfs/internal/markdown"
	"github.com/zarar/vaultfs/internal/output"
)

func runProperties(vaultPath, path string) (map[string]any, error) {
	data, err := os.ReadFile(filepath.Join(vaultPath, path))
	if err != nil {
		return nil, err
	}
	fm, _, err := markdown.ParseFrontmatter(data)
	return fm, err
}

func runPropertySet(vaultPath, path, name, value string) error {
	fullPath := filepath.Join(vaultPath, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	fm, body, err := markdown.ParseFrontmatter(data)
	if err != nil {
		return err
	}

	fm[name] = value

	return writeFrontmatterAndBody(fullPath, fm, body)
}

func runPropertyRemove(vaultPath, path, name string) error {
	fullPath := filepath.Join(vaultPath, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	fm, body, err := markdown.ParseFrontmatter(data)
	if err != nil {
		return err
	}

	delete(fm, name)

	return writeFrontmatterAndBody(fullPath, fm, body)
}

func writeFrontmatterAndBody(fullPath string, fm map[string]any, body []byte) error {
	var content string

	if len(fm) > 0 {
		fmBytes, err := yaml.Marshal(fm)
		if err != nil {
			return err
		}
		content = "---\n" + string(fmBytes) + "---\n"
		if len(body) > 0 {
			content += "\n" + string(body)
		}
	} else {
		content = string(body)
	}

	return os.WriteFile(fullPath, []byte(content), 0644)
}

var propertiesCmd = &cobra.Command{
	Use:   "properties <path>",
	Short: "Read all frontmatter properties",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		props, err := runProperties(vaultPath, args[0])
		if err != nil {
			return err
		}
		format := output.ResolveFormat(formatFlag, true)
		if format == output.FormatJSON {
			return output.WriteJSON(os.Stdout, props)
		}
		for k, v := range props {
			fmt.Printf("%s: %v\n", k, v)
		}
		return nil
	},
}

var propertyCmd = &cobra.Command{
	Use:   "property",
	Short: "Property operations",
}

var propertySetCmd = &cobra.Command{
	Use:   "set <path>",
	Short: "Set a frontmatter property",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		value, _ := cmd.Flags().GetString("value")
		if err := runPropertySet(vaultPath, args[0], name, value); err != nil {
			return err
		}
		fmt.Printf("Set %s=%s on %s\n", name, value, args[0])
		return nil
	},
}

var propertyRemoveCmd = &cobra.Command{
	Use:   "remove <path>",
	Short: "Remove a frontmatter property",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		if err := runPropertyRemove(vaultPath, args[0], name); err != nil {
			return err
		}
		fmt.Printf("Removed %s from %s\n", name, args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(propertiesCmd)

	propertySetCmd.Flags().String("name", "", "Property name")
	propertySetCmd.Flags().String("value", "", "Property value")
	propertyRemoveCmd.Flags().String("name", "", "Property name")

	propertyCmd.AddCommand(propertySetCmd)
	propertyCmd.AddCommand(propertyRemoveCmd)
	rootCmd.AddCommand(propertyCmd)
}
