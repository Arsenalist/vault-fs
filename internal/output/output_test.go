package output

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestJSONFormat(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]string{"path": "notes/test.md", "body": "hello"}

	err := WriteJSON(&buf, data)
	if err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result["path"] != "notes/test.md" {
		t.Errorf("expected path=notes/test.md, got %s", result["path"])
	}
}

func TestJSONFormatIndented(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]int{"count": 42}

	err := WriteJSON(&buf, data)
	if err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	// Output should be indented for readability
	output := buf.String()
	if output[0] != '{' {
		t.Errorf("expected indented JSON starting with '{', got: %s", output[:20])
	}
	if !bytes.Contains(buf.Bytes(), []byte("\n")) {
		t.Error("expected indented JSON with newlines")
	}
}

func TestTextFormat(t *testing.T) {
	var buf bytes.Buffer
	WriteText(&buf, "Vault initialized at /tmp/vault (preset: basic, 9 directories created)")

	output := buf.String()
	if output != "Vault initialized at /tmp/vault (preset: basic, 9 directories created)\n" {
		t.Errorf("unexpected text output: %q", output)
	}
}

func TestErrorJSON(t *testing.T) {
	var buf bytes.Buffer
	WriteErrorJSON(&buf, "file not found: notes/missing.md")

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("error output is not valid JSON: %v", err)
	}
	if result["error"] != "file not found: notes/missing.md" {
		t.Errorf("expected error message, got %s", result["error"])
	}
}

func TestErrorText(t *testing.T) {
	var buf bytes.Buffer
	WriteErrorText(&buf, "file not found: notes/missing.md")

	output := buf.String()
	expected := "Error: file not found: notes/missing.md\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestDefaultFormatQuery(t *testing.T) {
	f := ResolveFormat("", true)
	if f != FormatJSON {
		t.Errorf("expected JSON for query command, got %s", f)
	}
}

func TestDefaultFormatAction(t *testing.T) {
	f := ResolveFormat("", false)
	if f != FormatText {
		t.Errorf("expected Text for action command, got %s", f)
	}
}

func TestExplicitFormatOverride(t *testing.T) {
	f := ResolveFormat("text", true)
	if f != FormatText {
		t.Errorf("expected Text override for query command, got %s", f)
	}

	f = ResolveFormat("json", false)
	if f != FormatJSON {
		t.Errorf("expected JSON override for action command, got %s", f)
	}
}
