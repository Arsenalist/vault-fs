package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverFromExplicitFlag(t *testing.T) {
	tmp := t.TempDir()
	os.MkdirAll(filepath.Join(tmp, ".vaultfs"), 0755)

	path, err := Discover(tmp, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != tmp {
		t.Errorf("expected %s, got %s", tmp, path)
	}
}

func TestDiscoverFromEnvVar(t *testing.T) {
	tmp := t.TempDir()
	os.MkdirAll(filepath.Join(tmp, ".vaultfs"), 0755)

	t.Setenv("VAULTFS_PATH", tmp)

	path, err := Discover("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != tmp {
		t.Errorf("expected %s, got %s", tmp, path)
	}
}

func TestDiscoverWalkUp(t *testing.T) {
	tmp := t.TempDir()
	os.MkdirAll(filepath.Join(tmp, ".vaultfs"), 0755)
	subdir := filepath.Join(tmp, "notes", "deep")
	os.MkdirAll(subdir, 0755)

	path, err := Discover("", subdir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != tmp {
		t.Errorf("expected %s, got %s", tmp, path)
	}
}

func TestDiscoverFallbackToDefault(t *testing.T) {
	// Unset env, use a CWD without .vaultfs
	t.Setenv("VAULTFS_PATH", "")
	tmp := t.TempDir()

	path, err := Discover("", tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, "vault-fs")
	if path != expected {
		t.Errorf("expected fallback %s, got %s", expected, path)
	}
}

func TestDiscoverFlagTakesPrecedenceOverEnv(t *testing.T) {
	flagDir := t.TempDir()
	envDir := t.TempDir()
	os.MkdirAll(filepath.Join(flagDir, ".vaultfs"), 0755)
	os.MkdirAll(filepath.Join(envDir, ".vaultfs"), 0755)

	t.Setenv("VAULTFS_PATH", envDir)

	path, err := Discover(flagDir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != flagDir {
		t.Errorf("expected flag path %s, got %s", flagDir, path)
	}
}

func TestLoadConfig(t *testing.T) {
	tmp := t.TempDir()
	configDir := filepath.Join(tmp, ".vaultfs")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(`
vault:
  path: /tmp/test-vault

presets:
  basic:
    directories:
      - Notes
      - Projects/Active
`), 0644)

	cfg, err := LoadConfig(tmp)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Vault.Path != "/tmp/test-vault" {
		t.Errorf("expected vault path /tmp/test-vault, got %s", cfg.Vault.Path)
	}
	preset, ok := cfg.Presets["basic"]
	if !ok {
		t.Fatal("expected 'basic' preset")
	}
	if len(preset.Directories) != 2 {
		t.Errorf("expected 2 directories, got %d", len(preset.Directories))
	}
	if preset.Directories[0] != "Notes" {
		t.Errorf("expected first dir to be Notes, got %s", preset.Directories[0])
	}
}

func TestLoadConfigFallsBackToDefaults(t *testing.T) {
	tmp := t.TempDir()
	// No .vaultfs/config.yaml exists

	cfg, err := LoadConfig(tmp)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	// Should have the embedded default preset
	preset, ok := cfg.Presets["basic"]
	if !ok {
		t.Fatal("expected 'basic' preset from defaults")
	}
	if len(preset.Directories) == 0 {
		t.Error("expected default preset to have directories")
	}
}

func TestLoadConfigMergesWithDefaults(t *testing.T) {
	tmp := t.TempDir()
	configDir := filepath.Join(tmp, ".vaultfs")
	os.MkdirAll(configDir, 0755)
	// Config that adds a custom preset but doesn't include basic
	os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(`
presets:
  custom:
    directories:
      - Custom/Dir
`), 0644)

	cfg, err := LoadConfig(tmp)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	// Should still have the embedded basic preset
	if _, ok := cfg.Presets["basic"]; !ok {
		t.Error("expected 'basic' preset from defaults to be preserved")
	}
	// And the custom one
	if _, ok := cfg.Presets["custom"]; !ok {
		t.Error("expected 'custom' preset")
	}
}
