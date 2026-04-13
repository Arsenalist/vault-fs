package output

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	FormatJSON = "json"
	FormatText = "text"
)

// ResolveFormat determines the output format. If explicit is set, use it.
// Otherwise, query commands default to JSON, action commands default to text.
func ResolveFormat(explicit string, isQuery bool) string {
	if explicit == FormatJSON || explicit == FormatText {
		return explicit
	}
	if isQuery {
		return FormatJSON
	}
	return FormatText
}

// WriteJSON writes data as indented JSON to the writer.
func WriteJSON(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// WriteText writes a plain text message with a trailing newline.
func WriteText(w io.Writer, msg string) {
	fmt.Fprintln(w, msg)
}

// WriteErrorJSON writes an error as a JSON object: {"error": "message"}.
func WriteErrorJSON(w io.Writer, msg string) {
	WriteJSON(w, map[string]string{"error": msg})
}

// WriteErrorText writes an error as "Error: message\n".
func WriteErrorText(w io.Writer, msg string) {
	fmt.Fprintf(w, "Error: %s\n", msg)
}
