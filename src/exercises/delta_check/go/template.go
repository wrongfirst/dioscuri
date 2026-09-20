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
	}
}
