## Purpose

Defines requirements for the search capability of vaultfs.

## ADDED Requirements

### Requirement: Full-text search
The system SHALL provide full-text search across vault files using a bleve index, with results ranked by relevance.

#### Scenario: Basic search
- **WHEN** user runs `vaultfs search "meeting notes"`
- **THEN** the system outputs a JSON array of matching files with paths and relevance scores, ranked by score descending

#### Scenario: Search scoped to folder
- **WHEN** user runs `vaultfs search "standup" --folder=projects`
- **THEN** only files within `projects/` are searched

#### Scenario: Search with limit
- **WHEN** user runs `vaultfs search "design" --limit=5`
- **THEN** at most 5 results are returned

### Requirement: Search with context
The system SHALL return matching lines with surrounding context when using the context subcommand.

#### Scenario: Context search
- **WHEN** user runs `vaultfs search:context "TODO"`
- **THEN** the system outputs matching lines with file path, line number, the matching line, and 2 lines of surrounding context

### Requirement: Fuzzy filename search
The system SHALL support fuzzy matching against filenames when the `--fuzzy` flag is used.

#### Scenario: Fuzzy search
- **WHEN** user runs `vaultfs search "mtng nots" --fuzzy`
- **THEN** files with names approximately matching "mtng nots" (e.g., "meeting-notes.md") are returned with match scores

### Requirement: Lazy index rebuild
The system SHALL automatically rebuild the search index when it is stale (index mtime older than most recent file mtime).

#### Scenario: Stale index
- **WHEN** user runs a search and files have been modified since the last index build
- **THEN** the index is rebuilt incrementally before returning results

#### Scenario: No index exists
- **WHEN** user runs a search and no `.vaultfs/index.bleve` exists
- **THEN** the index is built from scratch before returning results

### Requirement: Explicit index rebuild
The system SHALL allow users to explicitly rebuild the search index.

#### Scenario: Force rebuild
- **WHEN** user runs `vaultfs index rebuild`
- **THEN** the search index is fully rebuilt from all vault files and a confirmation message is printed
