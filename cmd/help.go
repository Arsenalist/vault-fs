package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helpText = `vault-fs — A platform-agnostic markdown vault manager for AI agents

USAGE
  vault-fs <command> [subcommand] [--flags]

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
  vault-fs init [--path=<path>] [--preset=<name>] [--dirs=<a,b/c>]   Initialize a new vault
  vault-fs init --list-presets                                         List available presets
  vault-fs info                                                        Vault metadata (path, file/folder counts)

FILE OPERATIONS
  vault-fs create <path> --content="..."                  Create file (.md auto-added, parent dirs auto-created)
  vault-fs create <path> --content="..." --append         Append to existing file instead of erroring
  vault-fs read <path>                                    Read file (JSON: path, properties, body, modified, size)
  vault-fs append <path> --content="..."                  Append to file (creates if missing)
  vault-fs prepend <path> --content="..."                 Prepend after frontmatter
  vault-fs move <path> --to=<target>                      Move/rename (auto-creates target dirs)
  vault-fs delete <path>                                  Delete file
  vault-fs list [--folder=<dir>] [--ext=<ext>]            List files (default: .md)
  vault-fs folders                                        List all directories (excludes .vaultfs/)
  vault-fs mkdir <path>                                   Create directories recursively

RECENT FILES
  vault-fs recent [--days=7] [--limit=20] [--folder=<dir>]   Recently modified files (newest first)

SEARCH
  vault-fs search <query> [--folder=<dir>] [--limit=10]   Full-text search (AND semantics: all terms required)
  vault-fs search <query> --exact                          Exact phrase match
  vault-fs search <query> --fuzzy                          Fuzzy filename matching
  vault-fs search:context <query> [--limit=10]             Search with matching line context
  vault-fs index rebuild                                   Force rebuild search index

  Search uses a bleve index stored in .vaultfs/index.bleve.
  The index is lazily rebuilt when stale (files changed since last index).

TAGS
  vault-fs tags [--counts] [--sort=count]                 List all tags (frontmatter + inline #tag)
  vault-fs tag <name>                                     List files with a specific tag

  Tags are extracted from:
    - YAML frontmatter: tags: [work, urgent]
    - Inline syntax:    #tag, #tag/sub (nested)

TASKS
  vault-fs tasks [--pending] [--done] [--folder=<dir>]    List tasks with metadata
  vault-fs task toggle <path> --line=<N>                  Toggle checkbox on a specific line

  Task metadata parsed from:
    - Priority:  🔴/⏫ = high, 🟡/🔼 = medium, 🔵/🔽 = low
    - Due date:  #due/YYYY-MM-DD or 📅 YYYY-MM-DD
    - Mentions:  @name
    - Tags:      #tag (inline, excludes #due/...)

FRONTMATTER PROPERTIES
  vault-fs properties <path>                              Read all frontmatter as JSON
  vault-fs property set <path> --name=<key> --value=<val> Set property (creates frontmatter if missing)
  vault-fs property remove <path> --name=<key>            Remove property (idempotent)

OUTLINE
  vault-fs outline <path>                                 Heading structure as nested JSON tree

CONFIGURATION
  Per-vault config stored in .vaultfs/config.yaml.
  Presets define directory scaffolds for vault-fs init.

EXAMPLES
  vault-fs init --preset=basic
  vault-fs create "Daily Plan/2026-04-13" --content="## Priorities\n\n- [ ] "
  vault-fs tasks --pending --folder=projects
  vault-fs search "quarterly review" --limit=5
  vault-fs tags --counts --sort=count
  vault-fs property set notes/standup.md --name=status --value=active
`

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show comprehensive usage information",
	Long:  "Display detailed help for all vault-fs commands, flags, and concepts.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(helpText)
	},
}

func init() {
	rootCmd.AddCommand(helpCmd)
}
