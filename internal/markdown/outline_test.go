package markdown

import (
	"testing"
)

func TestExtractOutlineNested(t *testing.T) {
	input := `# Title

## Section A

### Sub A1

## Section B`

	outline := ExtractOutline([]byte(input))

	if len(outline) != 1 {
		t.Fatalf("expected 1 root heading, got %d", len(outline))
	}
	if outline[0].Text != "Title" {
		t.Errorf("expected 'Title', got %s", outline[0].Text)
	}
	if outline[0].Level != 1 {
		t.Errorf("expected level 1, got %d", outline[0].Level)
	}
	if len(outline[0].Children) != 2 {
		t.Fatalf("expected 2 children of Title, got %d", len(outline[0].Children))
	}
	sectionA := outline[0].Children[0]
	if sectionA.Text != "Section A" {
		t.Errorf("expected 'Section A', got %s", sectionA.Text)
	}
	if len(sectionA.Children) != 1 {
		t.Fatalf("expected 1 child of Section A, got %d", len(sectionA.Children))
	}
	if sectionA.Children[0].Text != "Sub A1" {
		t.Errorf("expected 'Sub A1', got %s", sectionA.Children[0].Text)
	}
}

func TestExtractOutlineFlat(t *testing.T) {
	input := `## A
## B
## C`

	outline := ExtractOutline([]byte(input))

	if len(outline) != 3 {
		t.Fatalf("expected 3 headings at same level, got %d", len(outline))
	}
}

func TestExtractOutlineNoHeadings(t *testing.T) {
	input := `Just some text without headings.`

	outline := ExtractOutline([]byte(input))

	if len(outline) != 0 {
		t.Errorf("expected empty outline, got %d headings", len(outline))
	}
}

func TestExtractOutlineSkippedLevels(t *testing.T) {
	input := `# Title

### Skipped to H3

## Back to H2`

	outline := ExtractOutline([]byte(input))

	if len(outline) != 1 {
		t.Fatalf("expected 1 root, got %d", len(outline))
	}
	// H3 should still be a child of H1 even though H2 was skipped
	if len(outline[0].Children) != 2 {
		t.Fatalf("expected 2 children of Title, got %d", len(outline[0].Children))
	}
}
