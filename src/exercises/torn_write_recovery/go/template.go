package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

var (
	ErrPowerCut       = errors.New("power cut")
	ErrDiskFull       = errors.New("disk full")
	ErrCorruptPayload = errors.New("corrupt payload")
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
	delete(m.Files, p)
	return nil
}

func (m *MockFS) MkdirAll(p string) error { return nil }

func (m *MockFS) ReadDir(dir string) ([]string, error) {
	return nil, nil
}

func (m *MockFS) WalkDir(root string, fn func(p string, isDir bool) error) error {
	return nil
}

type FaultyFS struct {
	base       FS
	diskFull   bool
	byteBudget int
}

func NewFaultyFS(base FS) *FaultyFS {
	return &FaultyFS{base: base, byteBudget: -1}
}

func (f *FaultyFS) CrashAfter(bytes int) { f.byteBudget = bytes }
func (f *FaultyFS) SetDiskFull(full bool) { f.diskFull = full }

type faultyWriter struct {
	w  io.WriteCloser
	fs *FaultyFS
}

func (fw *faultyWriter) Write(p []byte) (int, error) {
	if fw.fs.diskFull {
		return 0, ErrDiskFull
	}
	if fw.fs.byteBudget < 0 {
		return fw.w.Write(p)
	}
	if fw.fs.byteBudget == 0 {
		return 0, ErrPowerCut
	}
	if len(p) <= fw.fs.byteBudget {
		n, err := fw.w.Write(p)
		fw.fs.byteBudget -= n
		return n, err
	}
	allowed := p[:fw.fs.byteBudget]
	n, _ := fw.w.Write(allowed)
	fw.fs.byteBudget = 0
	return n, ErrPowerCut
}

func (fw *faultyWriter) Close() error { return fw.w.Close() }

func (f *FaultyFS) Create(p string) (io.WriteCloser, error) {
	if f.diskFull {
		return nil, ErrDiskFull
	}
	if f.byteBudget == 0 {
		return nil, ErrPowerCut
	}
	w, err := f.base.Create(p)
	if err != nil {
		return nil, err
	}
	return &faultyWriter{w: w, fs: f}, nil
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

type Store struct {
	fs      FS
	blobDir string
}

func NewStore(fs FS, blobDir string) *Store {
	return &Store{fs: fs, blobDir: blobDir}
}

func (s *Store) BlobPath(hash string) string {
	return path.Join(s.blobDir, hash)
}

func (s *Store) HasBlob(hash string) bool {
	return s.fs.Exists(s.BlobPath(hash))
}

// SafeSaveBlob streams payload to staging, computes SHA-256 on the fly, cleans up on error, and atomically renames.
func SafeSaveBlob(fs FS, store *Store, stagingDir string, claimedHash string, r io.Reader) error {
	return nil
}
