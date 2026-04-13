# vaultfs — Full Command Reference

Complete reference for all vaultfs CLI commands.

**Syntax**: `vaultfs <command> [subcommand] [--flags]`

**Global Flags**:
| Flag | Type | Default | Description |
|---|---|---|---|
| `--vault` | string | (auto-detected) | Override vault path |
| `--format` | string | (depends on command) | Output format: `json` or `text` |

---

## Table of Contents

1. [Vault Management](#vault-management)
2. [File Operations](#file-operations)
3. [Recent Files](#recent-files)
4. [Search](#search)
5. [Tags](#tags)
6. [Tasks](#tasks)
7. [Properties](#properties)
8. [Outline](#outline)
9. [Help](#help)

---

## Vault Management

### `vaultfs init`

Initialize a new vault.

| Flag | Type | Default | Description |
|---|---|---|---|
| `--path` | string | `~/vault-fs` | Vault location |
| `--preset` | string | (none) | Preset name (e.g., `basic`) |
| `--dirs` | string | (none) | Comma-separated directories to create |
| `--list-presets` | bool | false | List available presets as JSON |

**Output**: Text (action command)

```bash
vaultfs init --path=/tmp/vault --preset=basic
vaultfs init --preset=basic --dirs="clients/acme,labs"
vaultfs init --list-presets
```

**Default presets**:
- `basic`: Daily Debrief, Daily Plan, Journal, Meeting Notes, Projects/Active, Projects/Archived, Reports, Scratchpad, Stakeholders

### `vaultfs info`

Display vault metadata.

**Parameters**: None

**Output**: JSON (query command)

```json
{
  "path": "/Users/user/vault-fs",
  "file_count": 42,
  "folder_count": 12,
  "config_path": "/Users/user/vault-fs/.vaultfs/config.yaml",
  "index_exists": true
}
```

---

## File Operations

### `vaultfs create <path>`

Create a new markdown file.

| Flag | Type | Default | Description |
|---|---|---|---|
| `--content` | string | `""` | File content |
| `--append` | bool | false | Append to existing file instead of erroring |

**Output**: Text (action command)

- Path should **not** include `.md` — added automatically
- Parent directories are created automatically
- Errors if file exists (unless `--append` is set)

### `vaultfs read <path>`

Read a file with parsed frontmatter.

**Output**: JSON (query command)

```json
{
  "path": "notes/standup.md",
  "properties": {"title": "Standup", "tags": ["work"]},
  "body": "# Standup\n\nContent here.",
  "modified": "2026-04-13T09:30:00Z",
  "size": 284
}
```

### `vaultfs append <path>`

Append content to a file. Creates file (and parent dirs) if missing.

| Flag | Type | Default | Description |
|---|---|---|---|
| `--content` | string | `""` | Content to append |

**Output**: Text (action command)

### `vaultfs prepend <path>`

Prepend content after frontmatter (not at byte 0).

| Flag | Type | Default | Description |
|---|---|---|---|
| `--content` | string | `""` | Content to prepend |

**Output**: Text (action command)

### `vaultfs move <path>`

Move or rename a file. Creates target directories automatically.

| Flag | Type | Default | Description |
|---|---|---|---|
| `--to` | string | (required) | Target path |

**Output**: Text (action command)

### `vaultfs delete <path>`

Permanently delete a file.

**Output**: Text (action command)

### `vaultfs list`

List files in the vault.

| Flag | Type | Default | Description |
|---|---|---|---|
| `--folder` | string | (all) | Filter to folder |
| `--ext` | string | `md` | Filter by extension |

**Output**: JSON array (query command)

```json
[
  {"path": "notes/standup.md", "size": 284, "modified": "2026-04-13T09:30:00Z"},
  {"path": "Journal/day1.md", "size": 120, "modified": "2026-04-12T15:00:00Z"}
]
```

### `vaultfs folders`

List all directories (excludes `.vaultfs/`).

**Output**: JSON array (query command)

### `vaultfs mkdir <path>`

Create directories recursively.

**Output**: Text (action command)

---

## Recent Files

### `vaultfs recent`

List recently modified files, sorted by modification time (newest first).

| Flag | Type | Default | Description |
|---|---|---|---|
| `--days` | int | `7` | Time window in days |
| `--limit` | int | `20` | Maximum results |
| `--folder` | string | (all) | Filter to folder |

**Output**: JSON array (query command)

Same structure as `vaultfs list` output.

---

## Search

### `vaultfs search <query>`

Full-text search using bleve index. Index is lazily rebuilt when stale.

Uses AND semantics by default: `"abc def"` requires both terms to be present. Use `--exact` for phrase matching where `"abc def"` must appear as that exact sequence.

| Flag | Type | Default | Description |
|---|---|---|---|
| `--folder` | string | (all) | Filter to folder |
| `--limit` | int | `10` | Maximum results |
| `--fuzzy` | bool | false | Use fuzzy filename matching instead of content search |
| `--exact` | bool | false | Match exact phrase instead of AND-ing terms |

**Output**: JSON array (query command)

```json
[
  {"path": "notes/meeting.md", "score": 0.85},
  {"path": "projects/alpha.md", "score": 0.42}
]
```

### `vaultfs search:context <query>`

Search with matching line context.

| Flag | Type | Default | Description |
|---|---|---|---|
| `--limit` | int | `10` | Maximum file results |

**Output**: JSON array (query command)

```json
[
  {
    "path": "notes/meeting.md",
    "matches": [
      {"line": 3, "content": "Discussed quarterly budget and allocations."}
    ]
  }
]
```

### `vaultfs index rebuild`

Force rebuild the search index.

**Output**: Text (action command)

---

## Tags

### `vaultfs tags`

List all tags from frontmatter and inline `#tag` syntax.

| Flag | Type | Default | Description |
|---|---|---|---|
| `--counts` | bool | false | Include usage counts |
| `--sort` | string | `name` | Sort by: `name` or `count` |

**Output**: JSON array (query command)

```json
[
  {"name": "work", "count": 12},
  {"name": "urgent", "count": 5}
]
```

### `vaultfs tag <name>`

List files containing a specific tag. Supports nested tags (e.g., `project/alpha`).

**Output**: JSON array of file paths (query command)

---

## Tasks

### `vaultfs tasks`

Extract checkbox tasks from all vault files.

| Flag | Type | Default | Description |
|---|---|---|---|
| `--pending` | bool | false | Only incomplete tasks |
| `--done` | bool | false | Only completed tasks |
| `--folder` | string | (all) | Filter to folder |

**Output**: JSON array (query command)

```json
[
  {
    "file": "projects/launch.md",
    "line": 12,
    "text": "Fix critical bug",
    "done": false,
    "priority": "high",
    "due": "2026-04-15",
    "tags": ["backend"],
    "mentions": ["alice"]
  }
]
```

**Priority parsing**: `🔴`/`⏫` → high, `🟡`/`🔼` → medium, `🔵`/`🔽` → low
**Due date parsing**: `#due/YYYY-MM-DD` or `📅 YYYY-MM-DD`
**Mentions**: `@name` patterns
**Tags**: Inline `#tag` (excludes `#due/...`)

### `vaultfs task toggle <path>`

Toggle a task checkbox between `- [ ]` and `- [x]`.

| Flag | Type | Default | Description |
|---|---|---|---|
| `--line` | int | (required) | Line number of the task |

**Output**: Text (action command)

---

## Properties

### `vaultfs properties <path>`

Read all YAML frontmatter properties.

**Output**: JSON object (query command)

### `vaultfs property set <path>`

Set a frontmatter property. Creates frontmatter block if missing.

| Flag | Type | Default | Description |
|---|---|---|---|
| `--name` | string | (required) | Property name |
| `--value` | string | (required) | Property value |

**Output**: Text (action command)

### `vaultfs property remove <path>`

Remove a frontmatter property. Idempotent (no error if missing).

| Flag | Type | Default | Description |
|---|---|---|---|
| `--name` | string | (required) | Property name |

**Output**: Text (action command)

---

## Outline

### `vaultfs outline <path>`

Extract heading structure as a nested JSON tree.

**Output**: JSON array (query command)

```json
[
  {
    "level": 1,
    "text": "Title",
    "children": [
      {"level": 2, "text": "Section A", "children": [
        {"level": 3, "text": "Sub A1"}
      ]},
      {"level": 2, "text": "Section B"}
    ]
  }
]
```

---

## Help

### `vaultfs help`

Display comprehensive built-in usage guide covering all commands, flags, vault discovery, output format behavior, search semantics, task metadata parsing, and examples.

**Output**: Text (always)

```bash
vaultfs help
```
