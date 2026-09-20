package main

import (
	"bytes"
	"io"
	"os"
	"strings"
)

type ReaderFS interface {
	Open(name string) (io.ReadCloser, error)
}

type MockFS struct {
	files map[string][]byte
}

func NewMockFS() *MockFS {
	return &MockFS{
		files: make(map[string][]byte),
	}
}

func (m *MockFS) AddFile(name string, data []byte) {
	m.files[name] = data
}

func (m *MockFS) Open(name string) (io.ReadCloser, error) {
	data, ok := m.files[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

// ReadFileString opens and streams file content to string.
func ReadFileString(fs ReaderFS, name string) (string, error) {
	r, err := fs.Open(name)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var sb strings.Builder
	if _, err := io.Copy(&sb, r); err != nil {
		return "", err
	}
	return sb.String(), nil
}
