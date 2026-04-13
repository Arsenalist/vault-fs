package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zarar/vaultfs/internal/output"
)

// VaultInfo holds vault metadata.
type VaultInfo struct {
	Path        string `json:"path"`
	FileCount   int    `json:"file_count"`
	FolderCount int    `json:"folder_count"`
	ConfigPath  string `json:"config_path"`
	IndexExists bool   `json:"index_exists"`
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display vault metadata",
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		info, err := getVaultInfo(vaultPath)
		if err != nil {
			return err
		}

		format := output.ResolveFormat(formatFlag, true)
		if format == output.FormatJSON {
			return output.WriteJSON(os.Stdout, info)
		}
		fmt.Printf("Vault: %s\nFiles: %d\nFolders: %d\n", info.Path, info.FileCount, info.FolderCount)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func getVaultInfo(vaultPath string) (*VaultInfo, error) {
	var fileCount, folderCount int

	err := filepath.WalkDir(vaultPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip .vaultfs directory
		rel, _ := filepath.Rel(vaultPath, path)
		if rel == ".vaultfs" || (len(rel) > 9 && rel[:9] == ".vaultfs"+string(filepath.Separator)) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			if path != vaultPath {
				folderCount++
			}
		} else {
			fileCount++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(vaultPath, ".vaultfs", "config.yaml")
	_, indexErr := os.Stat(filepath.Join(vaultPath, ".vaultfs", "index.bleve"))

	return &VaultInfo{
		Path:        vaultPath,
		FileCount:   fileCount,
		FolderCount: folderCount,
		ConfigPath:  configPath,
		IndexExists: indexErr == nil,
	}, nil
}
