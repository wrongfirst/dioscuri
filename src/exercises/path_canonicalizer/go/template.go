package main

import (
	"errors"
	"path"
	"strings"
)

var ErrPathEscapesRoot = errors.New("path escapes root")

// CanonicalizePath normalizes rawPath to a clean POSIX path starting with '/'.
// It returns ErrPathEscapesRoot if the path attempts to traverse above the root.
func CanonicalizePath(rawPath string) (string, error) {
	return "", nil
}
