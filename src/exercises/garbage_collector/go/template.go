package main

type GCFS interface {
	ListBlobs() ([]string, error)
	DeleteBlob(hash string) error
}

type MockGCFS struct {
	Blobs map[string]bool
}

func NewMockGCFS(hashes ...string) *MockGCFS {
	m := &MockGCFS{Blobs: make(map[string]bool)}
	for _, h := range hashes {
		m.Blobs[h] = true
	}
	return m
}

func (m *MockGCFS) ListBlobs() ([]string, error) {
	list := make([]string, 0, len(m.Blobs))
	for h := range m.Blobs {
		list = append(list, h)
	}
	return list, nil
}

func (m *MockGCFS) DeleteBlob(hash string) error {
	delete(m.Blobs, hash)
	return nil
}

// RunGC performs mark-and-sweep to delete unreferenced blobs.
func RunGC(fs GCFS, activeManifests []map[string]string) (int, error) {
	return 0, nil
}
