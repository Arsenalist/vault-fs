package markdown

import (
	"regexp"
	"strings"
)

// tagRegex matches inline #tags including nested tags like #project/alpha.
// Must be preceded by whitespace or start of line, not inside a word.
var tagRegex = regexp.MustCompile(`(?:^|[ \t])#([a-zA-Z][a-zA-Z0-9_/-]*)`)

// ExtractInlineTags extracts #tags from markdown body text,
// ignoring code blocks (fenced and inline) and headings.
func ExtractInlineTags(body []byte) []string {
	lines := strings.Split(string(body), "\n")
	seen := make(map[string]bool)
	var tags []string

	inCodeBlock := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Toggle fenced code blocks
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			continue
		}

		// Skip headings (lines starting with #)
		if strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "# ") || strings.HasPrefix(trimmed, "## ") || strings.HasPrefix(trimmed, "### ") || strings.HasPrefix(trimmed, "#### ") || strings.HasPrefix(trimmed, "##### ") || strings.HasPrefix(trimmed, "###### ") {
			// Still extract tags from heading content
			// Actually headings start with # followed by space — the # in heading syntax
			// should not be confused with tags. Let's only look for tags after heading prefix.
		}

		// Remove inline code
		cleaned := removeInlineCode(line)

		// Remove heading prefix so "# Heading" doesn't match
		if strings.HasPrefix(trimmed, "#") {
			// Check if it's a heading (# followed by space or alone)
			for i := 0; i < len(trimmed); i++ {
				if trimmed[i] != '#' {
					if trimmed[i] == ' ' {
						// It's a heading, strip the prefix
						cleaned = trimmed[i+1:]
					}
					break
				}
			}
		}

		matches := tagRegex.FindAllStringSubmatch(" "+cleaned, -1)
		for _, m := range matches {
			tag := m[1]
			if !seen[tag] {
				seen[tag] = true
				tags = append(tags, tag)
			}
		}
	}

	return tags
}

// removeInlineCode strips `inline code` from a line.
func removeInlineCode(line string) string {
	result := strings.Builder{}
	inCode := false
	for i := 0; i < len(line); i++ {
		if line[i] == '`' {
			inCode = !inCode
			continue
		}
		if !inCode {
			result.WriteByte(line[i])
		}
	}
	return result.String()
}

// ExtractAllTags extracts tags from both YAML frontmatter and inline body,
// returning a deduplicated list.
func ExtractAllTags(data []byte) ([]string, error) {
	fm, body, err := ParseFrontmatter(data)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var tags []string

	// Extract from frontmatter
	if fmTags, ok := fm["tags"]; ok {
		switch v := fmTags.(type) {
		case []any:
			for _, t := range v {
				if s, ok := t.(string); ok && !seen[s] {
					seen[s] = true
					tags = append(tags, s)
				}
			}
		case string:
			s := strings.TrimSpace(v)
			if s != "" && !seen[s] {
				seen[s] = true
				tags = append(tags, s)
			}
		}
	}

	// Extract inline tags from body
	for _, tag := range ExtractInlineTags(body) {
		if !seen[tag] {
			seen[tag] = true
			tags = append(tags, tag)
		}
	}

	return tags, nil
}
