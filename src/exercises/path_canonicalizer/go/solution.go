package main

import (
	"errors"
	"path"
	"strings"
)

var ErrPathEscapesRoot = errors.New("path escapes root")

// CanonicalizePath normalizes rawPath to a clean POSIX path starting with '/'.
func CanonicalizePath(rawPath string) (string, error) {
	posix := strings.ReplaceAll(rawPath, "\\", "/")

	rel := path.Clean(posix)
	if rel == ".." || strings.HasPrefix(rel, "../") {
		return "", ErrPathEscapesRoot
	}

	if strings.HasPrefix(posix, "/..") {
		return "", ErrPathEscapesRoot
	}

	return path.Clean("/" + posix), nil
}
