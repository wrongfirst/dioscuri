package main

import (
	"bytes"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
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

type MemoryFS struct {
	mu    sync.RWMutex
	files map[string][]byte
	dirs  map[string]bool
}

func NewMemoryFS() *MemoryFS {
	return &MemoryFS{
		files: make(map[string][]byte),
		dirs:  make(map[string]bool),
	}
}

type memWriter struct {
	buf  bytes.Buffer
	fs   *MemoryFS
	path string
}

func (w *memWriter) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func (w *memWriter) Close() error {
	w.fs.mu.Lock()
	defer w.fs.mu.Unlock()

	data := make([]byte, w.buf.Len())
	copy(data, w.buf.Bytes())
	w.fs.files[w.path] = data

	dir := path.Dir(w.path)
	for dir != "/" && dir != "." && dir != "" {
		w.fs.dirs[dir] = true
		dir = path.Dir(dir)
	}
	return nil
}

func (m *MemoryFS) Open(p string) (io.ReadCloser, error) {
	return nil, nil
}

func (m *MemoryFS) Create(p string) (io.WriteCloser, error) {
	return nil, nil
}

func (m *MemoryFS) Rename(oldPath, newPath string) error {
	return nil
}

func (m *MemoryFS) Remove(p string) error {
	return nil
}

func (m *MemoryFS) Exists(p string) bool {
	return false
}

func (m *MemoryFS) MkdirAll(p string) error {
	return nil
}

func (m *MemoryFS) ReadDir(dir string) ([]string, error) {
	return nil, nil
}

func (m *MemoryFS) WalkDir(root string, fn func(path string, isDir bool) error) error {
	return nil
}
