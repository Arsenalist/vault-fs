package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInitDefaultPath(t *testing.T) {
	tmp := t.TempDir()
	vaultPath := filepath.Join(tmp, "vault-fs")

	err := runInit(vaultPath, "", nil)
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// .vaultfs/config.yaml should exist
	if _, err := os.Stat(filepath.Join(vaultPath, ".vaultfs", "config.yaml")); err != nil {
		t.Error("expected .vaultfs/config.yaml to exist")
	}
	// README.md should exist
	if _, err := os.Stat(filepath.Join(vaultPath, "README.md")); err != nil {
		t.Error("expected README.md to exist")
	}
}

func TestInitWithPreset(t *testing.T) {
	tmp := t.TempDir()
	vaultPath := filepath.Join(tmp, "vault")

	err := runInit(vaultPath, "basic", nil)
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	expectedDirs := []string{
		"Daily Debrief",
		"Daily Plan",
		"Journal",
		"Meeting Notes",
		"Projects/Active",
		"Projects/Archived",
		"Reports",
		"Scratchpad",
		"Stakeholders",
	}

	for _, d := range expectedDirs {
		fullPath := filepath.Join(vaultPath, d)
		if info, err := os.Stat(fullPath); err != nil || !info.IsDir() {
			t.Errorf("expected directory %s to exist", d)
		}
	}
}

func TestInitWithCustomDirs(t *testing.T) {
	tmp := t.TempDir()
	vaultPath := filepath.Join(tmp, "vault")

	err := runInit(vaultPath, "", []string{"inbox", "projects/active", "projects/archive"})
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	for _, d := range []string{"inbox", "projects/active", "projects/archive"} {
		fullPath := filepath.Join(vaultPath, d)
		if info, err := os.Stat(fullPath); err != nil || !info.IsDir() {
			t.Errorf("expected directory %s to exist", d)
		}
	}
}

func TestInitPresetPlusDirs(t *testing.T) {
	tmp := t.TempDir()
	vaultPath := filepath.Join(tmp, "vault")

	err := runInit(vaultPath, "basic", []string{"clients/acme"})
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Check a preset dir and the custom dir both exist
	if info, err := os.Stat(filepath.Join(vaultPath, "Journal")); err != nil || !info.IsDir() {
		t.Error("expected preset dir Journal to exist")
	}
	if info, err := os.Stat(filepath.Join(vaultPath, "clients", "acme")); err != nil || !info.IsDir() {
		t.Error("expected custom dir clients/acme to exist")
	}
}

func TestInitErrorOnExistingVault(t *testing.T) {
	tmp := t.TempDir()
	vaultPath := filepath.Join(tmp, "vault")
	os.MkdirAll(filepath.Join(vaultPath, ".vaultfs"), 0755)

	err := runInit(vaultPath, "", nil)
	if err == nil {
		t.Fatal("expected error for existing vault, got nil")
	}
}

func TestInitListPresets(t *testing.T) {
	result, err := listPresets()
	if err != nil {
		t.Fatalf("listPresets failed: %v", err)
	}

	// Should be valid JSON with a "basic" key
	var presets map[string]any
	if err := json.Unmarshal([]byte(result), &presets); err != nil {
		t.Fatalf("listPresets output is not valid JSON: %v", err)
	}
	if _, ok := presets["basic"]; !ok {
		t.Error("expected 'basic' preset in output")
	}
}
