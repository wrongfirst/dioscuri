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
		if r.Method != http.MethodPut && r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		claimedHash := strings.TrimPrefix(r.URL.Path, "/api/blob/")
		if claimedHash == "" {
			claimedHash = r.URL.Query().Get("hash")
		}
		if claimedHash == "" || strings.Contains(claimedHash, "/") {
			http.Error(w, "invalid or missing hash", http.StatusBadRequest)
			return
		}

		hasher := sha256.New()
		tee := io.TeeReader(r.Body, hasher)

		if err := store.SaveBlob(claimedHash, tee); err != nil {
			http.Error(w, "storage error", http.StatusInternalServerError)
			return
		}

		actualHash := hex.EncodeToString(hasher.Sum(nil))
		if actualHash != claimedHash {
			store.RemoveBlob(claimedHash)
			http.Error(w, "hash mismatch", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
