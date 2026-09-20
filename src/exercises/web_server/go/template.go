package main

import (
	"bytes"
	"io"
	"os"
)

type BlobOpener interface {
	OpenBlob(hash string) (io.ReadCloser, error)
}

type RouterState struct {
	ActiveDeployID string
	Manifests      map[string]map[string]string
}

type Response struct {
	StatusCode  int
	ContentType string
	Body        []byte
}

type MockBlobOpener struct {
	Blobs map[string][]byte
}

func NewMockBlobOpener() *MockBlobOpener {
	return &MockBlobOpener{Blobs: make(map[string][]byte)}
}

func (m *MockBlobOpener) OpenBlob(hash string) (io.ReadCloser, error) {
	data, ok := m.Blobs[hash]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

// RouteRequest routes production and preview requests to the correct blob.
func RouteRequest(state *RouterState, store BlobOpener, urlPath string) Response {
	return Response{StatusCode: 404}
}
