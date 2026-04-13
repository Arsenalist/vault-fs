package markdown

import (
	"testing"
)

func TestExtractTasksBasic(t *testing.T) {
	input := `# Todo

- [ ] Buy groceries
- [x] Clean house
- Regular list item`

	tasks := ExtractTasks([]byte(input))

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].Text != "Buy groceries" {
		t.Errorf("expected 'Buy groceries', got %q", tasks[0].Text)
	}
	if tasks[0].Done {
		t.Error("expected first task to be not done")
	}
	if tasks[0].Line != 3 {
		t.Errorf("expected line 3, got %d", tasks[0].Line)
	}
	if !tasks[1].Done {
		t.Error("expected second task to be done")
	}
}

func TestExtractTasksPriorityHigh(t *testing.T) {
	input := `- [ ] 🔴 Critical bug fix
- [ ] ⏫ Also important`

	tasks := ExtractTasks([]byte(input))

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].Priority != "high" {
		t.Errorf("expected high priority for 🔴, got %q", tasks[0].Priority)
	}
	if tasks[1].Priority != "high" {
		t.Errorf("expected high priority for ⏫, got %q", tasks[1].Priority)
	}
}

func TestExtractTasksPriorityMedium(t *testing.T) {
	input := `- [ ] 🟡 Medium task
- [ ] 🔼 Also medium`

	tasks := ExtractTasks([]byte(input))

	for i, task := range tasks {
		if task.Priority != "medium" {
			t.Errorf("task %d: expected medium priority, got %q", i, task.Priority)
		}
	}
}

func TestExtractTasksPriorityLow(t *testing.T) {
	input := `- [ ] 🔵 Low task
- [ ] 🔽 Also low`

	tasks := ExtractTasks([]byte(input))

	for i, task := range tasks {
		if task.Priority != "low" {
			t.Errorf("task %d: expected low priority, got %q", i, task.Priority)
		}
	}
}

func TestExtractTasksDueDate(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"- [ ] Task #due/2026-04-15", "2026-04-15"},
		{"- [ ] Task 📅 2026-04-15", "2026-04-15"},
	}

	for _, tt := range tests {
		tasks := ExtractTasks([]byte(tt.input))
		if len(tasks) != 1 {
			t.Fatalf("expected 1 task, got %d", len(tasks))
		}
		if tasks[0].Due != tt.expected {
			t.Errorf("input %q: expected due=%s, got %s", tt.input, tt.expected, tasks[0].Due)
		}
	}
}

func TestExtractTasksMentions(t *testing.T) {
	input := `- [ ] Review PR with @alice and @bob`

	tasks := ExtractTasks([]byte(input))
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if len(tasks[0].Mentions) != 2 {
		t.Fatalf("expected 2 mentions, got %d", len(tasks[0].Mentions))
	}
	if tasks[0].Mentions[0] != "alice" || tasks[0].Mentions[1] != "bob" {
		t.Errorf("unexpected mentions: %v", tasks[0].Mentions)
	}
}

func TestExtractTasksInlineTags(t *testing.T) {
	input := `- [ ] Fix bug #urgent #backend`

	tasks := ExtractTasks([]byte(input))
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if len(tasks[0].Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d: %v", len(tasks[0].Tags), tasks[0].Tags)
	}
}

func TestExtractTasksNoPriority(t *testing.T) {
	input := `- [ ] Plain task`

	tasks := ExtractTasks([]byte(input))
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Priority != "" {
		t.Errorf("expected empty priority, got %q", tasks[0].Priority)
	}
	if tasks[0].Due != "" {
		t.Errorf("expected empty due, got %q", tasks[0].Due)
	}
}

func TestExtractTasksDueTagNotInTags(t *testing.T) {
	input := `- [ ] Task #due/2026-04-15 #urgent`

	tasks := ExtractTasks([]byte(input))
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	// #due/... should be parsed as due date, not included in tags
	for _, tag := range tasks[0].Tags {
		if tag == "due/2026-04-15" {
			t.Error("due date should not appear in tags list")
		}
	}
	if tasks[0].Due != "2026-04-15" {
		t.Errorf("expected due=2026-04-15, got %s", tasks[0].Due)
	}
	if len(tasks[0].Tags) != 1 || tasks[0].Tags[0] != "urgent" {
		t.Errorf("expected [urgent], got %v", tasks[0].Tags)
	}
}
