package main

import (
	"bytes"
	"io"
	"os"
	"path"
)

type BlobFS interface {
	Exists(p string) bool
	Create(p string) (io.WriteCloser, error)
	Open(p string) (io.ReadCloser, error)
}

type memBlobFile struct {
	buf  *bytes.Buffer
	fs   *MemBlobFS
	path string
}

func (f *memBlobFile) Write(p []byte) (int, error) {
	return f.buf.Write(p)
}

func (f *memBlobFile) Close() error {
	f.fs.Files[f.path] = f.buf.Bytes()
	return nil
}

type MemBlobFS struct {
	Files map[string][]byte
}

func NewMemBlobFS() *MemBlobFS {
	return &MemBlobFS{Files: make(map[string][]byte)}
}

func (m *MemBlobFS) Exists(p string) bool {
	_, ok := m.Files[p]
	return ok
}

func (m *MemBlobFS) Create(p string) (io.WriteCloser, error) {
	return &memBlobFile{buf: &bytes.Buffer{}, fs: m, path: p}, nil
}

func (m *MemBlobFS) Open(p string) (io.ReadCloser, error) {
	data, ok := m.Files[p]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

type Store struct {
	fs      BlobFS
	blobDir string
}

func NewStore(fs BlobFS, blobDir string) *Store {
	return &Store{fs: fs, blobDir: blobDir}
}

func (s *Store) HasBlob(hash string) bool {
	return false
}

func (s *Store) SaveBlob(hash string, r io.Reader) error {
	return nil
}

func (s *Store) OpenBlob(hash string) (io.ReadCloser, error) {
	return nil, nil
}
