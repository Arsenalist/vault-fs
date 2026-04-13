## ADDED Requirements

### Requirement: List recently modified files
The system SHALL list files sorted by modification time (most recent first), with configurable time window and result limit.

#### Scenario: Recent files with defaults
- **WHEN** user runs `vaultfs recent`
- **THEN** the system outputs a JSON array of files modified in the last 7 days, limited to 20 results, sorted by mtime descending

#### Scenario: Recent files with custom window
- **WHEN** user runs `vaultfs recent --days=30 --limit=50`
- **THEN** files modified in the last 30 days are listed, up to 50 results

#### Scenario: Recent files in folder
- **WHEN** user runs `vaultfs recent --folder=projects`
- **THEN** only recently modified files within `projects/` are listed

#### Scenario: No recent files
- **WHEN** no files have been modified within the specified time window
- **THEN** the system outputs an empty JSON array
