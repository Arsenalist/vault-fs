## Purpose

Defines requirements for the agent-skill capability of vaultfs.

## ADDED Requirements

### Requirement: Claude Code skill definition
The project SHALL include a Claude Code skill file (`skills/vaultfs/SKILL.md`) that teaches agents how to use the `vaultfs` CLI with all commands documented, common patterns, and tips.

#### Scenario: Skill triggers on vault operations
- **WHEN** a user asks an agent to interact with their markdown vault (create notes, search, list tasks, manage tags, etc.)
- **THEN** the skill is triggered and the agent has full context on `vaultfs` command syntax, flags, and output formats

#### Scenario: Skill documents all commands
- **WHEN** an agent reads the skill file
- **THEN** it contains a command overview table, quick reference with examples for every command group (vault management, file CRUD, search, tags, tasks, properties, outline, recent), output format details, and vault discovery behavior

#### Scenario: Skill includes common agent patterns
- **WHEN** an agent needs to perform multi-step vault operations
- **THEN** the skill provides ready-to-use patterns (e.g., create note + set properties, search + read results, task extraction workflows)

### Requirement: Command reference document
The project SHALL include a full command reference (`skills/vaultfs/references/command-reference.md`) with complete parameter tables for every command and subcommand.

#### Scenario: Agent needs specific flags
- **WHEN** an agent needs to know the exact flags for a command
- **THEN** the command reference provides parameter names, types, defaults, and descriptions for every command
