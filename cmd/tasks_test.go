package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTasksAll(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "todo.md"), []byte("- [ ] Buy milk\n- [x] Clean house\n- [ ] 🔴 Urgent task #due/2026-04-15 @alice #backend"), 0644)

	tasks, err := runTasks(vaultPath, "", "")
	if err != nil {
		t.Fatalf("tasks failed: %v", err)
	}

	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}
}

func TestTasksPending(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "todo.md"), []byte("- [ ] Pending\n- [x] Done"), 0644)

	tasks, err := runTasks(vaultPath, "pending", "")
	if err != nil {
		t.Fatalf("tasks failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 pending task, got %d", len(tasks))
	}
	if tasks[0].Done {
		t.Error("expected pending task")
	}
}

func TestTasksDone(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "todo.md"), []byte("- [ ] Pending\n- [x] Done"), 0644)

	tasks, err := runTasks(vaultPath, "done", "")
	if err != nil {
		t.Fatalf("tasks failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 done task, got %d", len(tasks))
	}
}

func TestTasksFolder(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "todo.md"), []byte("- [ ] Root task"), 0644)
	os.WriteFile(filepath.Join(vaultPath, "Journal", "day.md"), []byte("- [ ] Journal task"), 0644)

	tasks, err := runTasks(vaultPath, "", "Journal")
	if err != nil {
		t.Fatalf("tasks failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task in Journal, got %d", len(tasks))
	}
}

func TestTasksMetadata(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "todo.md"), []byte("- [ ] 🔴 Fix bug #due/2026-04-15 @alice #backend"), 0644)

	tasks, err := runTasks(vaultPath, "", "")
	if err != nil {
		t.Fatalf("tasks failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	task := tasks[0]
	if task.Priority != "high" {
		t.Errorf("expected high priority, got %s", task.Priority)
	}
	if task.Due != "2026-04-15" {
		t.Errorf("expected due 2026-04-15, got %s", task.Due)
	}
	if len(task.Mentions) != 1 || task.Mentions[0] != "alice" {
		t.Errorf("expected mention alice, got %v", task.Mentions)
	}
	if len(task.Tags) != 1 || task.Tags[0] != "backend" {
		t.Errorf("expected tag backend, got %v", task.Tags)
	}
}

func TestTaskToggle(t *testing.T) {
	vaultPath := setupVault(t)
	os.WriteFile(filepath.Join(vaultPath, "todo.md"), []byte("# Todo\n- [ ] Buy milk\n- [x] Done"), 0644)

	err := runTaskToggle(vaultPath, "todo.md", 2)
	if err != nil {
		t.Fatalf("toggle failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(vaultPath, "todo.md"))
	lines := splitLines(string(data))
	if lines[1] != "- [x] Buy milk" {
		t.Errorf("expected toggled to done, got %q", lines[1])
	}

	// Toggle back
	err = runTaskToggle(vaultPath, "todo.md", 2)
	if err != nil {
		t.Fatalf("toggle back failed: %v", err)
	}
	data, _ = os.ReadFile(filepath.Join(vaultPath, "todo.md"))
	lines = splitLines(string(data))
	if lines[1] != "- [ ] Buy milk" {
		t.Errorf("expected toggled back, got %q", lines[1])
	}
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
