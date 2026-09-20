package main

import (
	"os"
	"sort"
	"strings"
)

type MemoryFS struct {
	files map[string][]byte
}

func NewMemoryFS() *MemoryFS {
	return &MemoryFS{
		files: make(map[string][]byte),
	}
}

func (m *MemoryFS) WriteFile(path string, data []byte) error {
	return nil
}

func (m *MemoryFS) ReadFile(path string) ([]byte, error) {
	return nil, nil
}

func (m *MemoryFS) Rename(oldPath, newPath string) error {
	return nil
}

func (m *MemoryFS) Remove(path string) error {
	return nil
}

func (m *MemoryFS) Exists(path string) bool {
	return false
}

func (m *MemoryFS) ListDir(prefix string) []string {
	return nil
}
