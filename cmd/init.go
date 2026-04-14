package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zarar/vaultfs/internal/vault"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new vault",
	Long: `Initialize a new vault at the specified path (default: ~/vault-fs).

Use --preset to scaffold with a predefined directory structure.
Use --dirs to create additional custom directories (comma-separated, supports nesting like "projects/active").
Use --list-presets to see available presets.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		listPresetsFlag, _ := cmd.Flags().GetBool("list-presets")
		if listPresetsFlag {
			result, err := listPresets()
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		}

		pathFlag, _ := cmd.Flags().GetString("path")
		presetFlag, _ := cmd.Flags().GetString("preset")
		dirsFlag, _ := cmd.Flags().GetString("dirs")

		vaultPath := pathFlag
		if vaultPath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			vaultPath = filepath.Join(home, "vault-fs")
		}

		var dirs []string
		if dirsFlag != "" {
			dirs = strings.Split(dirsFlag, ",")
			for i := range dirs {
				dirs[i] = strings.TrimSpace(dirs[i])
			}
		}

		if err := runInit(vaultPath, presetFlag, dirs); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	initCmd.Flags().String("path", "", "Path for the new vault (default: ~/vault-fs)")
	initCmd.Flags().String("preset", "", "Preset to use for directory scaffolding (e.g., basic)")
	initCmd.Flags().String("dirs", "", "Comma-separated list of directories to create (supports nesting)")
	initCmd.Flags().Bool("list-presets", false, "List available presets as JSON")
	rootCmd.AddCommand(initCmd)
}

// runInit creates a new vault at the given path.
func runInit(vaultPath string, preset string, extraDirs []string) error {
	// Check if vault already exists
	if _, err := os.Stat(filepath.Join(vaultPath, ".vaultfs")); err == nil {
		fmt.Printf("Vault already exists at %s\n", vaultPath)
		return nil
	}

	// Create .vaultfs directory
	configDir := filepath.Join(vaultPath, ".vaultfs")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Load default config to get presets
	cfg, err := vault.LoadConfig(vaultPath)
	if err != nil {
		return fmt.Errorf("failed to load default config: %w", err)
	}

	// Collect directories to create
	var allDirs []string

	if preset != "" {
		p, ok := cfg.Presets[preset]
		if !ok {
			// Clean up the .vaultfs dir we created
			os.RemoveAll(configDir)
			return fmt.Errorf("unknown preset: %s", preset)
		}
		allDirs = append(allDirs, p.Directories...)
	}

	// Add extra dirs, deduplicating
	seen := make(map[string]bool)
	for _, d := range allDirs {
		seen[d] = true
	}
	for _, d := range extraDirs {
		if !seen[d] {
			allDirs = append(allDirs, d)
			seen[d] = true
		}
	}

	// Create all directories
	var createdDirs []string
	for _, d := range allDirs {
		fullPath := filepath.Join(vaultPath, d)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", d, err)
		}
		createdDirs = append(createdDirs, d)
	}

	// Write config.yaml
	configData, err := vault.MarshalConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), configData, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Write README.md
	readme := fmt.Sprintf("# Vault\n\nManaged by [vault-fs](https://github.com/zarar/vaultfs).\n")
	if err := os.WriteFile(filepath.Join(vaultPath, "README.md"), []byte(readme), 0644); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	// Output
	presetLabel := "none"
	if preset != "" {
		presetLabel = preset
	}
	fmt.Printf("Vault initialized at %s (preset: %s, %d directories created)\n", vaultPath, presetLabel, len(createdDirs))

	return nil
}

// listPresets returns available presets as JSON.
func listPresets() (string, error) {
	cfg, err := vault.LoadDefaultConfig()
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(cfg.Presets, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
