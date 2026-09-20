package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

type BlobStore interface {
	SaveBlob(hash string, r io.Reader) error
	RemoveBlob(hash string) error
}

type MockBlobStore struct {
	Blobs map[string][]byte
}

func NewMockBlobStore() *MockBlobStore {
	return &MockBlobStore{Blobs: make(map[string][]byte)}
}

func (m *MockBlobStore) SaveBlob(hash string, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	m.Blobs[hash] = data
	return nil
}

func (m *MockBlobStore) RemoveBlob(hash string) error {
	delete(m.Blobs, hash)
	return nil
}

// UploadBlobHandler streams r.Body directly to store while verifying content hash on the fly.
func UploadBlobHandler(store BlobStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
