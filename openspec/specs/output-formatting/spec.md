## Purpose

Defines requirements for the output-formatting capability of vaultfs.

## ADDED Requirements

### Requirement: JSON default for query commands
Query commands (read, list, folders, recent, search, search:context, tags, tag, tasks, properties, outline, info) SHALL output JSON by default.

#### Scenario: Default query output
- **WHEN** user runs `vaultfs list` with no `--format` flag
- **THEN** the output is valid JSON

#### Scenario: Text override for queries
- **WHEN** user runs `vaultfs list --format=text`
- **THEN** the output is human-readable plain text

### Requirement: Text default for action commands
Action commands (init, create, append, prepend, move, delete, mkdir, property set, property remove, task toggle, index rebuild) SHALL output human-readable text by default.

#### Scenario: Default action output
- **WHEN** user runs `vaultfs create notes/foo --content="hello"`
- **THEN** the output is a human-readable confirmation message

#### Scenario: JSON override for actions
- **WHEN** user runs `vaultfs create notes/foo --content="hello" --format=json`
- **THEN** the output is a structured JSON response

### Requirement: Consistent error format
All errors SHALL be written to stderr. When `--format=json` is active, errors SHALL also be JSON-formatted with an `error` field.

#### Scenario: Error in text mode
- **WHEN** a command fails with no format flag
- **THEN** a human-readable error message is written to stderr with exit code 1

#### Scenario: Error in JSON mode
- **WHEN** a command fails with `--format=json`
- **THEN** `{"error": "description"}` is written to stderr with exit code 1
