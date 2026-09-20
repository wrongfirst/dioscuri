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
	posix := strings.ReplaceAll(rawPath, "\\", "/")

	// Strip Windows drive letter if present (e.g. "C:")
	if len(posix) >= 2 && posix[1] == ':' {
		ch := posix[0]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			posix = posix[2:]
		}
	}

	trimmed := strings.TrimLeft(posix, "/")
	cleanRel := path.Clean(trimmed)

	if cleanRel == ".." || strings.HasPrefix(cleanRel, "../") {
		return "", ErrPathEscapesRoot
	}

	if cleanRel == "." || cleanRel == "" {
		return "/", nil
	}

	return "/" + cleanRel, nil
}
