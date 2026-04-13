package markdown

import (
	"sort"
	"testing"
)

func TestExtractInlineTags(t *testing.T) {
	input := `# Meeting Notes

Discussion about #project and #deadline.
Also mentioned #team/backend work.`

	tags := ExtractInlineTags([]byte(input))

	expected := []string{"deadline", "project", "team/backend"}
	sort.Strings(tags)
	if len(tags) != len(expected) {
		t.Fatalf("expected %d tags, got %d: %v", len(expected), len(tags), tags)
	}
	for i, tag := range tags {
		if tag != expected[i] {
			t.Errorf("expected tag %s, got %s", expected[i], tag)
		}
	}
}

func TestExtractInlineTagsIgnoresCodeBlocks(t *testing.T) {
	input := "Some text #real-tag\n\n```\n#not-a-tag\n```\n\nMore text `#also-not-a-tag`"

	tags := ExtractInlineTags([]byte(input))

	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d: %v", len(tags), tags)
	}
	if tags[0] != "real-tag" {
		t.Errorf("expected 'real-tag', got %s", tags[0])
	}
}

func TestExtractInlineTagsIgnoresHeadings(t *testing.T) {
	input := "# Heading\n## Another Heading\nText with #actual-tag"

	tags := ExtractInlineTags([]byte(input))

	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d: %v", len(tags), tags)
	}
	if tags[0] != "actual-tag" {
		t.Errorf("expected 'actual-tag', got %s", tags[0])
	}
}

func TestExtractAllTags(t *testing.T) {
	input := `---
tags: [work, urgent]
---

Some text #work #personal and #project/alpha.`

	tags, err := ExtractAllTags([]byte(input))
	if err != nil {
		t.Fatalf("ExtractAllTags failed: %v", err)
	}

	sort.Strings(tags)
	expected := []string{"personal", "project/alpha", "urgent", "work"}
	if len(tags) != len(expected) {
		t.Fatalf("expected %d tags, got %d: %v", len(expected), len(tags), tags)
	}
	for i, tag := range tags {
		if tag != expected[i] {
			t.Errorf("expected tag %s, got %s", expected[i], tag)
		}
	}
}

func TestExtractAllTagsDeduplication(t *testing.T) {
	input := `---
tags: [urgent]
---

This is #urgent and important.`

	tags, err := ExtractAllTags([]byte(input))
	if err != nil {
		t.Fatalf("ExtractAllTags failed: %v", err)
	}

	count := 0
	for _, tag := range tags {
		if tag == "urgent" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 'urgent' once, found %d times", count)
	}
}

func TestExtractAllTagsNoTags(t *testing.T) {
	input := `# Plain Markdown

No tags here at all.`

	tags, err := ExtractAllTags([]byte(input))
	if err != nil {
		t.Fatalf("ExtractAllTags failed: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected no tags, got %v", tags)
	}
}
