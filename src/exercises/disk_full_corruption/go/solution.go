package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrDiskFull = errors.New("disk full")

type AtomicFS interface {
	WriteFile(path string, data []byte) error
	Rename(oldPath, newPath string) error
	Remove(path string) error
	ReadFile(path string) ([]byte, error)
}

type MockAtomicFS struct {
	Files    map[string][]byte
	DiskFull bool
}

func NewMockAtomicFS() *MockAtomicFS {
	return &MockAtomicFS{Files: make(map[string][]byte)}
}

func (m *MockAtomicFS) WriteFile(p string, data []byte) error {
	if m.DiskFull {
		return ErrDiskFull
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	m.Files[p] = cp
	return nil
}

func (m *MockAtomicFS) Rename(oldPath, newPath string) error {
	data, ok := m.Files[oldPath]
	if !ok {
		return os.ErrNotExist
	}
	delete(m.Files, oldPath)
	m.Files[newPath] = data
	return nil
}

func (m *MockAtomicFS) Remove(p string) error {
	delete(m.Files, p)
	return nil
}

func (m *MockAtomicFS) ReadFile(p string) ([]byte, error) {
	data, ok := m.Files[p]
	if !ok {
		return nil, os.ErrNotExist
	}
	return data, nil
}

// UpdateProductionPointer writes to temp file and renames to production.json.
func UpdateProductionPointer(fs AtomicFS, deployID string) error {
	payload := []byte(fmt.Sprintf("{\"active\":\"%s\"}", deployID))
	tmpPath := "production.json.tmp"
	targetPath := "production.json"

	if err := fs.WriteFile(tmpPath, payload); err != nil {
		return err
	}

	if err := fs.Rename(tmpPath, targetPath); err != nil {
		fs.Remove(tmpPath)
		return err
	}

	return nil
}
