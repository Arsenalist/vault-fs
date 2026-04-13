## Purpose

Defines requirements for the file-operations capability of vaultfs.

## ADDED Requirements

### Requirement: Create a file
The system SHALL create a markdown file at the specified vault-relative path, automatically creating any parent directories that do not exist. The `.md` extension SHALL be appended automatically if omitted.

#### Scenario: Create with content
- **WHEN** user runs `vaultfs create notes/standup --content="# Standup

Today's notes"`
- **THEN** `notes/standup.md` is created with the specified content, and the `notes/` directory is created if it doesn't exist

#### Scenario: Create existing file defaults to error
- **WHEN** user runs `vaultfs create` for a path that already exists and no `--append` flag is set
- **THEN** the system SHALL report an error and not overwrite the existing file

#### Scenario: Create with append on existing file
- **WHEN** user runs `vaultfs create notes/standup --content="New content" --append`
- **THEN** the content is appended to the end of the existing file

### Requirement: Read a file
The system SHALL read a markdown file and return its frontmatter (parsed as structured data) and body separately.

#### Scenario: Read file with frontmatter
- **WHEN** user runs `vaultfs read notes/standup.md`
- **THEN** the system outputs JSON with `path`, `properties` (parsed frontmatter), `body` (markdown content), `modified` (ISO 8601 timestamp), and `size` (bytes)

#### Scenario: Read file without frontmatter
- **WHEN** user runs `vaultfs read` on a file with no YAML frontmatter
- **THEN** the `properties` field is an empty object and `body` contains the full file content

#### Scenario: Read nonexistent file
- **WHEN** user runs `vaultfs read` for a path that does not exist
- **THEN** the system SHALL report an error with exit code 1

### Requirement: Append content to a file
The system SHALL append content to the end of a markdown file, creating the file (and parent directories) if it does not exist.

#### Scenario: Append to existing file
- **WHEN** user runs `vaultfs append notes/standup.md --content="- New item"`
- **THEN** the content is added at the end of the file

#### Scenario: Append to nonexistent file
- **WHEN** user runs `vaultfs append notes/new.md --content="- First item"`
- **THEN** the file is created (with parent directories) and the content is written as the file body

### Requirement: Prepend content to a file
The system SHALL prepend content after the frontmatter block (not at byte 0).

#### Scenario: Prepend to file with frontmatter
- **WHEN** user runs `vaultfs prepend notes/standup.md --content="## New Section"`
- **THEN** the content is inserted after the closing `---` of the frontmatter, before the body

#### Scenario: Prepend to file without frontmatter
- **WHEN** user runs `vaultfs prepend` on a file with no frontmatter
- **THEN** the content is inserted at the beginning of the file

### Requirement: Move a file
The system SHALL move/rename a file within the vault, creating target directories as needed.

#### Scenario: Move file
- **WHEN** user runs `vaultfs move notes/old.md --to=archive/old.md`
- **THEN** the file is moved and the `archive/` directory is created if needed

### Requirement: Delete a file
The system SHALL delete a file from the vault.

#### Scenario: Delete file
- **WHEN** user runs `vaultfs delete notes/old.md`
- **THEN** the file is permanently removed from the filesystem

### Requirement: List files
The system SHALL list files in the vault with optional filtering by folder and extension.

#### Scenario: List all markdown files
- **WHEN** user runs `vaultfs list`
- **THEN** the system outputs a JSON array of all `.md` files in the vault with their paths, sizes, and modified timestamps

#### Scenario: List files in folder
- **WHEN** user runs `vaultfs list --folder=projects`
- **THEN** only files within the `projects/` directory (recursively) are listed

#### Scenario: List by extension
- **WHEN** user runs `vaultfs list --ext=txt`
- **THEN** only `.txt` files are listed

### Requirement: List folders
The system SHALL list all directories in the vault.

#### Scenario: List folders
- **WHEN** user runs `vaultfs folders`
- **THEN** the system outputs a JSON array of all directory paths in the vault, excluding `.vaultfs/`

### Requirement: Create directory
The system SHALL create directories recursively.

#### Scenario: Create nested directory
- **WHEN** user runs `vaultfs mkdir projects/2026/q2`
- **THEN** all intermediate directories are created
