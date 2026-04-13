package markdown

import (
	"regexp"
	"strings"
)

// Task represents a markdown checkbox task.
type Task struct {
	Line     int      `json:"line"`
	Text     string   `json:"text"`
	Done     bool     `json:"done"`
	Priority string   `json:"priority,omitempty"`
	Due      string   `json:"due,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Mentions []string `json:"mentions,omitempty"`
}

var (
	taskRegex    = regexp.MustCompile(`^(\s*)-\s+\[([ xX])\]\s+(.*)$`)
	dueTagRegex  = regexp.MustCompile(`#due/(\d{4}-\d{2}-\d{2})`)
	dueEmojiRegex = regexp.MustCompile(`📅\s*(\d{4}-\d{2}-\d{2})`)
	mentionRegex = regexp.MustCompile(`@([a-zA-Z][a-zA-Z0-9_-]*)`)
	taskTagRegex = regexp.MustCompile(`(?:^|[ \t])#([a-zA-Z][a-zA-Z0-9_/-]*)`)
)

// Priority emoji mapping
var priorityMap = map[string]string{
	"🔴": "high",
	"⏫": "high",
	"🟡": "medium",
	"🔼": "medium",
	"🔵": "low",
	"🔽": "low",
}

// ExtractTasks extracts checkbox tasks from markdown content.
func ExtractTasks(data []byte) []Task {
	lines := strings.Split(string(data), "\n")
	var tasks []Task

	for i, line := range lines {
		matches := taskRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		done := matches[2] == "x" || matches[2] == "X"
		rawText := matches[3]

		task := Task{
			Line: i + 1, // 1-indexed
			Text: cleanTaskText(rawText),
			Done: done,
		}

		// Extract priority
		task.Priority = extractPriority(rawText)

		// Extract due date
		task.Due = extractDue(rawText)

		// Extract mentions
		task.Mentions = extractMentions(rawText)

		// Extract tags (excluding #due/...)
		task.Tags = extractTaskTags(rawText)

		tasks = append(tasks, task)
	}

	return tasks
}

func extractPriority(text string) string {
	for emoji, level := range priorityMap {
		if strings.Contains(text, emoji) {
			return level
		}
	}
	return ""
}

func extractDue(text string) string {
	if m := dueTagRegex.FindStringSubmatch(text); m != nil {
		return m[1]
	}
	if m := dueEmojiRegex.FindStringSubmatch(text); m != nil {
		return m[1]
	}
	return ""
}

func extractMentions(text string) []string {
	matches := mentionRegex.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return nil
	}
	var mentions []string
	for _, m := range matches {
		mentions = append(mentions, m[1])
	}
	return mentions
}

func extractTaskTags(text string) []string {
	matches := taskTagRegex.FindAllStringSubmatch(" "+text, -1)
	if len(matches) == 0 {
		return nil
	}
	var tags []string
	for _, m := range matches {
		tag := m[1]
		// Exclude #due/... patterns
		if strings.HasPrefix(tag, "due/") {
			continue
		}
		tags = append(tags, tag)
	}
	if len(tags) == 0 {
		return nil
	}
	return tags
}

// cleanTaskText strips priority emojis and metadata from display text.
func cleanTaskText(text string) string {
	// Remove priority emojis
	for emoji := range priorityMap {
		text = strings.ReplaceAll(text, emoji, "")
	}
	// Trim whitespace
	text = strings.TrimSpace(text)
	return text
}
