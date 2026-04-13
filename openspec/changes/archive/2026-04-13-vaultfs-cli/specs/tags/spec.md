## ADDED Requirements

### Requirement: List all tags in vault
The system SHALL extract and list all tags from both YAML frontmatter `tags:` arrays and inline `#tag` syntax (including nested `#tag/sub`).

#### Scenario: List tags
- **WHEN** user runs `vaultfs tags`
- **THEN** the system outputs a JSON array of unique tag names found across all vault files

#### Scenario: List tags with counts
- **WHEN** user runs `vaultfs tags --counts`
- **THEN** each tag includes a usage count (number of files containing that tag)

#### Scenario: List tags sorted by count
- **WHEN** user runs `vaultfs tags --counts --sort=count`
- **THEN** tags are sorted by usage count descending

### Requirement: List files with a specific tag
The system SHALL list all files containing a given tag.

#### Scenario: Filter by tag
- **WHEN** user runs `vaultfs tag project`
- **THEN** the system outputs a JSON array of file paths that contain the tag `project` (from frontmatter or inline)

#### Scenario: Filter by nested tag
- **WHEN** user runs `vaultfs tag "project/alpha"`
- **THEN** files containing the nested tag `#project/alpha` or frontmatter tag `project/alpha` are listed

### Requirement: Tag extraction from both sources
The system SHALL recognize tags from YAML frontmatter `tags:` field (both array `[a, b]` and string formats) and inline `#tag` syntax in markdown body content.

#### Scenario: Frontmatter tags
- **WHEN** a file has frontmatter `tags: [work, urgent]`
- **THEN** both `work` and `urgent` are recognized as tags for that file

#### Scenario: Inline tags
- **WHEN** a file body contains `#meeting` and `#project/alpha`
- **THEN** both `meeting` and `project/alpha` are recognized as tags for that file

#### Scenario: Deduplicated tags
- **WHEN** a file has `urgent` in both frontmatter and as inline `#urgent`
- **THEN** the tag is counted once for that file
