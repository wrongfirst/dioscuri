package main

import (
	"time"
)

type StoredBlob struct {
	Hash      string
	CreatedAt time.Time
}

type GCStore interface {
	ListBlobs() ([]StoredBlob, error)
	DeleteBlob(hash string) error
}

type MockGCStore struct {
	blobs map[string]time.Time
}

func NewMockGCStore() *MockGCStore {
	return &MockGCStore{blobs: make(map[string]time.Time)}
}

func (m *MockGCStore) AddBlob(hash string, t time.Time) {
	m.blobs[hash] = t
}

func (m *MockGCStore) ListBlobs() ([]StoredBlob, error) {
	list := make([]StoredBlob, 0, len(m.blobs))
	for h, t := range m.blobs {
		list = append(list, StoredBlob{Hash: h, CreatedAt: t})
	}
	return list, nil
}

func (m *MockGCStore) DeleteBlob(hash string) error {
	delete(m.blobs, hash)
	return nil
}

// RunQuarantinedGC cleans up unreferenced blobs older than quarantine duration.
func RunQuarantinedGC(store GCStore, manifests []map[string]string, now time.Time, quarantine time.Duration) (int, error) {
	return 0, nil
}
