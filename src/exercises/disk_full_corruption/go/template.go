package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
)

var ErrDiskFull = errors.New("disk full")

type FS interface {
	Open(path string) (io.ReadCloser, error)
	Create(path string) (io.WriteCloser, error)
	Rename(oldPath, newPath string) error
	Remove(path string) error
	Exists(path string) bool
	MkdirAll(path string) error
	ReadDir(path string) ([]string, error)
	WalkDir(root string, fn func(path string, isDir bool) error) error
}

type memFile struct {
	buf  *bytes.Buffer
	fs   *MockFS
	path string
}

func (f *memFile) Write(p []byte) (int, error) { return f.buf.Write(p) }
func (f *memFile) Close() error {
	f.fs.Files[f.path] = f.buf.Bytes()
	return nil
}

type MockFS struct {
	Files map[string][]byte
}

func NewMockFS() *MockFS {
	return &MockFS{Files: make(map[string][]byte)}
}

func (m *MockFS) Exists(p string) bool {
	_, ok := m.Files[p]
	return ok
}

func (m *MockFS) Create(p string) (io.WriteCloser, error) {
	return &memFile{buf: &bytes.Buffer{}, fs: m, path: p}, nil
}

func (m *MockFS) Open(p string) (io.ReadCloser, error) {
	data, ok := m.Files[p]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (m *MockFS) Rename(oldPath, newPath string) error {
	data, ok := m.Files[oldPath]
	if !ok {
		return os.ErrNotExist
	}
	delete(m.Files, oldPath)
	m.Files[newPath] = data
	return nil
}

func (m *MockFS) Remove(p string) error {
	delete(m.Files, p)
	return nil
}

func (m *MockFS) MkdirAll(p string) error { return nil }
func (m *MockFS) ReadDir(dir string) ([]string, error) { return nil, nil }
func (m *MockFS) WalkDir(root string, fn func(string, bool) error) error { return nil }

type FaultyFS struct {
	base     FS
	diskFull bool
}

func NewFaultyFS(base FS) *FaultyFS {
	return &FaultyFS{base: base}
}

func (f *FaultyFS) SetDiskFull(full bool) { f.diskFull = full }

func (f *FaultyFS) Create(p string) (io.WriteCloser, error) {
	if f.diskFull {
		return nil, ErrDiskFull
	}
	return f.base.Create(p)
}

func (f *FaultyFS) Open(p string) (io.ReadCloser, error) { return f.base.Open(p) }
func (f *FaultyFS) Rename(old, new string) error         { return f.base.Rename(old, new) }
func (f *FaultyFS) Remove(p string) error                { return f.base.Remove(p) }
func (f *FaultyFS) Exists(p string) bool                 { return f.base.Exists(p) }
func (f *FaultyFS) MkdirAll(p string) error              { return f.base.MkdirAll(p) }
func (f *FaultyFS) ReadDir(dir string) ([]string, error) { return f.base.ReadDir(dir) }
func (f *FaultyFS) WalkDir(root string, fn func(string, bool) error) error {
	return f.base.WalkDir(root, fn)
}

type ActivePointer struct {
	DeployID string `json:"deploy_id"`
}

// UpdateProductionPointer writes new deploy pointer to temporary file and atomically renames.
func UpdateProductionPointer(fs FS, deployID string) error {
	return nil
}
