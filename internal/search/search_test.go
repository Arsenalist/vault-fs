package search

import (
	"os"
	"path/filepath"
	"testing"
)

func setupSearchVault(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	os.MkdirAll(filepath.Join(tmp, ".vaultfs"), 0755)
	os.MkdirAll(filepath.Join(tmp, "notes"), 0755)
	os.MkdirAll(filepath.Join(tmp, "projects"), 0755)

	os.WriteFile(filepath.Join(tmp, "notes", "meeting.md"), []byte("# Meeting Notes\n\nDiscussed the quarterly review and budget allocation.\nAction items for the team."), 0644)
	os.WriteFile(filepath.Join(tmp, "notes", "standup.md"), []byte("# Standup\n\nWorked on the API integration. Blocked on auth tokens."), 0644)
	os.WriteFile(filepath.Join(tmp, "projects", "alpha.md"), []byte("# Project Alpha\n\nBudget is approved. Starting next sprint."), 0644)
	return tmp
}

func TestIndexCreateAndOpen(t *testing.T) {
	vaultPath := setupSearchVault(t)
	idx, err := OpenOrCreate(vaultPath)
	if err != nil {
		t.Fatalf("OpenOrCreate failed: %v", err)
	}
	defer idx.Close()

	// Index directory should exist
	indexPath := filepath.Join(vaultPath, ".vaultfs", "index.bleve")
	if _, err := os.Stat(indexPath); err != nil {
		t.Error("expected index directory to exist")
	}
}

func TestIndexRebuild(t *testing.T) {
	vaultPath := setupSearchVault(t)
	idx, err := OpenOrCreate(vaultPath)
	if err != nil {
		t.Fatalf("OpenOrCreate failed: %v", err)
	}

	count, err := idx.Rebuild(vaultPath)
	if err != nil {
		t.Fatalf("Rebuild failed: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 indexed docs, got %d", count)
	}
	idx.Close()
}

func TestIndexQuery(t *testing.T) {
	vaultPath := setupSearchVault(t)
	idx, err := OpenOrCreate(vaultPath)
	if err != nil {
		t.Fatalf("OpenOrCreate failed: %v", err)
	}
	defer idx.Close()

	idx.Rebuild(vaultPath)

	results, err := idx.Search("budget", "", 10, false)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) < 1 {
		t.Fatal("expected at least 1 result for 'budget'")
	}
	// Both meeting.md and alpha.md mention budget
	if len(results) != 2 {
		t.Errorf("expected 2 results for 'budget', got %d", len(results))
	}
}

func TestIndexQueryWithFolder(t *testing.T) {
	vaultPath := setupSearchVault(t)
	idx, err := OpenOrCreate(vaultPath)
	if err != nil {
		t.Fatalf("OpenOrCreate failed: %v", err)
	}
	defer idx.Close()

	idx.Rebuild(vaultPath)

	results, err := idx.Search("budget", "projects", 10, false)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result in projects folder, got %d", len(results))
	}
}

func TestStaleDetection(t *testing.T) {
	vaultPath := setupSearchVault(t)
	idx, err := OpenOrCreate(vaultPath)
	if err != nil {
		t.Fatalf("OpenOrCreate failed: %v", err)
	}
	idx.Rebuild(vaultPath)
	idx.Close()

	// Should not be stale right after rebuild
	stale, err := IsStale(vaultPath)
	if err != nil {
		t.Fatalf("IsStale failed: %v", err)
	}
	if stale {
		t.Error("expected index to not be stale right after rebuild")
	}

	// Modify a file
	os.WriteFile(filepath.Join(vaultPath, "notes", "new.md"), []byte("new content"), 0644)

	stale, err = IsStale(vaultPath)
	if err != nil {
		t.Fatalf("IsStale failed: %v", err)
	}
	if !stale {
		t.Error("expected index to be stale after file modification")
	}
}
