package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

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

// DeltaCheckHandler returns an http.HandlerFunc that responds with missing hashes.
func DeltaCheckHandler(store BlobChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CheckRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		missing := make([]string, 0)
		for _, h := range req.Hashes {
			if !store.HasBlob(h) {
				missing = append(missing, h)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CheckResponse{Missing: missing})
	}
}
