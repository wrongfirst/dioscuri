package main

import (
	"bytes"
	"io"
	"os"
	"path"
	"strings"
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

func getContentType(p string) string {
	ext := path.Ext(p)
	switch ext {
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	default:
		return "application/octet-stream"
	}
}

// RouteRequest routes production and preview requests to the correct blob.
func RouteRequest(state *RouterState, store BlobOpener, urlPath string) Response {
	targetDeployID := state.ActiveDeployID
	targetPath := urlPath

	if strings.HasPrefix(urlPath, "/preview/") {
		trimmed := strings.TrimPrefix(urlPath, "/preview/")
		parts := strings.SplitN(trimmed, "/", 2)
		if len(parts) < 2 {
			return Response{StatusCode: 404}
		}
		targetDeployID = parts[0]
		targetPath = "/" + parts[1]
	}

	if targetPath == "/" {
		targetPath = "/index.html"
	}

	manifest, ok := state.Manifests[targetDeployID]
	if !ok {
		return Response{StatusCode: 404}
	}

	hash, ok := manifest[targetPath]
	if !ok {
		return Response{StatusCode: 404}
	}

	r, err := store.OpenBlob(hash)
	if err != nil {
		return Response{StatusCode: 404}
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return Response{StatusCode: 500}
	}

	return Response{
		StatusCode:  200,
		ContentType: getContentType(targetPath),
		Body:        data,
	}
}
