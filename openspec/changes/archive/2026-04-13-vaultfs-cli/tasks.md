## 1. Project Setup

- [x] 1.1 Initialize Go module (`go mod init`), create directory structure (`cmd/`, `internal/vault/`, `internal/markdown/`, `internal/search/`, `internal/output/`, `defaults/`), add `.gitignore` (with `index.bleve`, binary, etc.)
- [x] 1.2 Add dependencies: cobra, goldmark, goccy/go-yaml, sahilm/fuzzy, blevesearch/bleve
- [x] 1.3 Create `main.go` with cobra root command including `--format` and `--vault` global flags
- [x] 1.4 Create `defaults/config.yaml` with the `basic` preset (Daily Debrief, Daily Plan, Journal, Meeting Notes, Projects/Active, Projects/Archived, Reports, Scratchpad, Stakeholders) and embed it via `go:embed`

## 2. Output Formatting (TDD)

- [x] 2.1 Write tests for `internal/output/`: JSON formatter, text formatter, error formatting, query-vs-action default logic
- [x] 2.2 Implement `internal/output/` package to pass all tests

## 3. Vault Management (TDD)

- [x] 3.1 Write tests for `internal/vault/`: vault discovery chain (--vault → env → CWD walk-up → fallback), config loading, config merging with embedded defaults
- [x] 3.2 Implement `internal/vault/` package to pass all tests
- [x] 3.3 Write tests for `vaultfs init`: default path, custom path, preset dirs, custom --dirs, preset+dirs merge, --list-presets, error on existing vault
- [x] 3.4 Implement `vaultfs init` command to pass all tests
- [x] 3.5 Write tests for `vaultfs info`: vault path, file/folder counts, index status
- [x] 3.6 Implement `vaultfs info` command to pass all tests

## 4. Markdown Parsing Core (TDD)

- [x] 4.1 Write tests for frontmatter parser: split YAML from body, empty frontmatter, no frontmatter, malformed frontmatter
- [x] 4.2 Implement frontmatter parser in `internal/markdown/` to pass all tests
- [x] 4.3 Write tests for inline tag extraction: `#tag`, `#tag/sub`, tags in code blocks (should be ignored), deduplication with frontmatter tags
- [x] 4.4 Implement tag extraction to pass all tests
- [x] 4.5 Write tests for task extraction: `- [ ]`/`- [x]` with line numbers, priority emojis (🔴⏫→high, 🟡🔼→medium, 🔵🔽→low), due dates (`#due/YYYY-MM-DD`, `📅 YYYY-MM-DD`), @mentions, inline #tags, plain tasks with no metadata
- [x] 4.6 Implement task extraction to pass all tests
- [x] 4.7 Write tests for heading outline: nested heading tree, flat headings, no headings, skipped levels
- [x] 4.8 Implement heading outline extraction to pass all tests

## 5. File Operations (TDD)

- [x] 5.1 Write tests for `vaultfs create`: with content, auto `.md` extension, auto parent dirs, error on existing file, --append flag on existing file
- [x] 5.2 Implement `vaultfs create` to pass all tests
- [x] 5.3 Write tests for `vaultfs read`: file with frontmatter, without frontmatter, nonexistent file, JSON output structure (path, properties, body, modified, size)
- [x] 5.4 Implement `vaultfs read` to pass all tests
- [x] 5.5 Write tests for `vaultfs append`: append to existing, append creates nonexistent file with parent dirs; `vaultfs prepend`: prepend after frontmatter, prepend on file without frontmatter
- [x] 5.6 Implement `vaultfs append` and `vaultfs prepend` to pass all tests
- [x] 5.7 Write tests for `vaultfs move` (auto target dirs), `vaultfs delete`, `vaultfs list` (--folder, --ext filters), `vaultfs folders` (excludes .vaultfs/), `vaultfs mkdir` (recursive)
- [x] 5.8 Implement move, delete, list, folders, mkdir to pass all tests

## 6. Recent Files (TDD)

- [x] 6.1 Write tests for `vaultfs recent`: default 7 days/20 limit, custom --days/--limit, --folder filter, empty result, sorted by mtime descending
- [x] 6.2 Implement `vaultfs recent` to pass all tests

## 7. Tags (TDD)

- [x] 7.1 Write tests for `vaultfs tags`: all tags, --counts, --sort=count; `vaultfs tag <name>`: files with tag, nested tags
- [x] 7.2 Implement tags and tag commands to pass all tests

## 8. Tasks (TDD)

- [x] 8.1 Write tests for `vaultfs tasks`: all/--pending/--done/--folder filters, JSON structure with priority/due/tags/mentions; `vaultfs task toggle`: toggle [ ] to [x] and vice versa
- [x] 8.2 Implement tasks and task toggle commands to pass all tests

## 9. Frontmatter Properties (TDD)

- [x] 9.1 Write tests for `vaultfs properties`: read all as JSON, empty frontmatter; `vaultfs property set`: set on existing, create frontmatter if missing; `vaultfs property remove`: remove existing, idempotent on missing
- [x] 9.2 Implement properties, property set, property remove to pass all tests

## 10. Outline (TDD)

- [x] 10.1 Write tests for `vaultfs outline`: nested tree output, no headings, skipped levels
- [x] 10.2 Implement `vaultfs outline` to pass all tests

## 11. Search (TDD)

- [x] 11.1 Write tests for `internal/search/`: bleve index create/open/rebuild/close, document indexing, basic query, stale index detection
- [x] 11.2 Implement `internal/search/` package to pass all tests
- [x] 11.3 Write tests for `vaultfs search`: full-text query, --folder, --limit, relevance ranking; `vaultfs search:context`: matching lines with context
- [x] 11.4 Implement search and search:context commands to pass all tests
- [x] 11.5 Write tests for --fuzzy flag: fuzzy filename matching, score ordering
- [x] 11.6 Implement fuzzy search to pass all tests
- [x] 11.7 Write tests for `vaultfs index rebuild`: full reindex, confirmation output
- [x] 11.8 Implement index rebuild command to pass all tests

## 12. Agent Skill

- [x] 12.1 Create `skills/vaultfs/SKILL.md` with skill frontmatter (name, description, trigger conditions), command overview table, quick reference with examples for all command groups, output format details, vault discovery docs, common agent patterns, and tips
- [x] 12.2 Create `skills/vaultfs/references/command-reference.md` with full parameter tables for every command and subcommand

## 13. Final Integration Tests & Polish

- [x] 13.1 Write end-to-end integration tests: full workflow (init → create → read → search → tags → tasks → properties → outline)
- [x] 13.2 Verify cross-platform path handling (filepath.Join/Clean everywhere, `/` in output)
- [x] 13.3 Run `go vet`, `staticcheck`, ensure clean `go build` on all targets
