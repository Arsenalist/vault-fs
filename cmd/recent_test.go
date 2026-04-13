package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRecentDefaults(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "recent.md"), []byte("new"), 0644)

	files, err := runRecent(vaultPath, 7, 20, "")
	if err != nil {
		t.Fatalf("recent failed: %v", err)
	}
	if len(files) < 1 {
		t.Error("expected at least 1 recent file")
	}
	// Should be sorted by mtime descending
	if len(files) >= 2 {
		t1, _ := time.Parse(time.RFC3339, files[0].Modified)
		t2, _ := time.Parse(time.RFC3339, files[1].Modified)
		if t1.Before(t2) {
			t.Error("expected files sorted by mtime descending")
		}
	}
}

func TestRecentCustomWindow(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "old.md"), []byte("old"), 0644)

	// Set mtime to 60 days ago
	oldTime := time.Now().Add(-60 * 24 * time.Hour)
	os.Chtimes(filepath.Join(vaultPath, "old.md"), oldTime, oldTime)

	files, err := runRecent(vaultPath, 7, 20, "")
	if err != nil {
		t.Fatalf("recent failed: %v", err)
	}
	for _, f := range files {
		if f.Path == "old.md" {
			t.Error("old.md should not appear in 7-day window")
		}
	}
}

func TestRecentWithFolder(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "Journal", "today.md"), []byte("today"), 0644)
	os.WriteFile(filepath.Join(vaultPath, "Reports", "q1.md"), []byte("q1"), 0644)

	files, err := runRecent(vaultPath, 7, 20, "Journal")
	if err != nil {
		t.Fatalf("recent failed: %v", err)
	}
	for _, f := range files {
		if f.Path == "Reports/q1.md" {
			t.Error("should not include files outside Journal folder")
		}
	}
}

func TestRecentLimit(t *testing.T) {
	vaultPath := setupVault(t)
	for i := 0; i < 5; i++ {
		os.WriteFile(filepath.Join(vaultPath, filepath.Join("Journal", string(rune('a'+i))+".md")), []byte("x"), 0644)
	}

	files, err := runRecent(vaultPath, 7, 2, "")
	if err != nil {
		t.Fatalf("recent failed: %v", err)
	}
	if len(files) > 2 {
		t.Errorf("expected max 2 files, got %d", len(files))
	}
}

func TestRecentEmpty(t *testing.T) {
	vaultPath := setupVault(t)

	// Set all files to old
	oldTime := time.Now().Add(-60 * 24 * time.Hour)
	filepath.WalkDir(vaultPath, func(path string, d os.DirEntry, err error) error {
		if !d.IsDir() {
			os.Chtimes(path, oldTime, oldTime)
		}
		return nil
	})

	files, err := runRecent(vaultPath, 1, 20, "")
	if err != nil {
		t.Fatalf("recent failed: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}
