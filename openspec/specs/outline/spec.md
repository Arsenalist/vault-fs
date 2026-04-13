## Purpose

Defines requirements for the outline capability of vaultfs.

## ADDED Requirements

### Requirement: Extract heading outline
The system SHALL parse a markdown file and return its heading structure as a nested tree.

#### Scenario: Outline of file with headings
- **WHEN** user runs `vaultfs outline notes/design.md`
- **THEN** the system outputs a JSON tree of headings with `level`, `text`, and `children` fields

#### Scenario: File with no headings
- **WHEN** user runs `vaultfs outline` on a file with no headings
- **THEN** the system outputs an empty JSON array

#### Scenario: Nested heading structure
- **WHEN** a file has `# Title`, `## Section A`, `### Sub A1`, `## Section B`
- **THEN** the outline reflects the nesting: Title contains Section A (which contains Sub A1) and Section B
