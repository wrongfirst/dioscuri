package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
)

type FS interface {
	Exists(path string) bool
	Open(path string) (io.ReadCloser, error)
}

type BlobOpener interface {
	OpenBlob(hash string) (io.ReadCloser, error)
}

type ActivePointer struct {
	DeployID string `json:"deploy_id"`
}

type MockFS struct {
	Files map[string][]byte
}

func NewMockFS() *MockFS {
	return &MockFS{Files: make(map[string][]byte)}
}

func (m *MockFS) Exists(p string) bool {
	_, ok := m.Files[p]
	return ok
}

func (m *MockFS) Open(p string) (io.ReadCloser, error) {
	data, ok := m.Files[p]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

type MockStore struct {
	Blobs map[string][]byte
}

func NewMockStore() *MockStore {
	return &MockStore{Blobs: make(map[string][]byte)}
}

func (m *MockStore) OpenBlob(hash string) (io.ReadCloser, error) {
	data, ok := m.Blobs[hash]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func detectContentType(filePath string) string {
	switch strings.ToLower(path.Ext(filePath)) {
	case ".html", ".htm":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}

// NewWebServer returns an http.Handler that routes production and preview requests via manifests.
func NewWebServer(fs FS, store BlobOpener) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var deployID string
		var filePath string

		cleanPath := path.Clean(r.URL.Path)

		if strings.HasPrefix(cleanPath, "/preview/") {
			rest := strings.TrimPrefix(cleanPath, "/preview/")
			parts := strings.SplitN(rest, "/", 2)
			deployID = parts[0]
			if len(parts) > 1 && parts[1] != "" {
				filePath = "/" + parts[1]
			} else {
				filePath = "/"
			}
		} else {
			prodReader, err := fs.Open("/production.json")
			if err != nil {
				http.NotFound(w, r)
				return
			}
			var active ActivePointer
			err = json.NewDecoder(prodReader).Decode(&active)
			prodReader.Close()
			if err != nil || active.DeployID == "" {
				http.NotFound(w, r)
				return
			}
			deployID = active.DeployID
			filePath = cleanPath
		}

		if filePath == "/" || strings.HasSuffix(filePath, "/") {
			filePath = path.Join(filePath, "index.html")
		}
		if !strings.HasPrefix(filePath, "/") {
			filePath = "/" + filePath
		}

		manifestReader, err := fs.Open("/manifests/" + deployID + ".json")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		var manifest map[string]string
		err = json.NewDecoder(manifestReader).Decode(&manifest)
		manifestReader.Close()
		if err != nil {
			http.NotFound(w, r)
			return
		}

		hash, ok := manifest[filePath]
		if !ok {
			http.NotFound(w, r)
			return
		}

		blobReader, err := store.OpenBlob(hash)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer blobReader.Close()

		w.Header().Set("Content-Type", detectContentType(filePath))
		w.WriteHeader(http.StatusOK)
		io.Copy(w, blobReader)
	})
}
