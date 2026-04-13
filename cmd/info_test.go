package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func setupVault(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	runInit(tmp, "basic", nil)
	// Create some files
	os.WriteFile(filepath.Join(tmp, "notes.md"), []byte("# Notes"), 0644)
	os.WriteFile(filepath.Join(tmp, "Journal", "day1.md"), []byte("# Day 1"), 0644)
	return tmp
}

func TestInfoReturnsVaultMetadata(t *testing.T) {
	vaultPath := setupVault(t)

	info, err := getVaultInfo(vaultPath)
	if err != nil {
		t.Fatalf("getVaultInfo failed: %v", err)
	}

	if info.Path != vaultPath {
		t.Errorf("expected path %s, got %s", vaultPath, info.Path)
	}
	if info.FileCount < 3 { // README.md + notes.md + Journal/day1.md
		t.Errorf("expected at least 3 files, got %d", info.FileCount)
	}
	if info.FolderCount < 9 { // basic preset dirs
		t.Errorf("expected at least 9 folders, got %d", info.FolderCount)
	}

	// Should be marshallable to JSON
	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("failed to marshal info: %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("info is not valid JSON: %v", err)
	}
	if _, ok := result["path"]; !ok {
		t.Error("expected 'path' in JSON output")
	}
}
