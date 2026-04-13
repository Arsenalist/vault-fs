package markdown

import (
	"bytes"

	"github.com/goccy/go-yaml"
)

var fmDelimiter = []byte("---")

// ParseFrontmatter splits a markdown file into YAML frontmatter (as a map) and body.
// Frontmatter is only recognized at byte 0 of the file, delimited by "---".
// Returns empty map and full content if no valid frontmatter is found.
func ParseFrontmatter(data []byte) (map[string]any, []byte, error) {
	// Must start with ---
	if !bytes.HasPrefix(data, fmDelimiter) {
		return map[string]any{}, data, nil
	}

	// Find closing ---
	rest := data[3:]
	// Skip the newline after opening ---
	if len(rest) > 0 && rest[0] == '\n' {
		rest = rest[1:]
	} else if len(rest) > 1 && rest[0] == '\r' && rest[1] == '\n' {
		rest = rest[2:]
	}

	idx := bytes.Index(rest, fmDelimiter)
	if idx < 0 {
		// No closing delimiter — treat as no frontmatter
		return map[string]any{}, data, nil
	}

	// Ensure the closing --- is at the start of a line
	if idx > 0 && rest[idx-1] != '\n' {
		return map[string]any{}, data, nil
	}

	yamlContent := rest[:idx]
	body := bytes.TrimLeft(rest[idx+3:], "\r\n") // skip closing --- and leading whitespace

	var fm map[string]any
	if err := yaml.Unmarshal(yamlContent, &fm); err != nil {
		// Malformed YAML — treat as no frontmatter
		return map[string]any{}, data, nil
	}
	if fm == nil {
		fm = map[string]any{}
	}

	return fm, body, nil
}
