# vault-fs

A platform-agnostic CLI for managing markdown-based knowledge vaults. Built for AI agents, usable by humans.

## What is it?

vault-fs gives you a single binary to create, read, search, and organize markdown files in a structured vault — no Obsidian, no GUI, no runtime dependencies. Every query command returns JSON by default so agents can parse results without extra tooling.

## Install

```bash
go install github.com/zarar/vaultfs@latest
```

Or build from source:

```bash
git clone https://github.com/zarar/vaultfs.git
cd vaultfs
go build -o vault-fs .
```

## Quick Start

```bash
# Initialize a vault with the basic preset
vault-fs init --preset=basic

# Create a note (auto-adds .md, auto-creates parent dirs)
vault-fs create "Daily Plan/2026-04-13" --content="## Priorities\n\n- [ ] Ship feature\n- [ ] Review PRs"

# Create from a file (for large content)
vault-fs create "imports/report" --input=~/Downloads/report.md

# Pipe content from stdin
cat draft.md | vault-fs create "notes/draft" --input=-

# Read it back (JSON with parsed frontmatter)
vault-fs read "Daily Plan/2026-04-13.md"

# Search across all files (AND semantics)
vault-fs search "ship feature"

# List pending tasks
vault-fs tasks --pending

# See all tags with usage counts
vault-fs tags --counts --sort=count

# Get comprehensive help
vault-fs help
```

## Features

### Vault Management
- **`init`** — Scaffold a vault with config-driven presets (or custom directories)
- **`info`** — Vault metadata: file count, folder count, index status
- **Vault discovery** — `--vault` flag → `VAULTFS_PATH` env → walk up from CWD → `~/vault-fs`

### File Operations
- **`create`** — Create files with content, auto `.md` extension, auto parent dirs
- **`read`** — Returns parsed frontmatter + body as structured JSON
- **`append`** / **`prepend`** — Append creates missing files; prepend inserts after frontmatter
- **`--input`** — All write commands accept `--input=<file>` to read content from a file (use `-` for stdin)
- **`move`** / **`delete`** — Move with auto target dirs, permanent delete
- **`list`** / **`folders`** / **`mkdir`** — File discovery and directory management

### Search
- **`search`** — Full-text search via bleve index with AND semantics (`--exact` for phrase match)
- **`search:context`** — Returns matching lines with file path and line numbers
- **`--fuzzy`** — Fuzzy filename matching
- **Lazy indexing** — Index auto-rebuilds when stale; `index rebuild` for manual refresh

### Tags
- **`tags`** — List all tags with optional counts and sorting
- **`tag <name>`** — Find files with a specific tag (supports nested `#tag/sub`)
- Extracts from both YAML frontmatter `tags:` and inline `#tag` syntax

### Tasks
- **`tasks`** — Extract checkbox tasks with rich metadata parsing
- **`task toggle`** — Toggle `- [ ]` ↔ `- [x]` at a specific line
- Parses priority (`🔴`/`⏫` high, `🟡`/`🔼` medium, `🔵`/`🔽` low), due dates (`#due/YYYY-MM-DD`, `📅 YYYY-MM-DD`), `@mentions`, and inline `#tags`

### Frontmatter Properties
- **`properties`** — Read all YAML frontmatter as JSON
- **`property set`** — Set a property (creates frontmatter block if missing)
- **`property remove`** — Remove a property (idempotent)

### Outline
- **`outline`** — Heading structure as a nested JSON tree

## Output Format

| Command Type | Default | Override |
|---|---|---|
| Query commands (read, list, search, tags, tasks, ...) | JSON | `--format=text` |
| Action commands (init, create, append, delete, ...) | Text | `--format=json` |

## Configuration

Per-vault config lives in `.vaultfs/config.yaml`:

```yaml
vault:
  path: ~/vault-fs

presets:
  basic:
    directories:
      - Daily Debrief
      - Daily Plan
      - Journal
      - Meeting Notes
      - Projects/Active
      - Projects/Archived
      - Reports
      - Scratchpad
      - Stakeholders
```

Add your own presets and use them with `vault-fs init --preset=<name>`.

## Claude Code Skill

vault-fs ships with a [Claude Code skill](.claude/skills/vaultfs/SKILL.md) that teaches AI agents how to use every command. When working in a project with vault-fs, Claude automatically knows how to create notes, search content, extract tasks, manage tags, and more.

The skill includes:
- Full command overview with examples
- Common agent workflow patterns (daily journals, task dashboards, search-and-read)
- A [complete command reference](.claude/skills/vaultfs/references/command-reference.md) with all flags and output formats

## Tech Stack

| Component | Library |
|---|---|
| CLI framework | [cobra](https://github.com/spf13/cobra) |
| Markdown parsing | [goldmark](https://github.com/yuin/goldmark) |
| YAML frontmatter | [goccy/go-yaml](https://github.com/goccy/go-yaml) |
| Full-text search | [bleve](https://github.com/blevesearch/bleve) |
| Fuzzy matching | [sahilm/fuzzy](https://github.com/sahilm/fuzzy) |

## Development

```bash
# Run all tests (103 tests across 4 packages)
go test ./...

# Build
go build -o vault-fs .

# Vet
go vet ./...
```

## License

MIT
