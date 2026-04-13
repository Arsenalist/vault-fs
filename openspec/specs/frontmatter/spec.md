## Purpose

Defines requirements for the frontmatter capability of vaultfs.

## ADDED Requirements

### Requirement: Read all frontmatter properties
The system SHALL parse and return all YAML frontmatter properties from a markdown file.

#### Scenario: Read properties
- **WHEN** user runs `vaultfs properties notes/standup.md`
- **THEN** the system outputs a JSON object of all frontmatter key-value pairs

#### Scenario: No frontmatter
- **WHEN** user runs `vaultfs properties` on a file with no frontmatter
- **THEN** the system outputs an empty JSON object `{}`

### Requirement: Set a frontmatter property
The system SHALL set a YAML frontmatter property on a markdown file, creating the frontmatter block if it doesn't exist.

#### Scenario: Set property on file with frontmatter
- **WHEN** user runs `vaultfs property set notes/standup.md --name=status --value=active`
- **THEN** the `status` property is set to `active` in the existing frontmatter

#### Scenario: Set property on file without frontmatter
- **WHEN** user runs `vaultfs property set notes/plain.md --name=tags --value="[work, urgent]"`
- **THEN** a frontmatter block is created at the top of the file with the property set

### Requirement: Remove a frontmatter property
The system SHALL remove a specified property from a file's frontmatter.

#### Scenario: Remove property
- **WHEN** user runs `vaultfs property remove notes/standup.md --name=draft`
- **THEN** the `draft` key is removed from the frontmatter

#### Scenario: Remove nonexistent property
- **WHEN** user runs `vaultfs property remove` for a property that doesn't exist
- **THEN** the system succeeds silently (idempotent)
