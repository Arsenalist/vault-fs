package vfs

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"
)

func TestNotFoundErrorMessage(t *testing.T) {
	err := NewNotFound("notes/missing.md")
	want := "file not found: notes/missing.md"
	if err.Error() != want {
		t.Errorf("expected %q, got %q", want, err.Error())
	}
}

func TestNotFoundErrorIsFsErrNotExist(t *testing.T) {
	err := NewNotFound("x.md")
	if !errors.Is(err, fs.ErrNotExist) {
		t.Error("expected errors.Is(err, fs.ErrNotExist) to be true")
	}
}

func TestNotFoundErrorAsExtractsPath(t *testing.T) {
	wrapped := fmt.Errorf("doing thing: %w", NewNotFound("a/b.md"))

	var nf *NotFoundError
	if !errors.As(wrapped, &nf) {
		t.Fatal("expected errors.As to extract *NotFoundError")
	}
	if nf.Path != "a/b.md" {
		t.Errorf("expected Path=a/b.md, got %q", nf.Path)
	}
}
