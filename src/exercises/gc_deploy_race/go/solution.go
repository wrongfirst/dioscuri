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
	live := make(map[string]bool)
	for _, m := range manifests {
		for _, hash := range m {
			live[hash] = true
		}
	}

	blobs, err := store.ListBlobs()
	if err != nil {
		return 0, err
	}

	deleted := 0
	for _, b := range blobs {
		if live[b.Hash] {
			continue
		}
		if now.Sub(b.CreatedAt) > quarantine {
			if err := store.DeleteBlob(b.Hash); err != nil {
				return deleted, err
			}
			deleted++
		}
	}

	return deleted, nil
}
