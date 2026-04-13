package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helpText = `vaultfs — A platform-agnostic markdown vault manager for AI agents

USAGE
  vaultfs <command> [subcommand] [--flags]

GLOBAL FLAGS
  --vault <path>        Override vault path (takes precedence over env/auto-detection)
  --format <json|text>  Override output format

VAULT DISCOVERY
  1. --vault flag
  2. VAULTFS_PATH environment variable
  3. Walk up from CWD looking for .vaultfs/ directory
  4. Fall back to ~/vault-fs

OUTPUT FORMAT
  Query commands default to JSON:  read, list, folders, recent, search, tags, tasks, properties, outline, info
  Action commands default to text:  init, create, append, prepend, move, delete, mkdir, property set/remove, task toggle, index rebuild

VAULT MANAGEMENT
  vaultfs init [--path=<path>] [--preset=<name>] [--dirs=<a,b/c>]   Initialize a new vault
  vaultfs init --list-presets                                         List available presets
  vaultfs info                                                        Vault metadata (path, file/folder counts)

FILE OPERATIONS
  vaultfs create <path> --content="..."                  Create file (.md auto-added, parent dirs auto-created)
  vaultfs create <path> --content="..." --append         Append to existing file instead of erroring
  vaultfs read <path>                                    Read file (JSON: path, properties, body, modified, size)
  vaultfs append <path> --content="..."                  Append to file (creates if missing)
  vaultfs prepend <path> --content="..."                 Prepend after frontmatter
  vaultfs move <path> --to=<target>                      Move/rename (auto-creates target dirs)
  vaultfs delete <path>                                  Delete file
  vaultfs list [--folder=<dir>] [--ext=<ext>]            List files (default: .md)
  vaultfs folders                                        List all directories (excludes .vaultfs/)
  vaultfs mkdir <path>                                   Create directories recursively

RECENT FILES
  vaultfs recent [--days=7] [--limit=20] [--folder=<dir>]   Recently modified files (newest first)

SEARCH
  vaultfs search <query> [--folder=<dir>] [--limit=10]   Full-text search (AND semantics: all terms required)
  vaultfs search <query> --exact                          Exact phrase match
  vaultfs search <query> --fuzzy                          Fuzzy filename matching
  vaultfs search:context <query> [--limit=10]             Search with matching line context
  vaultfs index rebuild                                   Force rebuild search index

  Search uses a bleve index stored in .vaultfs/index.bleve.
  The index is lazily rebuilt when stale (files changed since last index).

TAGS
  vaultfs tags [--counts] [--sort=count]                 List all tags (frontmatter + inline #tag)
  vaultfs tag <name>                                     List files with a specific tag

  Tags are extracted from:
    - YAML frontmatter: tags: [work, urgent]
    - Inline syntax:    #tag, #tag/sub (nested)

TASKS
  vaultfs tasks [--pending] [--done] [--folder=<dir>]    List tasks with metadata
  vaultfs task toggle <path> --line=<N>                  Toggle checkbox on a specific line

  Task metadata parsed from:
    - Priority:  🔴/⏫ = high, 🟡/🔼 = medium, 🔵/🔽 = low
    - Due date:  #due/YYYY-MM-DD or 📅 YYYY-MM-DD
    - Mentions:  @name
    - Tags:      #tag (inline, excludes #due/...)

FRONTMATTER PROPERTIES
  vaultfs properties <path>                              Read all frontmatter as JSON
  vaultfs property set <path> --name=<key> --value=<val> Set property (creates frontmatter if missing)
  vaultfs property remove <path> --name=<key>            Remove property (idempotent)

OUTLINE
  vaultfs outline <path>                                 Heading structure as nested JSON tree

CONFIGURATION
  Per-vault config stored in .vaultfs/config.yaml.
  Presets define directory scaffolds for vaultfs init.

EXAMPLES
  vaultfs init --preset=basic
  vaultfs create "Daily Plan/2026-04-13" --content="## Priorities\n\n- [ ] "
  vaultfs tasks --pending --folder=projects
  vaultfs search "quarterly review" --limit=5
  vaultfs tags --counts --sort=count
  vaultfs property set notes/standup.md --name=status --value=active
`

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show comprehensive usage information",
	Long:  "Display detailed help for all vaultfs commands, flags, and concepts.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(helpText)
	},
}

func init() {
	rootCmd.AddCommand(helpCmd)
}
