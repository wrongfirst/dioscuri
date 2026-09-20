package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
)

type FS interface {
	Exists(path string) bool
	Open(path string) (io.ReadCloser, error)
}

type BlobOpener interface {
	OpenBlob(hash string) (io.ReadCloser, error)
}

type ActivePointer struct {
	DeployID string `json:"deploy_id"`
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

func (m *MockFS) Open(p string) (io.ReadCloser, error) {
	data, ok := m.Files[p]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

type MockStore struct {
	Blobs map[string][]byte
}

func NewMockStore() *MockStore {
	return &MockStore{Blobs: make(map[string][]byte)}
}

func (m *MockStore) OpenBlob(hash string) (io.ReadCloser, error) {
	data, ok := m.Blobs[hash]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

// NewWebServer returns an http.Handler that routes production and preview requests via manifests.
func NewWebServer(fs FS, store BlobOpener) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
