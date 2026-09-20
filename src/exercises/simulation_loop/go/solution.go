package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	ErrIntegrityViolation = errors.New("blob content does not match hash")
	ErrDanglingReference  = errors.New("manifest references non-existent blob")
	ErrProductionInvalid  = errors.New("production pointer invalid")
)

type FS interface {
	Open(path string) (io.ReadCloser, error)
	Create(path string) (io.WriteCloser, error)
	Rename(oldPath, newPath string) error
	Remove(path string) error
	Exists(path string) bool
	MkdirAll(path string) error
	ReadDir(path string) ([]string, error)
	WalkDir(root string, fn func(path string, isDir bool) error) error
}

type MemoryFS struct {
	mu    sync.RWMutex
	files map[string][]byte
	dirs  map[string]bool
}

func NewMemoryFS() *MemoryFS {
	return &MemoryFS{
		files: make(map[string][]byte),
		dirs:  make(map[string]bool),
	}
}

type memWriter struct {
	buf  bytes.Buffer
	fs   *MemoryFS
	path string
}

func (w *memWriter) Write(p []byte) (int, error) { return w.buf.Write(p) }
func (w *memWriter) Close() error {
	w.fs.mu.Lock()
	defer w.fs.mu.Unlock()
	data := make([]byte, w.buf.Len())
	copy(data, w.buf.Bytes())
	w.fs.files[w.path] = data
	dir := path.Dir(w.path)
	for dir != "/" && dir != "." && dir != "" {
		w.fs.dirs[dir] = true
		dir = path.Dir(dir)
	}
	return nil
}

func (m *MemoryFS) Open(p string) (io.ReadCloser, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, ok := m.files[p]
	if !ok {
		return nil, os.ErrNotExist
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	return io.NopCloser(bytes.NewReader(cp)), nil
}

func (m *MemoryFS) Create(p string) (io.WriteCloser, error) {
	return &memWriter{fs: m, path: p}, nil
}

func (m *MemoryFS) Rename(old, new string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, ok := m.files[old]
	if !ok {
		return os.ErrNotExist
	}
	delete(m.files, old)
	m.files[new] = data
	return nil
}

func (m *MemoryFS) Remove(p string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.files, p)
	return nil
}

func (m *MemoryFS) Exists(p string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.files[p]
	return ok || m.dirs[p]
}

func (m *MemoryFS) MkdirAll(p string) error { return nil }

func (m *MemoryFS) ReadDir(dir string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cleanDir := strings.TrimRight(dir, "/")
	prefix := cleanDir + "/"
	seen := make(map[string]bool)
	var list []string
	for p := range m.files {
		if strings.HasPrefix(p, prefix) {
			rel := strings.TrimPrefix(p, prefix)
			parts := strings.Split(rel, "/")
			if len(parts) > 0 && !seen[parts[0]] {
				seen[parts[0]] = true
				list = append(list, parts[0])
			}
		}
	}
	sort.Strings(list)
	return list, nil
}

func (m *MemoryFS) WalkDir(root string, fn func(string, bool) error) error {
	return nil
}

type Store struct {
	fs      FS
	blobDir string
}

func NewStore(fs FS, blobDir string) *Store {
	return &Store{fs: fs, blobDir: blobDir}
}

func (s *Store) BlobPath(h string) string {
	return path.Join(s.blobDir, h)
}

func (s *Store) HasBlob(h string) bool {
	return s.fs.Exists(s.BlobPath(h))
}

func (s *Store) SaveBlob(h string, r io.Reader) error {
	if s.HasBlob(h) {
		return nil
	}
	w, err := s.fs.Create(s.BlobPath(h))
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, r)
	return err
}

func (s *Store) OpenBlob(h string) (io.ReadCloser, error) {
	return s.fs.Open(s.BlobPath(h))
}

type ActivePointer struct {
	DeployID string `json:"deploy_id"`
}

type Oracle struct {
	fs          FS
	store       *Store
	manifestDir string
}

func NewOracle(fs FS, store *Store, manifestDir string) *Oracle {
	return &Oracle{fs: fs, store: store, manifestDir: manifestDir}
}

func (o *Oracle) AssertBlobIntegrity() error {
	hashes, err := o.fs.ReadDir(o.store.blobDir)
	if err != nil {
		return err
	}
	for _, hash := range hashes {
		r, err := o.store.OpenBlob(hash)
		if err != nil {
			return err
		}
		hasher := sha256.New()
		_, copyErr := io.Copy(hasher, r)
		r.Close()
		if copyErr != nil {
			return copyErr
		}
		if hex.EncodeToString(hasher.Sum(nil)) != hash {
			return ErrIntegrityViolation
		}
	}
	return nil
}

func (o *Oracle) AssertProductionSafety() error {
	r, err := o.fs.Open("/production.json")
	if err != nil {
		return ErrProductionInvalid
	}
	var ptr ActivePointer
	err = json.NewDecoder(r).Decode(&ptr)
	r.Close()
	if err != nil || ptr.DeployID == "" {
		return ErrProductionInvalid
	}
	manifestPath := path.Join(o.manifestDir, ptr.DeployID+".json")
	if !o.fs.Exists(manifestPath) {
		return ErrProductionInvalid
	}
	return nil
}

func (o *Oracle) AssertNoDanglingReferences() error {
	manifestFiles, err := o.fs.ReadDir(o.manifestDir)
	if err != nil {
		return err
	}
	for _, mf := range manifestFiles {
		if !strings.HasSuffix(mf, ".json") {
			continue
		}
		r, err := o.fs.Open(path.Join(o.manifestDir, mf))
		if err != nil {
			return err
		}
		var manifest map[string]string
		err = json.NewDecoder(r).Decode(&manifest)
		r.Close()
		if err != nil {
			return err
		}
		for _, hash := range manifest {
			if !o.store.HasBlob(hash) {
				return ErrDanglingReference
			}
		}
	}
	return nil
}

func (o *Oracle) AssertAll() error {
	if err := o.AssertBlobIntegrity(); err != nil {
		return err
	}
	if err := o.AssertProductionSafety(); err != nil {
		return err
	}
	if err := o.AssertNoDanglingReferences(); err != nil {
		return err
	}
	return nil
}

type SimResult struct {
	StepsExecuted int
	DeploysCount  int
	RollbackCount int
	ReadsCount    int
}

// RunSimulation runs steps operations and checks invariants after each step.
func RunSimulation(rng *rand.Rand, steps int) (SimResult, error) {
	fs := NewMemoryFS()
	store := NewStore(fs, "/blobs")
	oracle := NewOracle(fs, store, "/manifests")

	var history []string
	res := SimResult{}

	for step := 0; step < steps; step++ {
		action := 0
		if len(history) > 0 {
			action = rng.Intn(3)
		}

		switch action {
		case 0:
			deployID := fmt.Sprintf("dep-%04d", len(history)+1)
			manifest := make(map[string]string)
			for i := 0; i < 2; i++ {
				content := fmt.Sprintf("content_%s_%d", deployID, i)
				h := sha256.Sum256([]byte(content))
				hStr := hex.EncodeToString(h[:])
				if err := store.SaveBlob(hStr, strings.NewReader(content)); err != nil {
					return res, err
				}
				manifest[fmt.Sprintf("/page_%d.html", i)] = hStr
			}

			// Save manifest
			mw, err := fs.Create(path.Join("/manifests", deployID+".json"))
			if err != nil {
				return res, err
			}
			if err := json.NewEncoder(mw).Encode(manifest); err != nil {
				mw.Close()
				return res, err
			}
			if err := mw.Close(); err != nil {
				return res, err
			}

			// Update production pointer atomically
			pw, err := fs.Create("/production.json.tmp")
			if err != nil {
				return res, err
			}
			if err := json.NewEncoder(pw).Encode(ActivePointer{DeployID: deployID}); err != nil {
				pw.Close()
				return res, err
			}
			if err := pw.Close(); err != nil {
				return res, err
			}
			if err := fs.Rename("/production.json.tmp", "/production.json"); err != nil {
				return res, err
			}

			history = append(history, deployID)
			res.DeploysCount++

		case 1:
			target := history[rng.Intn(len(history))]
			pw, err := fs.Create("/production.json.tmp")
			if err != nil {
				return res, err
			}
			if err := json.NewEncoder(pw).Encode(ActivePointer{DeployID: target}); err != nil {
				pw.Close()
				return res, err
			}
			if err := pw.Close(); err != nil {
				return res, err
			}
			if err := fs.Rename("/production.json.tmp", "/production.json"); err != nil {
				return res, err
			}
			res.RollbackCount++

		case 2:
			r, err := fs.Open("/production.json")
			if err != nil {
				return res, err
			}
			var ptr ActivePointer
			err = json.NewDecoder(r).Decode(&ptr)
			r.Close()
			if err != nil {
				return res, err
			}

			mr, err := fs.Open(path.Join("/manifests", ptr.DeployID+".json"))
			if err != nil {
				return res, err
			}
			var manifest map[string]string
			err = json.NewDecoder(mr).Decode(&manifest)
			mr.Close()
			if err != nil {
				return res, err
			}

			for _, hash := range manifest {
				br, err := store.OpenBlob(hash)
				if err != nil {
					return res, err
				}
				io.Copy(io.Discard, br)
				br.Close()
			}
			res.ReadsCount++
		}

		if err := oracle.AssertAll(); err != nil {
			return res, err
		}
		res.StepsExecuted++
	}

	return res, nil
}
