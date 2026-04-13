package markdown

import (
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	input := `---
title: Test Note
tags: [work, urgent]
---

# Body Content

Some text here.`

	fm, body, err := ParseFrontmatter([]byte(input))
	if err != nil {
		t.Fatalf("ParseFrontmatter failed: %v", err)
	}

	if fm["title"] != "Test Note" {
		t.Errorf("expected title='Test Note', got %v", fm["title"])
	}

	tags, ok := fm["tags"].([]any)
	if !ok {
		t.Fatalf("expected tags to be a slice, got %T", fm["tags"])
	}
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}

	expectedBody := "# Body Content\n\nSome text here."
	if string(body) != expectedBody {
		t.Errorf("unexpected body: %q", string(body))
	}
}

func TestParseFrontmatterEmpty(t *testing.T) {
	input := `---
---

Body only.`

	fm, body, err := ParseFrontmatter([]byte(input))
	if err != nil {
		t.Fatalf("ParseFrontmatter failed: %v", err)
	}
	if len(fm) != 0 {
		t.Errorf("expected empty frontmatter, got %v", fm)
	}
	if string(body) != "Body only." {
		t.Errorf("unexpected body: %q", string(body))
	}
}

func TestParseFrontmatterNone(t *testing.T) {
	input := `# Just Markdown

No frontmatter here.`

	fm, body, err := ParseFrontmatter([]byte(input))
	if err != nil {
		t.Fatalf("ParseFrontmatter failed: %v", err)
	}
	if len(fm) != 0 {
		t.Errorf("expected empty frontmatter, got %v", fm)
	}
	if string(body) != input {
		t.Errorf("expected full content as body, got %q", string(body))
	}
}

func TestParseFrontmatterMalformed(t *testing.T) {
	input := `---
title: unclosed
no closing delimiter

# Body`

	fm, body, err := ParseFrontmatter([]byte(input))
	if err != nil {
		t.Fatalf("ParseFrontmatter failed: %v", err)
	}
	// Malformed = treat as no frontmatter
	if len(fm) != 0 {
		t.Errorf("expected empty frontmatter for malformed input, got %v", fm)
	}
	if string(body) != input {
		t.Errorf("expected full content as body")
	}
}

func TestParseFrontmatterNotAtStart(t *testing.T) {
	input := `Some text before
---
title: Not Frontmatter
---`

	fm, body, err := ParseFrontmatter([]byte(input))
	if err != nil {
		t.Fatalf("ParseFrontmatter failed: %v", err)
	}
	// --- not at byte 0, so not frontmatter
	if len(fm) != 0 {
		t.Errorf("expected empty frontmatter, got %v", fm)
	}
	if string(body) != input {
		t.Errorf("expected full content as body")
	}
}
