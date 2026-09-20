package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

var ErrCorruptPayload = errors.New("corrupt payload")

type StagedFS interface {
	WriteFile(path string, data []byte) error
	Rename(oldPath, newPath string) error
	Remove(path string) error
	Exists(path string) bool
	ReadFile(path string) ([]byte, error)
}

type MockStagedFS struct {
	files map[string][]byte
}

func NewMockStagedFS() *MockStagedFS {
	return &MockStagedFS{files: make(map[string][]byte)}
}

func (m *MockStagedFS) WriteFile(p string, data []byte) error {
	cp := make([]byte, len(data))
	copy(cp, data)
	m.files[p] = cp
	return nil
}

func (m *MockStagedFS) Rename(oldPath, newPath string) error {
	data, ok := m.files[oldPath]
	if !ok {
		return os.ErrNotExist
	}
	delete(m.files, oldPath)
	m.files[newPath] = data
	return nil
}

func (m *MockStagedFS) Remove(p string) error {
	delete(m.files, p)
	return nil
}

func (m *MockStagedFS) Exists(p string) bool {
	_, ok := m.files[p]
	return ok
}

func (m *MockStagedFS) ReadFile(p string) ([]byte, error) {
	data, ok := m.files[p]
	if !ok {
		return nil, os.ErrNotExist
	}
	return data, nil
}

// SafeSaveBlob writes to staging, verifies digest, and atomically renames to blobs dir.
func SafeSaveBlob(fs StagedFS, claimedHash string, r io.Reader) error {
	return nil
}
