package main

type BlobChecker interface {
	HasBlob(hash string) bool
}

type CheckRequest struct {
	Hashes []string `json:"hashes"`
}

type CheckResponse struct {
	Missing []string `json:"missing"`
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

// DeltaCheck returns the subset of hashes missing from store.
func DeltaCheck(store BlobChecker, req CheckRequest) CheckResponse {
	missing := make([]string, 0)
	for _, h := range req.Hashes {
		if !store.HasBlob(h) {
			missing = append(missing, h)
		}
	}
	return CheckResponse{Missing: missing}
}
