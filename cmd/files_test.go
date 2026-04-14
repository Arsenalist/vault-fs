package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// --- create ---

func TestCreateWithContent(t *testing.T) {
	vaultPath := setupVault(t)

	err := runCreate(vaultPath, "notes/standup", "# Standup\n\nToday's notes", false)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(vaultPath, "notes", "standup.md"))
	if err != nil {
		t.Fatal("file not created")
	}
	if string(data) != "# Standup\n\nToday's notes" {
		t.Errorf("unexpected content: %q", string(data))
	}
}

func TestCreateAutoMdExtension(t *testing.T) {
	vaultPath := setupVault(t)

	err := runCreate(vaultPath, "test", "hello", false)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(vaultPath, "test.md")); err != nil {
		t.Error("expected test.md to exist")
	}
}

func TestCreateAutoParentDirs(t *testing.T) {
	vaultPath := setupVault(t)

	err := runCreate(vaultPath, "deep/nested/note", "content", false)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(vaultPath, "deep", "nested", "note.md")); err != nil {
		t.Error("expected deep/nested/note.md to exist")
	}
}

func TestCreateErrorOnExisting(t *testing.T) {
	vaultPath := setupVault(t)

	runCreate(vaultPath, "existing", "first", false)
	err := runCreate(vaultPath, "existing", "second", false)
	if err == nil {
		t.Error("expected error for existing file")
	}
}

func TestCreateAppendOnExisting(t *testing.T) {
	vaultPath := setupVault(t)

	runCreate(vaultPath, "note", "first", false)
	err := runCreate(vaultPath, "note", "\nsecond", true)
	if err != nil {
		t.Fatalf("create --append failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(vaultPath, "note.md"))
	if string(data) != "first\nsecond" {
		t.Errorf("unexpected content: %q", string(data))
	}
}

// --- read ---

func TestReadWithFrontmatter(t *testing.T) {
	vaultPath := setupVault(t)
	content := "---\ntitle: Test\ntags: [work]\n---\n\n# Body\n\nContent here."
	os.WriteFile(filepath.Join(vaultPath, "note.md"), []byte(content), 0644)

	result, err := runRead(vaultPath, "note.md")
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if result.Path != "note.md" {
		t.Errorf("expected path 'note.md', got %s", result.Path)
	}
	if result.Properties["title"] != "Test" {
		t.Errorf("expected title 'Test', got %v", result.Properties["title"])
	}
	if result.Body != "# Body\n\nContent here." {
		t.Errorf("unexpected body: %q", result.Body)
	}
	if result.Size == 0 {
		t.Error("expected non-zero size")
	}
	if result.Modified == "" {
		t.Error("expected non-empty modified timestamp")
	}
}

func TestReadWithoutFrontmatter(t *testing.T) {
	vaultPath := setupVault(t)
	content := "# Plain\n\nNo frontmatter."
	os.WriteFile(filepath.Join(vaultPath, "plain.md"), []byte(content), 0644)

	result, err := runRead(vaultPath, "plain.md")
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if len(result.Properties) != 0 {
		t.Errorf("expected empty properties, got %v", result.Properties)
	}
	if result.Body != content {
		t.Errorf("expected full content as body, got %q", result.Body)
	}
}

func TestReadNonexistent(t *testing.T) {
	vaultPath := setupVault(t)

	_, err := runRead(vaultPath, "missing.md")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestReadJSONMarshal(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "test.md"), []byte("hello"), 0644)

	result, _ := runRead(vaultPath, "test.md")
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	var m map[string]any
	json.Unmarshal(data, &m)
	for _, key := range []string{"path", "properties", "body", "modified", "size"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing key %s in JSON", key)
		}
	}
}

// --- append/prepend ---

func TestAppendToExisting(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "note.md"), []byte("first"), 0644)

	err := runAppend(vaultPath, "note.md", "\nsecond")
	if err != nil {
		t.Fatalf("append failed: %v", err)
	}
	data, _ := os.ReadFile(filepath.Join(vaultPath, "note.md"))
	if string(data) != "first\nsecond" {
		t.Errorf("unexpected: %q", string(data))
	}
}

func TestAppendCreatesNonexistentFile(t *testing.T) {
	vaultPath := setupVault(t)

	err := runAppend(vaultPath, "new/deep/note.md", "created by append")
	if err != nil {
		t.Fatalf("append failed: %v", err)
	}
	data, _ := os.ReadFile(filepath.Join(vaultPath, "new", "deep", "note.md"))
	if string(data) != "created by append" {
		t.Errorf("unexpected: %q", string(data))
	}
}

func TestPrependWithFrontmatter(t *testing.T) {
	vaultPath := setupVault(t)
	content := "---\ntitle: Test\n---\n\nOriginal body."
	os.WriteFile(filepath.Join(vaultPath, "note.md"), []byte(content), 0644)

	err := runPrepend(vaultPath, "note.md", "## New Section\n\n")
	if err != nil {
		t.Fatalf("prepend failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(vaultPath, "note.md"))
	result := string(data)
	// Should have frontmatter, then new content, then original body
	if result[:3] != "---" {
		t.Error("expected frontmatter preserved at start")
	}
	if !contains(result, "## New Section") {
		t.Error("expected prepended content")
	}
	if !contains(result, "Original body.") {
		t.Error("expected original body preserved")
	}
}

func TestPrependWithoutFrontmatter(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "note.md"), []byte("existing content"), 0644)

	err := runPrepend(vaultPath, "note.md", "prepended\n")
	if err != nil {
		t.Fatalf("prepend failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(vaultPath, "note.md"))
	if string(data) != "prepended\nexisting content" {
		t.Errorf("unexpected: %q", string(data))
	}
}

// --- input flag (read content from file) ---

func TestResolveContentFromFlag(t *testing.T) {
	content, err := resolveContent("hello world", "")
	if err != nil {
		t.Fatalf("resolveContent failed: %v", err)
	}
	if content != "hello world" {
		t.Errorf("expected 'hello world', got %q", content)
	}
}

func TestResolveContentFromInputFile(t *testing.T) {
	tmp := t.TempDir()
	inputFile := filepath.Join(tmp, "source.md")
	os.WriteFile(inputFile, []byte("# From File\n\nLarge content here."), 0644)

	content, err := resolveContent("", inputFile)
	if err != nil {
		t.Fatalf("resolveContent failed: %v", err)
	}
	if content != "# From File\n\nLarge content here." {
		t.Errorf("unexpected content: %q", content)
	}
}

func TestResolveContentInputOverridesContent(t *testing.T) {
	tmp := t.TempDir()
	inputFile := filepath.Join(tmp, "source.md")
	os.WriteFile(inputFile, []byte("from file"), 0644)

	content, err := resolveContent("from flag", inputFile)
	if err != nil {
		t.Fatalf("resolveContent failed: %v", err)
	}
	if content != "from file" {
		t.Errorf("expected --input to take precedence, got %q", content)
	}
}

func TestResolveContentInputFileNotFound(t *testing.T) {
	_, err := resolveContent("", "/nonexistent/path/file.md")
	if err == nil {
		t.Error("expected error for nonexistent input file")
	}
}

func TestCreateFromInputFile(t *testing.T) {
	vaultPath := setupVault(t)
	tmp := t.TempDir()
	inputFile := filepath.Join(tmp, "large-doc.md")
	os.WriteFile(inputFile, []byte("# Large Document\n\nLots of content here."), 0644)

	content, _ := resolveContent("", inputFile)
	err := runCreate(vaultPath, "imported/doc", content, false)
	if err != nil {
		t.Fatalf("create from input failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(vaultPath, "imported", "doc.md"))
	if string(data) != "# Large Document\n\nLots of content here." {
		t.Errorf("unexpected content: %q", string(data))
	}
}

func TestAppendFromInputFile(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "note.md"), []byte("existing"), 0644)

	tmp := t.TempDir()
	inputFile := filepath.Join(tmp, "extra.md")
	os.WriteFile(inputFile, []byte("\nappended from file"), 0644)

	content, _ := resolveContent("", inputFile)
	err := runAppend(vaultPath, "note.md", content)
	if err != nil {
		t.Fatalf("append from input failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(vaultPath, "note.md"))
	if string(data) != "existing\nappended from file" {
		t.Errorf("unexpected content: %q", string(data))
	}
}

// --- move/delete/list/folders/mkdir ---

func TestMoveAutoTargetDirs(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "old.md"), []byte("content"), 0644)

	err := runMove(vaultPath, "old.md", "archive/2026/old.md")
	if err != nil {
		t.Fatalf("move failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(vaultPath, "old.md")); err == nil {
		t.Error("old file should not exist")
	}
	data, err := os.ReadFile(filepath.Join(vaultPath, "archive", "2026", "old.md"))
	if err != nil {
		t.Fatal("moved file not found")
	}
	if string(data) != "content" {
		t.Error("content changed during move")
	}
}

func TestDelete(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "delete-me.md"), []byte("bye"), 0644)

	err := runDelete(vaultPath, "delete-me.md")
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(vaultPath, "delete-me.md")); err == nil {
		t.Error("file should be deleted")
	}
}

func TestListAll(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "a.md"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(vaultPath, "Journal", "b.md"), []byte("b"), 0644)

	files, err := runList(vaultPath, "", "")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	// Should include at least README.md, a.md, b.md
	if len(files) < 3 {
		t.Errorf("expected at least 3 files, got %d", len(files))
	}
}

func TestListWithFolder(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "a.md"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(vaultPath, "Journal", "b.md"), []byte("b"), 0644)

	files, err := runList(vaultPath, "Journal", "")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(files) != 2 { // day1.md from setupVault + b.md
		t.Errorf("expected 2 files in Journal, got %d", len(files))
	}
}

func TestListWithExt(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "note.md"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(vaultPath, "data.txt"), []byte("b"), 0644)

	files, err := runList(vaultPath, "", "txt")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 txt file, got %d", len(files))
	}
}

func TestFoldersExcludesVaultfs(t *testing.T) {
	vaultPath := setupVault(t)

	folders, err := runFolders(vaultPath)
	if err != nil {
		t.Fatalf("folders failed: %v", err)
	}

	for _, f := range folders {
		if f == ".vaultfs" {
			t.Error(".vaultfs should be excluded")
		}
	}
	if len(folders) < 9 {
		t.Errorf("expected at least 9 folders (basic preset), got %d", len(folders))
	}
}

func TestMkdirRecursive(t *testing.T) {
	vaultPath := setupVault(t)

	err := runMkdir(vaultPath, "deep/nested/dir")
	if err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	if info, err := os.Stat(filepath.Join(vaultPath, "deep", "nested", "dir")); err != nil || !info.IsDir() {
		t.Error("expected deep/nested/dir to exist")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
