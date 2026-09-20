package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
)

var (
	ErrMissingBlobs   = errors.New("missing blobs in store")
	ErrDeployNotFound = errors.New("deployment not found")
)

type FS interface {
	Exists(path string) bool
	Create(path string) (io.WriteCloser, error)
	Open(path string) (io.ReadCloser, error)
	Rename(oldPath, newPath string) error
}

type BlobChecker interface {
	HasBlob(hash string) bool
}

type Manifest map[string]string

type ActivePointer struct {
	DeployID string `json:"deploy_id"`
}

type memFile struct {
	buf  *bytes.Buffer
	fs   *MockFS
	path string
}

func (f *memFile) Write(p []byte) (int, error) {
	return f.buf.Write(p)
}

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

type MockBlobChecker struct {
	Existing map[string]bool
}

func NewMockBlobChecker(hashes ...string) *MockBlobChecker {
	m := &MockBlobChecker{Existing: make(map[string]bool)}
	for _, h := range hashes {
		m.Existing[h] = true
	}
	return m
}

func (m *MockBlobChecker) HasBlob(hash string) bool {
	return m.Existing[hash]
}

// CommitDeploy validates blobs, persists the manifest, and updates production.json atomically.
func CommitDeploy(fs FS, checker BlobChecker, deployID string, manifest Manifest) error {
	return nil
}

// RollbackDeploy updates production.json atomically to targetDeployID.
func RollbackDeploy(fs FS, targetDeployID string) error {
	return nil
}
