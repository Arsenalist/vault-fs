package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestTagsAll(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "a.md"), []byte("---\ntags: [work, urgent]\n---\n\n#meeting"), 0644)
	os.WriteFile(filepath.Join(vaultPath, "b.md"), []byte("#work #personal"), 0644)

	result, err := runTags(vaultPath, false, "")
	if err != nil {
		t.Fatalf("tags failed: %v", err)
	}

	var names []string
	for _, tr := range result {
		names = append(names, tr.Name)
	}
	sort.Strings(names)

	expected := []string{"meeting", "personal", "urgent", "work"}
	if len(names) != len(expected) {
		t.Fatalf("expected %d tags, got %d: %v", len(expected), len(names), names)
	}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], name)
		}
	}
}

func TestTagsWithCounts(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "a.md"), []byte("#work"), 0644)
	os.WriteFile(filepath.Join(vaultPath, "b.md"), []byte("#work #personal"), 0644)

	result, err := runTags(vaultPath, true, "count")
	if err != nil {
		t.Fatalf("tags failed: %v", err)
	}

	// "work" should be first (most used)
	if result[0].Name != "work" {
		t.Errorf("expected 'work' first, got %s", result[0].Name)
	}
	if result[0].Count != 2 {
		t.Errorf("expected count 2 for 'work', got %d", result[0].Count)
	}
}

func TestTagByName(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "a.md"), []byte("#project"), 0644)
	os.WriteFile(filepath.Join(vaultPath, "b.md"), []byte("#other"), 0644)
	os.WriteFile(filepath.Join(vaultPath, "c.md"), []byte("---\ntags: [project]\n---\n"), 0644)

	files, err := runTagByName(vaultPath, "project")
	if err != nil {
		t.Fatalf("tag failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("expected 2 files with tag 'project', got %d: %v", len(files), files)
	}
}

func TestTagByNameNested(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "a.md"), []byte("#project/alpha"), 0644)

	files, err := runTagByName(vaultPath, "project/alpha")
	if err != nil {
		t.Fatalf("tag failed: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}
}
