package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupSearchableVault(t *testing.T) string {
	t.Helper()
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "notes.md"), []byte("# Meeting Notes\n\nDiscussed quarterly budget and team allocation.\nAction items pending."), 0644)
	os.WriteFile(filepath.Join(vaultPath, "Journal", "day1.md"), []byte("# Day 1\n\nBudget review went well. Next steps decided."), 0644)
	os.WriteFile(filepath.Join(vaultPath, "Reports", "q1.md"), []byte("# Q1 Report\n\nPerformance exceeded expectations."), 0644)
	return vaultPath
}

func TestSearchFullText(t *testing.T) {
	vaultPath := setupSearchableVault(t)

	results, err := runSearch(vaultPath, "budget", "", 10, false, false)
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(results) < 1 {
		t.Fatal("expected at least 1 result for 'budget'")
	}
}

func TestSearchWithFolder(t *testing.T) {
	vaultPath := setupSearchableVault(t)

	results, err := runSearch(vaultPath, "budget", "Journal", 10, false, false)
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	for _, r := range results {
		if !strings.HasPrefix(r.Path, "Journal/") {
			t.Errorf("expected result in Journal, got %s", r.Path)
		}
	}
}

func TestSearchWithLimit(t *testing.T) {
	vaultPath := setupSearchableVault(t)

	results, err := runSearch(vaultPath, "budget", "", 1, false, false)
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(results) > 1 {
		t.Errorf("expected max 1 result, got %d", len(results))
	}
}

func TestSearchContext(t *testing.T) {
	vaultPath := setupSearchableVault(t)

	results, err := runSearchContext(vaultPath, "budget", 10)
	if err != nil {
		t.Fatalf("search:context failed: %v", err)
	}
	if len(results) < 1 {
		t.Fatal("expected at least 1 context result")
	}
	// Each result should have matching lines
	for _, r := range results {
		if len(r.Matches) == 0 {
			t.Errorf("expected matches for %s", r.Path)
		}
	}
}

func TestSearchFuzzy(t *testing.T) {
	vaultPath := setupSearchableVault(t)

	results, err := runSearch(vaultPath, "notes", "", 10, true, false)
	if err != nil {
		t.Fatalf("fuzzy search failed: %v", err)
	}
	// Should find files with "notes" in name
	if len(results) < 1 {
		t.Error("expected fuzzy match to find files containing 'notes'")
	}
}

func TestIndexRebuildCmd(t *testing.T) {
	vaultPath := setupSearchableVault(t)

	count, err := runIndexRebuild(vaultPath)
	if err != nil {
		t.Fatalf("index rebuild failed: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 indexed files, got %d", count)
	}
}
