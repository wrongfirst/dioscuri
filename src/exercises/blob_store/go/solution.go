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

func (s *Store) blobPath(hash string) string {
	return path.Join(s.blobDir, hash)
}

func (s *Store) HasBlob(hash string) bool {
	return s.fs.Exists(s.blobPath(hash))
}

func (s *Store) SaveBlob(hash string, r io.Reader) error {
	if s.HasBlob(hash) {
		return nil
	}
	w, err := s.fs.Create(s.blobPath(hash))
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	return err
}

func (s *Store) OpenBlob(hash string) (io.ReadCloser, error) {
	return s.fs.Open(s.blobPath(hash))
}
