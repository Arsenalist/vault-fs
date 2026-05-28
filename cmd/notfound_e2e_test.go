package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// execRoot invokes the cobra root with given args, capturing stdout/stderr
// output written via cmd.OutOrStdout()/cmd.ErrOrStderr(). It returns the
// stdout buffer and the error returned by cobra (nil means RunE returned
// nil → cobra would exit 0).
func execRoot(t *testing.T, args ...string) (stdout string, err error) {
	t.Helper()
	var outBuf, errBuf bytes.Buffer
	rootCmd.SetOut(&outBuf)
	rootCmd.SetErr(&errBuf)
	rootCmd.SetArgs(args)
	t.Cleanup(func() {
		// Reset global flag state so other tests aren't affected.
		formatFlag = ""
		vaultFlag = ""
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
		rootCmd.SetArgs(nil)
	})
	err = rootCmd.Execute()
	return outBuf.String(), err
}

func TestE2EReadMissingFileJSON(t *testing.T) {
	vault := setupVault(t)

	out, err := execRoot(t, "read", "missing.md", "--vault", vault, "--format", "json")
	if err != nil {
		t.Fatalf("expected nil error (exit 0), got: %v", err)
	}

	var result map[string]any
	if jerr := json.Unmarshal([]byte(out), &result); jerr != nil {
		t.Fatalf("expected JSON output on stdout, got %q: %v", out, jerr)
	}
	if found, ok := result["found"].(bool); !ok || found {
		t.Errorf("expected found=false, got %v", result["found"])
	}
	if result["path"] != "missing.md" {
		t.Errorf("expected path=missing.md, got %v", result["path"])
	}
}

func TestE2EMoveMissingSourceText(t *testing.T) {
	vault := setupVault(t)

	out, err := execRoot(t, "move", "ghost.md", "--to", "dest.md", "--vault", vault)
	if err != nil {
		t.Fatalf("expected nil error (exit 0), got: %v", err)
	}
	if out != "not found: ghost.md\n" {
		t.Errorf("expected 'not found: ghost.md\\n', got %q", out)
	}
}

func TestE2EDeleteMissingExitsZero(t *testing.T) {
	vault := setupVault(t)

	_, err := execRoot(t, "delete", "never-existed.md", "--vault", vault)
	if err != nil {
		t.Fatalf("delete of missing file should exit 0, got: %v", err)
	}
}

func TestE2EPropertyRemoveNotFoundJSON(t *testing.T) {
	vault := setupVault(t)

	out, err := execRoot(t, "property", "remove", "missing.md", "--name", "tag", "--vault", vault, "--format", "json")
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	var result map[string]any
	if jerr := json.Unmarshal([]byte(out), &result); jerr != nil {
		t.Fatalf("expected JSON output, got %q: %v", out, jerr)
	}
	if result["path"] != "missing.md" {
		t.Errorf("expected path=missing.md, got %v", result["path"])
	}
}

func TestE2ETaskToggleNotFoundText(t *testing.T) {
	vault := setupVault(t)

	out, err := execRoot(t, "task", "toggle", "missing.md", "--line", "1", "--vault", vault)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if !strings.Contains(out, "not found: missing.md") {
		t.Errorf("expected text not-found message, got %q", out)
	}
}

func TestE2EOutlineNotFoundJSON(t *testing.T) {
	vault := setupVault(t)

	out, err := execRoot(t, "outline", "missing.md", "--vault", vault)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	var result map[string]any
	if jerr := json.Unmarshal([]byte(out), &result); jerr != nil {
		t.Fatalf("expected JSON output (query default), got %q: %v", out, jerr)
	}
	if found, ok := result["found"].(bool); !ok || found {
		t.Errorf("expected found=false, got %v", result["found"])
	}
}
