## ADDED Requirements

### Requirement: Initialize a new vault
The system SHALL create a new vault at a specified path (defaulting to `~/vault-fs`) with a `.vaultfs/` configuration directory and a `README.md` file.

#### Scenario: Init with default path
- **WHEN** user runs `vaultfs init` with no `--path` flag
- **THEN** a vault is created at `~/vault-fs` with `.vaultfs/config.yaml` and `README.md`

#### Scenario: Init with custom path
- **WHEN** user runs `vaultfs init --path=/tmp/my-vault`
- **THEN** a vault is created at `/tmp/my-vault` with `.vaultfs/config.yaml` and `README.md`

#### Scenario: Init with preset
- **WHEN** user runs `vaultfs init --preset=basic`
- **THEN** the vault is created with all directories defined in the `basic` preset from the embedded default config, and the preset definition is written into `.vaultfs/config.yaml`

#### Scenario: Init with custom directories
- **WHEN** user runs `vaultfs init --dirs="inbox,projects/active,projects/archive"`
- **THEN** all specified directories are created recursively (including nested paths), in addition to any preset directories

#### Scenario: Init with preset and extra directories
- **WHEN** user runs `vaultfs init --preset=basic --dirs="clients/acme"`
- **THEN** both the preset directories and the extra custom directories are created, with no duplicates

#### Scenario: Init on existing vault
- **WHEN** user runs `vaultfs init` on a path that already contains `.vaultfs/`
- **THEN** the system SHALL report an error indicating the vault already exists

### Requirement: List available presets
The system SHALL allow users to list all available init presets.

#### Scenario: List presets
- **WHEN** user runs `vaultfs init --list-presets`
- **THEN** the system outputs all preset names and their directory lists in JSON format

### Requirement: Config-driven presets
Presets SHALL be defined in `.vaultfs/config.yaml`, not hardcoded in the binary. The binary SHALL embed a default config containing the `basic` preset which is written on init.

#### Scenario: Basic preset contents
- **WHEN** the `basic` preset is used
- **THEN** the following directories are created: Daily Debrief, Daily Plan, Journal, Meeting Notes, Projects/Active, Projects/Archived, Reports, Scratchpad, Stakeholders

#### Scenario: Custom preset in config
- **WHEN** a user has added a custom preset named `team` to `.vaultfs/config.yaml`
- **THEN** `vaultfs init --preset=team` in a new vault uses that preset's directory list

### Requirement: Display vault info
The system SHALL display vault metadata including path, file count, folder count, and index status.

#### Scenario: Vault info
- **WHEN** user runs `vaultfs info`
- **THEN** the system outputs vault path, total file count, total folder count, last index time, and config path as JSON

### Requirement: Vault discovery
The system SHALL resolve the vault location using a priority chain: `--vault` flag, `VAULTFS_PATH` env var, walk up from CWD for `.vaultfs/`, fall back to `~/vault-fs`.

#### Scenario: Vault from env var
- **WHEN** `VAULTFS_PATH` is set to `/tmp/my-vault` and no `--vault` flag is provided
- **THEN** the system operates on `/tmp/my-vault`

#### Scenario: Vault from CWD
- **WHEN** CWD is `/tmp/my-vault/notes/` and `/tmp/my-vault/.vaultfs/` exists
- **THEN** the system operates on `/tmp/my-vault`

#### Scenario: Vault fallback
- **WHEN** no flag, env var, or `.vaultfs/` ancestor exists
- **THEN** the system operates on `~/vault-fs`
