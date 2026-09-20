package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"strings"
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
}

func NewMockFS() *MockFS {
	return &MockFS{Files: make(map[string][]byte)}
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

func (m *MockFS) Remove(p string) error {
	if _, ok := m.Files[p]; !ok {
		return os.ErrNotExist
	}
	delete(m.Files, p)
	return nil
}

func (m *MockFS) MkdirAll(p string) error { return nil }

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
	for f := range m.Files {
		if strings.HasPrefix(f, root) {
			if err := fn(f, false); err != nil {
				return err
			}
		}
	}
	return nil
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

func (s *Store) OpenBlob(hash string) (io.ReadCloser, error) {
	return s.fs.Open(s.BlobPath(hash))
}

type Oracle struct {
	fs          FS
	store       *Store
	manifestDir string
}

func NewOracle(fs FS, store *Store, manifestDir string) *Oracle {
	return &Oracle{
		fs:          fs,
		store:       store,
		manifestDir: manifestDir,
	}
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
	var ptr struct {
		DeployID string `json:"deploy_id"`
	}
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
