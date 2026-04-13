package markdown

import (
	"regexp"
	"strings"
)

// Heading represents a heading in the outline tree.
type Heading struct {
	Level    int        `json:"level"`
	Text     string     `json:"text"`
	Children []*Heading `json:"children,omitempty"`
}

var headingRegex = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

// ExtractOutline parses markdown and returns a nested heading tree.
func ExtractOutline(data []byte) []*Heading {
	lines := strings.Split(string(data), "\n")

	var allHeadings []Heading
	for _, line := range lines {
		m := headingRegex.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		allHeadings = append(allHeadings, Heading{
			Level: len(m[1]),
			Text:  strings.TrimSpace(m[2]),
		})
	}

	if len(allHeadings) == 0 {
		return nil
	}

	return buildTree(allHeadings)
}

// buildTree converts a flat list of headings into a nested tree.
func buildTree(headings []Heading) []*Heading {
	var roots []*Heading
	var stack []*Heading

	for i := range headings {
		h := &headings[i]

		// Pop stack until we find a parent (lower level)
		for len(stack) > 0 && stack[len(stack)-1].Level >= h.Level {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			roots = append(roots, h)
		} else {
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, h)
		}

		stack = append(stack, h)
	}

	return roots
}
