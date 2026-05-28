// Package vfs holds shared vault primitives, including the not-found
// signal used so AI agents can distinguish "file does not exist" from
// real errors.
package vfs

import "io/fs"

// NotFoundError signals that a requested path does not exist. It is
// not a fatal error — callers convert it to a structured "not found"
// response so the agent driving the CLI can branch on the absence of
// a file without parsing stderr or special-casing exit codes.
type NotFoundError struct {
	Path string
}

func (e *NotFoundError) Error() string {
	return "file not found: " + e.Path
}

// Is lets errors.Is(err, fs.ErrNotExist) match a NotFoundError, so the
// value plays nicely with stdlib not-exist checks anywhere else.
func (e *NotFoundError) Is(target error) bool {
	return target == fs.ErrNotExist
}

func NewNotFound(path string) error {
	return &NotFoundError{Path: path}
}
