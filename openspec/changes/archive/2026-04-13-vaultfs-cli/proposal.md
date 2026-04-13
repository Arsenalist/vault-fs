## Why

AI agents need a simple, platform-agnostic way to manage markdown-based knowledge vaults — creating notes, extracting tasks/tags, searching content — without depending on Obsidian or any specific editor. No tool exists today that is filesystem-native, CLI-first, and designed for agent consumption (JSON output, predictable structure, zero GUI dependency).

## What Changes

- New Go CLI binary `vaultfs` with cobra-based command structure
- Vault initialization with config-driven directory presets (e.g., `basic` preset with Daily Debrief, Journal, Meeting Notes, etc.)
- File CRUD operations (create, read, append, prepend, move, delete) with automatic parent directory creation
- Directory listing and creation commands
- Recent files query by modification timestamp
- Tag extraction from both YAML frontmatter and inline `#tag` syntax (including nested `#tag/sub`)
- Task extraction from markdown checkboxes (`- [ ]` / `- [x]`) with priority, due date, and mention parsing
- Full-text search via bleve with lazy index rebuilding, plus fuzzy filename matching
- YAML frontmatter parsing and property get/set/remove commands
- Heading outline extraction
- JSON-default output for query commands, human-readable for action commands, with `--format` override

## Capabilities

### New Capabilities
- `vault-management`: Vault init with config-driven presets, vault info, per-vault config in `.vaultfs/config.yaml`
- `file-operations`: Create/read/append/prepend/move/delete files with auto-dir creation; list files/folders; mkdir
- `recent-files`: Query recently modified files by mtime with configurable window and limits
- `search`: Full-text search via bleve index (lazy rebuild), fuzzy filename matching, context-aware search results
- `tags`: Extract and query tags from frontmatter `tags:` arrays and inline `#tag`/`#tag/sub` syntax across vault
- `tasks`: Extract markdown checkboxes with priority emoji parsing, due date extraction, @mentions, and toggle support
- `frontmatter`: Parse/read/set/remove YAML frontmatter properties on markdown files
- `outline`: Extract heading structure from markdown files as structured output
- `output-formatting`: JSON-default for queries, text-default for actions, with `--format` flag override

### Modified Capabilities

(none — greenfield project)

## Impact

- **New binary**: `vaultfs` Go CLI, installable via `go install`
- **Dependencies**: cobra, goldmark, goccy/go-yaml, sahilm/fuzzy, blevesearch/bleve
- **Filesystem**: Creates `.vaultfs/` directory inside each vault for config and search index
- **Platforms**: macOS, Linux, Windows — all filesystem operations must be platform-agnostic
