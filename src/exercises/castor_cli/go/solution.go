package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
)

var (
	ErrMissingBlobs   = errors.New("missing blobs in store")
	ErrDeployNotFound = errors.New("deployment not found")
)

type FS interface {
	Exists(path string) bool
	Create(path string) (io.WriteCloser, error)
	Open(path string) (io.ReadCloser, error)
	Rename(oldPath, newPath string) error
	ReadDir(path string) ([]string, error)
	WalkDir(root string, fn func(path string, isDir bool) error) error
}

type Manifest map[string]string

type ActivePointer struct {
	DeployID string `json:"deploy_id"`
}

type memFile struct {
	buf  *bytes.Buffer
	fs   *MockFS
	path string
}

func (f *memFile) Write(p []byte) (int, error) { return f.buf.Write(p) }
func (f *memFile) Close() error {
	f.fs.Files[f.path] = f.buf.Bytes()
	return nil
}

type MockFS struct {
	Files map[string][]byte
	Dirs  []string
}

func NewMockFS() *MockFS {
	return &MockFS{Files: make(map[string][]byte), Dirs: make([]string, 0)}
}

func (m *MockFS) Exists(p string) bool {
	_, ok := m.Files[p]
	return ok
}

func (m *MockFS) Create(p string) (io.WriteCloser, error) {
	return &memFile{buf: &bytes.Buffer{}, fs: m, path: p}, nil
}

func (m *MockFS) Open(p string) (io.ReadCloser, error) {
	data, ok := m.Files[p]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (m *MockFS) Rename(oldPath, newPath string) error {
	data, ok := m.Files[oldPath]
	if !ok {
		return os.ErrNotExist
	}
	delete(m.Files, oldPath)
	m.Files[newPath] = data
	return nil
}

func (m *MockFS) ReadDir(dir string) ([]string, error) {
	cleanDir := strings.TrimRight(dir, "/")
	prefix := cleanDir + "/"
	seen := make(map[string]bool)
	var list []string
	for p := range m.Files {
		if strings.HasPrefix(p, prefix) {
			rel := strings.TrimPrefix(p, prefix)
			parts := strings.Split(rel, "/")
			if len(parts) > 0 && !seen[parts[0]] {
				seen[parts[0]] = true
				list = append(list, parts[0])
			}
		}
	}
	return list, nil
}

func (m *MockFS) WalkDir(root string, fn func(p string, isDir bool) error) error {
	cleanRoot := strings.TrimRight(strings.ReplaceAll(root, "\\", "/"), "/")
	for _, d := range m.Dirs {
		norm := strings.ReplaceAll(d, "\\", "/")
		if norm == cleanRoot || strings.HasPrefix(norm, cleanRoot+"/") {
			if err := fn(norm, true); err != nil {
				return err
			}
		}
	}
	for f := range m.Files {
		norm := strings.ReplaceAll(f, "\\", "/")
		if strings.HasPrefix(norm, cleanRoot+"/") || norm == cleanRoot {
			if err := fn(norm, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func HashFile(fs FS, path string) (string, error) {
	f, err := fs.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func BuildManifest(fs FS, rootDir string) (Manifest, error) {
	manifest := make(Manifest)
	cleanRoot := strings.TrimRight(strings.ReplaceAll(rootDir, "\\", "/"), "/")

	err := fs.WalkDir(rootDir, func(p string, isDir bool) error {
		if isDir {
			return nil
		}
		normPath := strings.ReplaceAll(p, "\\", "/")
		rel := strings.TrimPrefix(normPath, cleanRoot)
		if !strings.HasPrefix(rel, "/") {
			rel = "/" + rel
		}

		h, err := HashFile(fs, p)
		if err != nil {
			return err
		}
		manifest[rel] = h
		return nil
	})
	if err != nil {
		return nil, err
	}
	return manifest, nil
}

type Store struct {
	fs      FS
	blobDir string
}

func NewStore(fs FS, blobDir string) *Store {
	return &Store{fs: fs, blobDir: blobDir}
}

func (s *Store) BlobPath(hash string) string {
	return path.Join(s.blobDir, hash)
}

func (s *Store) HasBlob(hash string) bool {
	return s.fs.Exists(s.BlobPath(hash))
}

func (s *Store) SaveBlob(hash string, r io.Reader) error {
	if s.HasBlob(hash) {
		return nil
	}
	w, err := s.fs.Create(s.BlobPath(hash))
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	return err
}

func (s *Store) RemoveBlob(hash string) error {
	return nil
}

func (s *Store) OpenBlob(hash string) (io.ReadCloser, error) {
	return s.fs.Open(s.BlobPath(hash))
}

func DeltaCheckHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Hashes []string `json:"hashes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		missing := make([]string, 0)
		for _, h := range req.Hashes {
			if !store.HasBlob(h) {
				missing = append(missing, h)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]string{"missing": missing})
	}
}

func UploadBlobHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut && r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		claimedHash := strings.TrimPrefix(r.URL.Path, "/api/blob/")
		if claimedHash == "" {
			claimedHash = r.URL.Query().Get("hash")
		}
		if claimedHash == "" || strings.Contains(claimedHash, "/") {
			http.Error(w, "invalid hash", http.StatusBadRequest)
			return
		}

		hasher := sha256.New()
		tee := io.TeeReader(r.Body, hasher)
		if err := store.SaveBlob(claimedHash, tee); err != nil {
			http.Error(w, "storage error", http.StatusInternalServerError)
			return
		}
		if hex.EncodeToString(hasher.Sum(nil)) != claimedHash {
			store.RemoveBlob(claimedHash)
			http.Error(w, "hash mismatch", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func CommitDeploy(fs FS, store *Store, deployID string, manifest Manifest) error {
	for _, hash := range manifest {
		if !store.HasBlob(hash) {
			return ErrMissingBlobs
		}
	}
	manifestPath := "/manifests/" + deployID + ".json"
	w, err := fs.Create(manifestPath)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(manifest); err != nil {
		w.Close()
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	tmpPath := "/production.json.tmp"
	pw, err := fs.Create(tmpPath)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(pw).Encode(ActivePointer{DeployID: deployID}); err != nil {
		pw.Close()
		return err
	}
	if err := pw.Close(); err != nil {
		return err
	}

	return fs.Rename(tmpPath, "/production.json")
}

func NewWebServer(fs FS, store *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		filePath := path.Clean(r.URL.Path)
		if filePath == "/" || strings.HasSuffix(filePath, "/") {
			filePath = path.Join(filePath, "index.html")
		}

		mReader, err := fs.Open("/manifests/" + active.DeployID + ".json")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		var manifest map[string]string
		err = json.NewDecoder(mReader).Decode(&manifest)
		mReader.Close()
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

		if strings.HasSuffix(filePath, ".html") {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		} else if strings.HasSuffix(filePath, ".css") {
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
		}
		w.WriteHeader(http.StatusOK)
		io.Copy(w, blobReader)
	})
}

// Client communicates with a remote Castor deployment server.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		HTTPClient: http.DefaultClient,
	}
}

func (c *Client) Check(hashes []string) ([]string, error) {
	payload, _ := json.Marshal(map[string][]string{"hashes": hashes})
	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/check", "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res struct {
		Missing []string `json:"missing"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Missing, nil
}

func (c *Client) UploadBlob(hash string, r io.Reader) error {
	req, err := http.NewRequest(http.MethodPut, c.BaseURL+"/api/blob/"+hash, r)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload status %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) Deploy(deployID string, manifest map[string]string) error {
	payload, _ := json.Marshal(map[string]any{
		"deploy_id": deployID,
		"manifest":  manifest,
	})
	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/deploy", "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("deploy status %d", resp.StatusCode)
	}
	return nil
}

type DeployStats struct {
	DeployID      string
	TotalFiles    int
	UploadedFiles int
}

// ====================================================================
// ASSEMBLED SERVER MUX AND CLIENT PIPELINE
// ====================================================================

// NewServer mounts delta checking, blob uploading, deploy committing, and web serving onto a ServeMux.
func NewServer(fs FS, store *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/check", DeltaCheckHandler(store))
	mux.HandleFunc("/api/blob/", UploadBlobHandler(store))
	mux.HandleFunc("/api/deploy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			DeployID string            `json:"deploy_id"`
			Manifest map[string]string `json:"manifest"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if err := CommitDeploy(fs, store, req.DeployID, req.Manifest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("/", NewWebServer(fs, store))
	return mux
}

// DeployPipeline coordinates manifest creation, delta queries, missing blob uploads, and deploy activation.
func DeployPipeline(client *Client, fs FS, localDir string, deployID string) (*DeployStats, error) {
	manifest, err := BuildManifest(fs, localDir)
	if err != nil {
		return nil, err
	}

	hashes := make([]string, 0, len(manifest))
	hashToRel := make(map[string]string)
	for relPath, hash := range manifest {
		hashes = append(hashes, hash)
		hashToRel[hash] = relPath
	}

	missing, err := client.Check(hashes)
	if err != nil {
		return nil, err
	}

	cleanRoot := strings.TrimRight(strings.ReplaceAll(localDir, "\\", "/"), "/")
	for _, hash := range missing {
		relPath := hashToRel[hash]
		filePath := cleanRoot + relPath

		r, err := fs.Open(filePath)
		if err != nil {
			return nil, err
		}
		uploadErr := client.UploadBlob(hash, r)
		r.Close()
		if uploadErr != nil {
			return nil, uploadErr
		}
	}

	if err := client.Deploy(deployID, manifest); err != nil {
		return nil, err
	}

	return &DeployStats{
		DeployID:      deployID,
		TotalFiles:    len(manifest),
		UploadedFiles: len(missing),
	}, nil
}
