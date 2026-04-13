package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zarar/vaultfs/internal/vault"
)

var (
	formatFlag string
	vaultFlag  string
)

var rootCmd = &cobra.Command{
	Use:   "vaultfs",
	Short: "A platform-agnostic markdown vault manager for AI agents",
	Long:  "vaultfs is a CLI tool for managing markdown-based knowledge vaults. It provides file CRUD, search, tag/task extraction, and frontmatter management with JSON-first output for agent consumption.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// resolveVault determines the vault path using the discovery chain.
func resolveVault() (string, error) {
	cwd, _ := os.Getwd()
	return vault.Discover(vaultFlag, cwd)
}

func init() {
	rootCmd.PersistentFlags().StringVar(&formatFlag, "format", "", "Output format: json or text (default depends on command type)")
	rootCmd.PersistentFlags().StringVar(&vaultFlag, "vault", "", "Path to vault (overrides VAULTFS_PATH and auto-detection)")
}
