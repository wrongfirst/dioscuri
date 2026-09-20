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
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, ok := m.files[p]
	if !ok {
		return nil, os.ErrNotExist
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	return io.NopCloser(bytes.NewReader(cp)), nil
}

func (m *MemoryFS) Create(p string) (io.WriteCloser, error) {
	return &memWriter{
		fs:   m,
		path: p,
	}, nil
}

func (m *MemoryFS) Rename(oldPath, newPath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, ok := m.files[oldPath]
	if !ok {
		return os.ErrNotExist
	}
	delete(m.files, oldPath)
	m.files[newPath] = data

	dir := path.Dir(newPath)
	for dir != "/" && dir != "." && dir != "" {
		m.dirs[dir] = true
		dir = path.Dir(dir)
	}
	return nil
}

func (m *MemoryFS) Remove(p string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.files[p]; !ok {
		return os.ErrNotExist
	}
	delete(m.files, p)
	return nil
}

func (m *MemoryFS) Exists(p string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.files[p]; ok {
		return true
	}
	return m.dirs[p]
}

func (m *MemoryFS) MkdirAll(p string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	clean := path.Clean(p)
	for clean != "/" && clean != "." && clean != "" {
		m.dirs[clean] = true
		clean = path.Dir(clean)
	}
	return nil
}

func (m *MemoryFS) ReadDir(dir string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cleanDir := strings.TrimRight(dir, "/")
	prefix := cleanDir + "/"
	seen := make(map[string]bool)
	var list []string

	for p := range m.files {
		if strings.HasPrefix(p, prefix) {
			rel := strings.TrimPrefix(p, prefix)
			parts := strings.Split(rel, "/")
			if len(parts) > 0 && !seen[parts[0]] {
				seen[parts[0]] = true
				list = append(list, parts[0])
			}
		}
	}
	for d := range m.dirs {
		if strings.HasPrefix(d, prefix) {
			rel := strings.TrimPrefix(d, prefix)
			parts := strings.Split(rel, "/")
			if len(parts) > 0 && parts[0] != "" && !seen[parts[0]] {
				seen[parts[0]] = true
				list = append(list, parts[0])
			}
		}
	}

	sort.Strings(list)
	return list, nil
}

func (m *MemoryFS) WalkDir(root string, fn func(path string, isDir bool) error) error {
	m.mu.RLock()
	cleanRoot := strings.TrimRight(strings.ReplaceAll(root, "\\", "/"), "/")
	var dirs []string
	var files []string

	for d := range m.dirs {
		norm := strings.ReplaceAll(d, "\\", "/")
		if norm == cleanRoot || strings.HasPrefix(norm, cleanRoot+"/") {
			dirs = append(dirs, norm)
		}
	}
	for f := range m.files {
		norm := strings.ReplaceAll(f, "\\", "/")
		if norm == cleanRoot || strings.HasPrefix(norm, cleanRoot+"/") {
			files = append(files, norm)
		}
	}
	m.mu.RUnlock()

	sort.Strings(dirs)
	sort.Strings(files)

	for _, d := range dirs {
		if err := fn(d, true); err != nil {
			return err
		}
	}
	for _, f := range files {
		if err := fn(f, false); err != nil {
			return err
		}
	}
	return nil
}
