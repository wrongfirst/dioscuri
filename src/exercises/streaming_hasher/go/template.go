package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

type FS interface {
	Open(path string) (io.ReadCloser, error)
}

type MockFS struct {
	Files map[string][]byte
}

func NewMockFS() *MockFS {
	return &MockFS{Files: make(map[string][]byte)}
}

func (m *MockFS) Open(path string) (io.ReadCloser, error) {
	data, ok := m.Files[path]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

// HashFile opens path from fs, streams it through SHA-256, and returns the hex digest.
func HashFile(fs FS, path string) (string, error) {
	return "", nil
}
