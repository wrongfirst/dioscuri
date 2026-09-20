package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

var (
	ErrPowerCut = errors.New("power cut")
	ErrDiskFull = errors.New("disk full")
)

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
	if _, ok := m.Files[p]; !ok {
		return os.ErrNotExist
	}
	delete(m.Files, p)
	return nil
}

func (m *MockFS) MkdirAll(p string) error { return nil }

func (m *MockFS) ReadDir(dir string) ([]string, error) {
	return nil, nil
}

func (m *MockFS) WalkDir(root string, fn func(p string, isDir bool) error) error {
	for f := range m.Files {
		if strings.HasPrefix(f, root) {
			if err := fn(f, false); err != nil {
				return err
			}
		}
	}
	return nil
}

type FaultyFS struct {
	base       FS
	diskFull   bool
	byteBudget int // -1 if disabled
}

func NewFaultyFS(base FS) *FaultyFS {
	return &FaultyFS{
		base:       base,
		diskFull:   false,
		byteBudget: -1,
	}
}

func (f *FaultyFS) SetDiskFull(full bool) {
	f.diskFull = full
}

func (f *FaultyFS) CrashAfter(bytes int) {
	f.byteBudget = bytes
}

type faultyWriter struct {
	w  io.WriteCloser
	fs *FaultyFS
}

func (fw *faultyWriter) Write(p []byte) (int, error) {
	return 0, nil
}

func (fw *faultyWriter) Close() error {
	return fw.w.Close()
}

func (f *FaultyFS) Create(p string) (io.WriteCloser, error) {
	return nil, nil
}

func (f *FaultyFS) CorruptByte(p string, offset int) error {
	return nil
}

func (f *FaultyFS) Open(p string) (io.ReadCloser, error)  { return f.base.Open(p) }
func (f *FaultyFS) Rename(old, new string) error          { return f.base.Rename(old, new) }
func (f *FaultyFS) Remove(p string) error                 { return f.base.Remove(p) }
func (f *FaultyFS) Exists(p string) bool                  { return f.base.Exists(p) }
func (f *FaultyFS) MkdirAll(p string) error               { return f.base.MkdirAll(p) }
func (f *FaultyFS) ReadDir(dir string) ([]string, error)  { return f.base.ReadDir(dir) }
func (f *FaultyFS) WalkDir(root string, fn func(string, bool) error) error {
	return f.base.WalkDir(root, fn)
}
