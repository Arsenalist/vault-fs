## Purpose

Defines requirements for the tasks capability of vaultfs.

## ADDED Requirements

### Requirement: Extract tasks from vault
The system SHALL extract markdown checkbox tasks (`- [ ]` and `- [x]`) from vault files with optional metadata parsing.

#### Scenario: List all tasks
- **WHEN** user runs `vaultfs tasks`
- **THEN** the system outputs a JSON array of all tasks with fields: `file`, `line`, `text`, `done`, `priority`, `due`, `tags`, `mentions`

#### Scenario: Filter pending tasks
- **WHEN** user runs `vaultfs tasks --pending`
- **THEN** only tasks with `- [ ]` (not done) are returned

#### Scenario: Filter completed tasks
- **WHEN** user runs `vaultfs tasks --done`
- **THEN** only tasks with `- [x]` (done) are returned

#### Scenario: Filter by folder
- **WHEN** user runs `vaultfs tasks --folder=projects`
- **THEN** only tasks from files within `projects/` are returned

### Requirement: Parse task priority
The system SHALL extract priority from emoji markers in task text.

#### Scenario: High priority
- **WHEN** a task contains `🔴` or `⏫`
- **THEN** the `priority` field is `"high"`

#### Scenario: Medium priority
- **WHEN** a task contains `🟡` or `🔼`
- **THEN** the `priority` field is `"medium"`

#### Scenario: Low priority
- **WHEN** a task contains `🔵` or `🔽`
- **THEN** the `priority` field is `"low"`

#### Scenario: No priority
- **WHEN** a task has no priority emoji
- **THEN** the `priority` field is `null`

### Requirement: Parse task due date
The system SHALL extract due dates from task text using common patterns.

#### Scenario: Due tag format
- **WHEN** a task contains `#due/2026-04-15`
- **THEN** the `due` field is `"2026-04-15"`

#### Scenario: Calendar emoji format
- **WHEN** a task contains `📅 2026-04-15`
- **THEN** the `due` field is `"2026-04-15"`

#### Scenario: No due date
- **WHEN** a task has no due date pattern
- **THEN** the `due` field is `null`

### Requirement: Parse task mentions and tags
The system SHALL extract inline `@mentions` and `#tags` from task text.

#### Scenario: Task with mention
- **WHEN** a task contains `@alice`
- **THEN** `"alice"` appears in the `mentions` array

#### Scenario: Task with tags
- **WHEN** a task contains `#urgent #backend`
- **THEN** both `"urgent"` and `"backend"` appear in the `tags` array (excluding `#due/...` patterns)

### Requirement: Toggle task status
The system SHALL toggle a task's checkbox between `- [ ]` and `- [x]` at a specific file and line.

#### Scenario: Toggle incomplete to complete
- **WHEN** user runs `vaultfs task toggle notes/todo.md --line=5` and line 5 is `- [ ] Buy milk`
- **THEN** line 5 becomes `- [x] Buy milk`

#### Scenario: Toggle complete to incomplete
- **WHEN** user runs `vaultfs task toggle notes/todo.md --line=5` and line 5 is `- [x] Buy milk`
- **THEN** line 5 becomes `- [ ] Buy milk`
