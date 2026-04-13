## Context

This is a greenfield Go CLI tool (`vaultfs`) for managing markdown-based knowledge vaults from the command line. The primary consumers are AI agents that need structured (JSON) output, but humans can use it too with `--format=text`. There is no existing codebase — we're building from scratch in the current directory (`/Users/zarar/dev/vault-fs/`).

The tool is inspired by obsidian-cli but removes the Obsidian dependency entirely. It operates directly on the filesystem, making it platform-agnostic (macOS, Linux, Windows) and usable with any markdown vault regardless of editor.

## Goals / Non-Goals

**Goals:**
- Provide a complete CLI for vault CRUD, search, tag/task extraction, and frontmatter management
- JSON-default output for query commands so agents can parse results without extra tooling
- Config-driven vault initialization with presets stored in `.vaultfs/config.yaml`
- Lazy-rebuilt bleve search index for fast full-text search without manual index management
- Platform-agnostic filesystem operations (no Unix-specific paths or assumptions)
- Single static binary with no runtime dependencies

**Non-Goals:**
- Wikilink parsing or link graph analysis (future consideration)
- File watching or live index updates
- GUI, TUI, or interactive modes
- Obsidian plugin/theme/sync management
- Real-time collaboration or multi-user access
- Template variable substitution (just file-based templates)

## Decisions

### 1. CLI Framework: cobra

**Choice**: `spf13/cobra`
**Over**: `urfave/cli/v2`, hand-rolled flag parsing
**Why**: Industry standard (43.7k stars), excellent subcommand support, auto-generated help text, shell completion. Every major Go CLI uses it. Agents benefit from predictable `--help` output.

### 2. Markdown Parser: goldmark

**Choice**: `yuin/goldmark`
**Over**: `gomarkdown/markdown`, `russross/blackfriday`
**Why**: Actively maintained, CommonMark compliant, extensible AST. We need AST walking for heading extraction (outline), task checkbox detection, and inline tag parsing — goldmark's extension system handles all three.

### 3. YAML Frontmatter: goccy/go-yaml

**Choice**: `goccy/go-yaml`
**Over**: `go-yaml/yaml.v3` (archived April 2025), `go.yaml.in/yaml/v4`
**Why**: Actively maintained, API-compatible migration path, good performance. The original go-yaml/yaml is archived and unmaintained.

### 4. Search: bleve with lazy rebuild

**Choice**: `blevesearch/bleve` with on-demand index rebuilding
**Over**: Simple `strings.Contains` loops, sqlite FTS5, external search services
**Why**: Embedded (no external deps), supports full-text search with relevance scoring, handles thousands of files well. Lazy rebuild means: on search, check if index mtime < most recent file mtime, rebuild if stale. Also provide explicit `vaultfs index rebuild`.

### 5. Fuzzy Matching: sahilm/fuzzy

**Choice**: `sahilm/fuzzy` for filename fuzzy matching
**Over**: `lithammer/fuzzysearch` (less maintained), Levenshtein DIY
**Why**: Simple API, good enough for filename matching. Used alongside bleve (bleve for content, fuzzy for filenames). Library is stable/complete even if not actively developed — small surface area, unlikely to need updates.

### 6. Output Strategy: split by command type

**Choice**: JSON default for queries, text default for actions
**Over**: JSON-everything, text-everything
**Why**: Agents parse query results (search, list, tags, tasks, read, properties, recent, outline). Action commands (init, create, append, delete, move, mkdir, property set/remove, task toggle, index rebuild) return human-readable confirmations. All commands accept `--format=json` or `--format=text` to override.

### 7. Per-vault Config with Embedded Defaults

**Choice**: `.vaultfs/config.yaml` inside the vault, with defaults embedded in the binary via `embed`
**Over**: Global-only config, hardcoded presets
**Why**: Vaults are self-contained and portable. The binary embeds a default config (containing the `basic` preset) which gets written to `.vaultfs/config.yaml` on `init`. Users can edit it freely — add presets, change defaults. No global config needed initially.

### 8. Project Layout

```
vaultfs/
├── cmd/                    # cobra command definitions
│   ├── root.go             # root command, --format, --vault flags
│   ├── init.go
│   ├── info.go
│   ├── create.go
│   ├── read.go
│   ├── append.go
│   ├── prepend.go
│   ├── move.go
│   ├── delete.go
│   ├── list.go
│   ├── folders.go
│   ├── mkdir.go
│   ├── recent.go
│   ├── search.go
│   ├── tags.go
│   ├── tasks.go
│   ├── properties.go
│   ├── outline.go
│   └── index.go
├── internal/
│   ├── vault/              # vault discovery, config loading, path resolution
│   ├── markdown/           # goldmark parsing, frontmatter extraction, tag/task parsing
│   ├── search/             # bleve index management, fuzzy filename matching
│   └── output/             # JSON/text formatting, response structs
├── defaults/
│   └── config.yaml         # embedded default config (basic preset)
├── main.go
├── go.mod
└── go.sum
```

### 9. Vault Discovery

When `--vault` is not specified, `vaultfs` resolves the vault by:
1. Check `VAULTFS_PATH` environment variable
2. Walk up from CWD looking for `.vaultfs/` directory
3. Fall back to `~/vault-fs`

This lets agents set an env var once, or operate naturally inside a vault directory.

## Risks / Trade-offs

**[Bleve index size]** → Bleve creates non-trivial index files for large vaults. Mitigation: index lives in `.vaultfs/index.bleve`, can be deleted and rebuilt. Add to `.gitignore` by default.

**[Lazy rebuild performance]** → First search after many file changes triggers a full reindex, which could be slow on very large vaults (10k+ files). Mitigation: incremental reindex (only changed files based on mtime comparison). If still too slow, user can run `vaultfs index rebuild` proactively.

**[sahilm/fuzzy maintenance]** → Library hasn't had a release since 2024. Mitigation: it's 300 lines of code with zero dependencies — can vendor or fork if needed. Stable algorithm, unlikely to need updates.

**[Cross-platform path handling]** → Windows uses `\` separators, macOS/Linux use `/`. Mitigation: use `filepath.Join`, `filepath.Clean` everywhere. Store paths with `/` in config and output, convert at OS boundary.

**[Frontmatter parsing edge cases]** → YAML frontmatter delimited by `---` can conflict with markdown horizontal rules. Mitigation: only parse frontmatter at byte 0 of file, require opening and closing `---` before any non-YAML content.
