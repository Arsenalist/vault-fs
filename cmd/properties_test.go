package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zarar/vaultfs/internal/vfs"
)

func TestPropertiesRead(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "note.md"), []byte("---\ntitle: Test\nstatus: active\n---\n\nBody"), 0644)

	props, err := runProperties(vaultPath, "note.md")
	if err != nil {
		t.Fatalf("properties failed: %v", err)
	}
	if props["title"] != "Test" {
		t.Errorf("expected title=Test, got %v", props["title"])
	}
	if props["status"] != "active" {
		t.Errorf("expected status=active, got %v", props["status"])
	}
}

func TestPropertiesReadNoFrontmatter(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "plain.md"), []byte("No frontmatter"), 0644)

	props, err := runProperties(vaultPath, "plain.md")
	if err != nil {
		t.Fatalf("properties failed: %v", err)
	}
	if len(props) != 0 {
		t.Errorf("expected empty, got %v", props)
	}
}

func TestPropertySet(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "note.md"), []byte("---\ntitle: Test\n---\n\nBody"), 0644)

	err := runPropertySet(vaultPath, "note.md", "status", "active")
	if err != nil {
		t.Fatalf("property set failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(vaultPath, "note.md"))
	if !strings.Contains(string(data), "status: active") {
		t.Errorf("expected status in frontmatter, got: %s", string(data))
	}
	if !strings.Contains(string(data), "Body") {
		t.Error("expected body preserved")
	}
}

func TestPropertySetCreatesIfMissing(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "plain.md"), []byte("No frontmatter content"), 0644)

	err := runPropertySet(vaultPath, "plain.md", "status", "draft")
	if err != nil {
		t.Fatalf("property set failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(vaultPath, "plain.md"))
	content := string(data)
	if !strings.HasPrefix(content, "---") {
		t.Error("expected frontmatter to be created")
	}
	if !strings.Contains(content, "status: draft") {
		t.Error("expected status in frontmatter")
	}
	if !strings.Contains(content, "No frontmatter content") {
		t.Error("expected original body preserved")
	}
}

func TestPropertyRemove(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "note.md"), []byte("---\ntitle: Test\ndraft: true\n---\n\nBody"), 0644)

	err := runPropertyRemove(vaultPath, "note.md", "draft")
	if err != nil {
		t.Fatalf("property remove failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(vaultPath, "note.md"))
	if strings.Contains(string(data), "draft") {
		t.Error("expected draft to be removed")
	}
	if !strings.Contains(string(data), "title: Test") {
		t.Error("expected title to be preserved")
	}
}

func TestPropertyRemoveIdempotent(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "note.md"), []byte("---\ntitle: Test\n---\n\nBody"), 0644)

	err := runPropertyRemove(vaultPath, "note.md", "nonexistent")
	if err != nil {
		t.Fatalf("property remove should be idempotent, got: %v", err)
	}
}

func TestPropertiesNotFound(t *testing.T) {
	vaultPath := setupVault(t)

	_, err := runProperties(vaultPath, "missing.md")
	var nf *vfs.NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("expected *vfs.NotFoundError, got %T: %v", err, err)
	}
	if nf.Path != "missing.md" {
		t.Errorf("expected Path=missing.md, got %q", nf.Path)
	}
}

func TestPropertySetAutoCreatesMissingFile(t *testing.T) {
	vaultPath := setupVault(t)

	err := runPropertySet(vaultPath, "new/note.md", "tag", "draft")
	if err != nil {
		t.Fatalf("property set on missing file should auto-create, got: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(vaultPath, "new", "note.md"))
	if err != nil {
		t.Fatalf("expected file created: %v", err)
	}
	content := string(data)
	if !strings.HasPrefix(content, "---") {
		t.Errorf("expected frontmatter at start, got: %q", content)
	}
	if !strings.Contains(content, "tag: draft") {
		t.Errorf("expected tag in frontmatter, got: %q", content)
	}
}

func TestPropertyRemoveNotFound(t *testing.T) {
	vaultPath := setupVault(t)

	err := runPropertyRemove(vaultPath, "missing.md", "any")
	var nf *vfs.NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("expected *vfs.NotFoundError, got %T: %v", err, err)
	}
	if nf.Path != "missing.md" {
		t.Errorf("expected Path=missing.md, got %q", nf.Path)
	}
}
