package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path"
	"strings"
)

type FS interface {
	ReadDir(dir string) ([]string, error)
	Open(path string) (io.ReadCloser, error)
	Remove(path string) error
}

type MockFS struct {
	Files map[string][]byte
}

func NewMockFS() *MockFS {
	return &MockFS{Files: make(map[string][]byte)}
}

func (m *MockFS) ReadDir(dir string) ([]string, error) {
	cleanDir := strings.TrimRight(dir, "/")
	prefix := cleanDir + "/"
	seen := make(map[string]bool)
	var list []string
	for p := range m.Files {
		if strings.HasPrefix(p, prefix) {
			rel := strings.TrimPrefix(p, prefix)
			parts := strings.Split(rel, "/")
			if len(parts) > 0 && !seen[parts[0]] {
				seen[parts[0]] = true
				list = append(list, parts[0])
			}
		}
	}
	return list, nil
}

func (m *MockFS) Open(p string) (io.ReadCloser, error) {
	data, ok := m.Files[p]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (m *MockFS) Remove(p string) error {
	if _, ok := m.Files[p]; !ok {
		return os.ErrNotExist
	}
	delete(m.Files, p)
	return nil
}

// RunGC performs mark-and-sweep across manifestDir and blobDir, deleting unreferenced blobs.
func RunGC(fs FS, manifestDir string, blobDir string) (int, error) {
	return 0, nil
}
