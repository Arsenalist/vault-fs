---
name: vault-fs
description: >
  Use this skill whenever the user wants Claude to interact with a markdown vault
  — creating or reading notes, searching content, listing or extracting tasks,
  managing tags, setting frontmatter properties, or organizing files. Without this
  skill, Claude has no way to perform vault operations. Treat any request that
  implies "create a note", "search my vault", "list my tasks", "find recent files",
  or "set properties on a file" as a trigger. Also trigger for vault initialization,
  file organization, and automation workflows involving markdown vaults. Skip for
  pure conceptual questions about markdown syntax or general file management that
  doesn't involve a vault.
---

# vault-fs

A platform-agnostic CLI for managing markdown-based knowledge vaults. Designed for AI agent consumption with JSON-first output.

> Read `references/command-reference.md` when you need specific flags, output formats,
> or details for any command.

## Prerequisites

| Requirement | Details |
|---|---|
| Go binary | `vault-fs` must be installed and in PATH |
| Vault initialized | Run `vault-fs init` first, or it falls back to `~/vault-fs` |

## Syntax

```bash
vault-fs <command> [subcommand] [--flags]
```

### Global Flags

| Flag | Description |
|---|---|
| `--vault <path>` | Override vault path (takes precedence over all) |
| `--format <json\|text>` | Override output format |

### Vault Discovery

When `--vault` is not specified, vault-fs resolves the vault by:
1. `VAULTFS_PATH` environment variable
2. Walk up from CWD looking for `.vaultfs/` directory
3. Fall back to `~/vault-fs`

## Output Format

- **Query commands** (read, list, folders, recent, search, tags, tasks, properties, outline, info) → **JSON by default**
- **Action commands** (init, create, append, prepend, move, delete, mkdir, property set/remove, task toggle, index rebuild) → **Text by default**
- Use `--format=json` or `--format=text` to override

## Command Overview

| Group | Key Commands | Purpose |
|---|---|---|
| **vault** | `init`, `info` | Initialize vault, show metadata |
| **files** | `create`, `read`, `append`, `prepend`, `move`, `delete`, `list`, `folders`, `mkdir` | File CRUD and discovery |
| **recent** | `recent` | Recently modified files |
| **search** | `search`, `search:context`, `index rebuild` | Full-text and fuzzy search |
| **tags** | `tags`, `tag` | Tag listing and filtering |
| **tasks** | `tasks`, `task toggle` | Task extraction and toggling |
| **properties** | `properties`, `property set`, `property remove` | Frontmatter management |
| **outline** | `outline` | Heading structure |
| **help** | `help` | Comprehensive usage guide |

## Quick Reference

### Vault Setup

```bash
# Initialize with default preset
vault-fs init --preset=basic

# Initialize at custom path with extra dirs
vault-fs init --path=/tmp/my-vault --preset=basic --dirs="clients/acme,labs"

# List available presets
vault-fs init --list-presets

# Vault info
vault-fs info
```

### Creating & Reading Files

```bash
# Create (auto-adds .md, auto-creates parent dirs)
vault-fs create notes/standup --content="# Standup\n\nToday's notes"

# Create from a file (useful for large content)
vault-fs create notes/imported --input=/path/to/source.md

# Pipe content from stdin
cat draft.md | vault-fs create notes/draft --input=-

# Create with append mode (appends if file exists)
vault-fs create notes/standup --content="\n- New item" --append

# Read (returns JSON with frontmatter + body separated)
vault-fs read notes/standup.md

# Append (creates file if missing)
vault-fs append notes/standup.md --content="\n## Update\nNew section"

# Append from a file
vault-fs append notes/standup.md --input=/path/to/extra-content.md

# Prepend (inserts after frontmatter)
vault-fs prepend notes/standup.md --content="**Priority: High**\n\n"

# Prepend from a file
vault-fs prepend notes/standup.md --input=/path/to/header.md

# Move / rename
vault-fs move notes/old.md --to=archive/old.md

# Delete
vault-fs delete notes/old.md
```

### File Discovery

```bash
# List all markdown files
vault-fs list

# Filter by folder
vault-fs list --folder=projects

# Filter by extension
vault-fs list --ext=txt

# List folders
vault-fs folders

# Create directory
vault-fs mkdir projects/2026/q2
```

### Recent Files

```bash
# Last 7 days, max 20
vault-fs recent

# Custom window
vault-fs recent --days=30 --limit=50 --folder=projects
```

### Search

```bash
# Full-text search (AND semantics: all terms required, lazy-rebuilds index)
vault-fs search "quarterly review"

# Exact phrase match
vault-fs search "quarterly review" --exact

# Scoped to folder
vault-fs search "budget" --folder=projects --limit=5

# Fuzzy filename matching
vault-fs search "standup" --fuzzy

# Search with matching line context
vault-fs search:context "TODO"

# Force rebuild index
vault-fs index rebuild
```

### Tags

```bash
# All tags
vault-fs tags

# Tags with counts, sorted
vault-fs tags --counts --sort=count

# Files with specific tag
vault-fs tag project
vault-fs tag "project/alpha"
```

### Tasks

```bash
# All tasks (returns JSON with priority, due, tags, mentions)
vault-fs tasks

# Filter
vault-fs tasks --pending
vault-fs tasks --done
vault-fs tasks --folder=projects

# Toggle checkbox
vault-fs task toggle notes/todo.md --line=5
```

### Frontmatter Properties

```bash
# Read all properties
vault-fs properties notes/standup.md

# Set property (creates frontmatter if missing)
vault-fs property set notes/standup.md --name=status --value=active

# Remove property
vault-fs property remove notes/standup.md --name=draft
```

### Outline

```bash
# Heading tree
vault-fs outline notes/design.md
```

## Common Agent Patterns

### Daily Journal Entry

```bash
DATE=$(date +%Y-%m-%d)
vault-fs create "Daily Plan/$DATE" --content="---\ndate: $DATE\n---\n\n## Plan\n\n- [ ] \n\n## Notes\n"
```

### Create Note with Properties

```bash
vault-fs create projects/new-feature --content="# New Feature\n\nDescription here."
vault-fs property set projects/new-feature.md --name=status --value=planning
vault-fs property set projects/new-feature.md --name=priority --value=high
```

### Task Dashboard

```bash
# Get all pending tasks as JSON
vault-fs tasks --pending

# Get tasks from a specific area
vault-fs tasks --pending --folder=projects
```

### Search and Read Workflow

```bash
# Find relevant files
vault-fs search "authentication" --limit=5

# Read the top result
vault-fs read path/to/result.md
```

### Vault Analytics

```bash
vault-fs info                          # File/folder counts
vault-fs tags --counts --sort=count    # Most used tags
vault-fs tasks --pending               # Open tasks
vault-fs recent --days=1               # Today's activity
```

## Tips

1. **Paths are vault-relative** — use `notes/standup.md`, not absolute paths.
2. **`create` auto-adds `.md`** — pass `notes/standup`, not `notes/standup.md`.
3. **`append` creates files** — safe to append to nonexistent files (creates with parent dirs).
4. **`prepend` respects frontmatter** — content goes after `---` block, not at byte 0.
5. **Search uses AND** — `"abc def"` requires both terms. Use `--exact` for phrase matching.
6. **Search is lazy** — first search after file changes triggers auto-reindex.
7. **Tags come from two sources** — frontmatter `tags:` array AND inline `#tag` in body.
8. **Tasks parse metadata** — priority emojis, `#due/YYYY-MM-DD`, `@mentions`, `#tags`.
9. **Config is per-vault** — `.vaultfs/config.yaml` stores presets and settings.
10. **Use `--format=text`** when you want human-readable output from query commands.
11. **Directories with spaces work** — the basic preset includes "Daily Debrief", "Meeting Notes", etc.
12. **Run `vault-fs help`** for comprehensive built-in usage guide covering all commands and flags.
13. **Use `--input` for large content** — `create`, `append`, and `prepend` accept `--input=<file>` to read content from a file instead of `--content`. Use `--input=-` to read from stdin.
