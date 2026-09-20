package main

import (
	"bytes"
	"io"
	"os"
	"path"
)

type FS interface {
	Exists(p string) bool
	Create(p string) (io.WriteCloser, error)
	Open(p string) (io.ReadCloser, error)
}

type memBlobFile struct {
	buf  *bytes.Buffer
	fs   *MockFS
	path string
}

func (f *memBlobFile) Write(p []byte) (int, error) {
	return f.buf.Write(p)
}

func (f *memBlobFile) Close() error {
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
	return &memBlobFile{buf: &bytes.Buffer{}, fs: m, path: p}, nil
}

func (m *MockFS) Open(p string) (io.ReadCloser, error) {
	data, ok := m.Files[p]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(data)), nil
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

func (s *Store) SaveBlob(hash string, r io.Reader) error {
	if s.HasBlob(hash) {
		return nil
	}
	w, err := s.fs.Create(s.BlobPath(hash))
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	return err
}

func (s *Store) OpenBlob(hash string) (io.ReadCloser, error) {
	return s.fs.Open(s.BlobPath(hash))
}
