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
}

func (m *MockFS) Open(name string) (io.ReadCloser, error) {
	return nil, os.ErrNotExist
}

// ReadFileString opens and streams file content to string.
func ReadFileString(fs ReaderFS, name string) (string, error) {
	return "", nil
}
