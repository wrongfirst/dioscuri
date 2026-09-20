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
	cp := make([]byte, len(data))
	copy(cp, data)
	m.files[path] = cp
	return nil
}

func (m *MemoryFS) ReadFile(path string) ([]byte, error) {
	data, ok := m.files[path]
	if !ok {
		return nil, os.ErrNotExist
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	return cp, nil
}

func (m *MemoryFS) Rename(oldPath, newPath string) error {
	data, ok := m.files[oldPath]
	if !ok {
		return os.ErrNotExist
	}
	delete(m.files, oldPath)
	m.files[newPath] = data
	return nil
}

func (m *MemoryFS) Remove(path string) error {
	if _, ok := m.files[path]; !ok {
		return os.ErrNotExist
	}
	delete(m.files, path)
	return nil
}

func (m *MemoryFS) Exists(path string) bool {
	_, ok := m.files[path]
	return ok
}

func (m *MemoryFS) ListDir(prefix string) []string {
	res := make([]string, 0)
	for p := range m.files {
		if strings.HasPrefix(p, prefix) {
			res = append(res, p)
		}
	}
	sort.Strings(res)
	return res
}
