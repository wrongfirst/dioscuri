package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strings"
)

type WalkableFS interface {
	Walk(root string, fn func(path string, isDir bool) error) error
	Open(path string) (io.ReadCloser, error)
}

type Manifest map[string]string

type MockWalkableFS struct {
	Files map[string][]byte
	Dirs  []string
}

func NewMockWalkableFS() *MockWalkableFS {
	return &MockWalkableFS{
		Files: make(map[string][]byte),
		Dirs:  make([]string, 0),
	}
}

func (m *MockWalkableFS) Open(p string) (io.ReadCloser, error) {
	data, ok := m.Files[p]
	if !ok {
		return nil, io.EOF
	}
	return io.NopCloser(strings.NewReader(string(data))), nil
}

func (m *MockWalkableFS) Walk(root string, fn func(p string, isDir bool) error) error {
	for _, d := range m.Dirs {
		if strings.HasPrefix(d, root) {
			if err := fn(d, true); err != nil {
				return err
			}
		}
	}
	for f := range m.Files {
		if strings.HasPrefix(f, root) {
			if err := fn(f, false); err != nil {
				return err
			}
		}
	}
	return nil
}

// BuildManifest walks rootDir, hashes each file, and returns the manifest.
func BuildManifest(fs WalkableFS, rootDir string) (Manifest, error) {
	return nil, nil
}
