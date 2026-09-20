package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

type FS interface {
	Open(path string) (io.ReadCloser, error)
	WalkDir(root string, fn func(path string, isDir bool) error) error
}

type Manifest map[string]string

type MockFS struct {
	Files map[string][]byte
	Dirs  []string
}

func NewMockFS() *MockFS {
	return &MockFS{
		Files: make(map[string][]byte),
		Dirs:  make([]string, 0),
	}
}

func (m *MockFS) Open(p string) (io.ReadCloser, error) {
	data, ok := m.Files[p]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(strings.NewReader(string(data))), nil
}

func (m *MockFS) WalkDir(root string, fn func(p string, isDir bool) error) error {
	cleanRoot := strings.TrimRight(strings.ReplaceAll(root, "\\", "/"), "/")
	for _, d := range m.Dirs {
		norm := strings.ReplaceAll(d, "\\", "/")
		if norm == cleanRoot || strings.HasPrefix(norm, cleanRoot+"/") {
			if err := fn(norm, true); err != nil {
				return err
			}
		}
	}
	for f := range m.Files {
		norm := strings.ReplaceAll(f, "\\", "/")
		if strings.HasPrefix(norm, cleanRoot+"/") || norm == cleanRoot {
			if err := fn(norm, false); err != nil {
				return err
			}
		}
	}
	return nil
}

// HashFile streams path from fs through SHA-256 and returns the hex digest.
func HashFile(fs FS, path string) (string, error) {
	f, err := fs.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// BuildManifest walks rootDir, hashes each file, and returns the manifest.
func BuildManifest(fs FS, rootDir string) (Manifest, error) {
	manifest := make(Manifest)
	cleanRoot := strings.TrimRight(strings.ReplaceAll(rootDir, "\\", "/"), "/")

	err := fs.WalkDir(rootDir, func(p string, isDir bool) error {
		if isDir {
			return nil
		}
		normPath := strings.ReplaceAll(p, "\\", "/")
		rel := strings.TrimPrefix(normPath, cleanRoot)
		if !strings.HasPrefix(rel, "/") {
			rel = "/" + rel
		}

		h, err := HashFile(fs, p)
		if err != nil {
			return err
		}
		manifest[rel] = h
		return nil
	})

	if err != nil {
		return nil, err
	}
	return manifest, nil
}
