package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestE2EFullWorkflow(t *testing.T) {
	tmp := t.TempDir()
	vaultPath := filepath.Join(tmp, "vault")

	// Init
	err := runInit(vaultPath, "basic", nil)
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Create files
	err = runCreate(vaultPath, "notes/standup", "---\ntags: [work, daily]\n---\n\n# Standup\n\n- [ ] 🔴 Fix auth bug #due/2026-04-15 @alice #backend\n- [x] Deploy v2\n- [ ] Review PR #frontend", false)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	err = runCreate(vaultPath, "projects/alpha", "---\ntags: [project]\n---\n\n# Project Alpha\n\n## Goals\n\n### Q1 Targets\n\nBudget approved.", false)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Read
	result, err := runRead(vaultPath, "notes/standup.md")
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if result.Properties["tags"] == nil {
		t.Error("expected tags in properties")
	}
	if result.Body == "" {
		t.Error("expected non-empty body")
	}

	// Append
	err = runAppend(vaultPath, "notes/standup.md", "\n- [ ] New item")
	if err != nil {
		t.Fatalf("append failed: %v", err)
	}

	// Tags
	tags, err := runTags(vaultPath, true, "count")
	if err != nil {
		t.Fatalf("tags failed: %v", err)
	}
	if len(tags) == 0 {
		t.Error("expected tags")
	}
	// "work" should be in tags from frontmatter
	found := false
	for _, tag := range tags {
		if tag.Name == "work" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'work' tag")
	}

	// Tasks
	tasks, err := runTasks(vaultPath, "pending", "")
	if err != nil {
		t.Fatalf("tasks failed: %v", err)
	}
	if len(tasks) < 2 {
		t.Errorf("expected at least 2 pending tasks, got %d", len(tasks))
	}
	// First task should have metadata
	for _, task := range tasks {
		if task.Priority == "high" && task.Due == "2026-04-15" {
			if len(task.Mentions) == 0 || task.Mentions[0] != "alice" {
				t.Error("expected @alice mention")
			}
			break
		}
	}

	// Properties
	props, err := runProperties(vaultPath, "notes/standup.md")
	if err != nil {
		t.Fatalf("properties failed: %v", err)
	}
	if props["tags"] == nil {
		t.Error("expected tags property")
	}

	err = runPropertySet(vaultPath, "notes/standup.md", "status", "active")
	if err != nil {
		t.Fatalf("property set failed: %v", err)
	}
	props, _ = runProperties(vaultPath, "notes/standup.md")
	if props["status"] != "active" {
		t.Errorf("expected status=active, got %v", props["status"])
	}

	// Outline
	outline, err := runOutline(vaultPath, "projects/alpha.md")
	if err != nil {
		t.Fatalf("outline failed: %v", err)
	}
	if len(outline) == 0 {
		t.Error("expected outline headings")
	}

	// Search
	searchResults, err := runSearch(vaultPath, "budget", "", 10, false, false)
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(searchResults) == 0 {
		t.Error("expected search results for 'budget'")
	}

	// Search context
	contextResults, err := runSearchContext(vaultPath, "auth", 10)
	if err != nil {
		t.Fatalf("search context failed: %v", err)
	}
	if len(contextResults) == 0 {
		t.Error("expected context results for 'auth'")
	}

	// Recent
	recent, err := runRecent(vaultPath, 7, 20, "")
	if err != nil {
		t.Fatalf("recent failed: %v", err)
	}
	if len(recent) < 2 {
		t.Errorf("expected at least 2 recent files, got %d", len(recent))
	}

	// Info
	info, err := getVaultInfo(vaultPath)
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}
	if info.FileCount < 3 {
		t.Errorf("expected at least 3 files, got %d", info.FileCount)
	}

	// Move
	err = runMove(vaultPath, "notes/standup.md", "archive/standup.md")
	if err != nil {
		t.Fatalf("move failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(vaultPath, "archive", "standup.md")); err != nil {
		t.Error("expected moved file")
	}

	// Tag by name
	tagFiles, err := runTagByName(vaultPath, "project")
	if err != nil {
		t.Fatalf("tag by name failed: %v", err)
	}
	if len(tagFiles) == 0 {
		t.Error("expected files with 'project' tag")
	}

	// Task toggle
	err = runCreate(vaultPath, "toggle-test", "- [ ] Test task\n- [x] Done task", false)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	err = runTaskToggle(vaultPath, "toggle-test.md", 1)
	if err != nil {
		t.Fatalf("toggle failed: %v", err)
	}
	data, _ := os.ReadFile(filepath.Join(vaultPath, "toggle-test.md"))
	if string(data)[:5] != "- [x]" {
		t.Error("expected task toggled to done")
	}
}
