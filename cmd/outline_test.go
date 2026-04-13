package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOutlineNested(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "doc.md"), []byte("# Title\n## Section A\n### Sub A1\n## Section B"), 0644)

	outline, err := runOutline(vaultPath, "doc.md")
	if err != nil {
		t.Fatalf("outline failed: %v", err)
	}
	if len(outline) != 1 {
		t.Fatalf("expected 1 root, got %d", len(outline))
	}
	if outline[0].Text != "Title" {
		t.Errorf("expected Title, got %s", outline[0].Text)
	}
	if len(outline[0].Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(outline[0].Children))
	}
}

func TestOutlineNoHeadings(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "plain.md"), []byte("Just text"), 0644)

	outline, err := runOutline(vaultPath, "plain.md")
	if err != nil {
		t.Fatalf("outline failed: %v", err)
	}
	if len(outline) != 0 {
		t.Errorf("expected empty, got %d", len(outline))
	}
}

func TestOutlineSkippedLevels(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "skip.md"), []byte("# Title\n### Skipped to H3\n## Back to H2"), 0644)

	outline, err := runOutline(vaultPath, "skip.md")
	if err != nil {
		t.Fatalf("outline failed: %v", err)
	}
	if len(outline) != 1 {
		t.Fatalf("expected 1 root, got %d", len(outline))
	}
	if len(outline[0].Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(outline[0].Children))
	}
}
