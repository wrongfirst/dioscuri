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
	live := make(map[string]bool)
	for _, manifest := range activeManifests {
		for _, hash := range manifest {
			live[hash] = true
		}
	}

	existing, err := fs.ListBlobs()
	if err != nil {
		return 0, err
	}

	deleted := 0
	for _, hash := range existing {
		if !live[hash] {
			if err := fs.DeleteBlob(hash); err != nil {
				return deleted, err
			}
			deleted++
		}
	}

	return deleted, nil
}
