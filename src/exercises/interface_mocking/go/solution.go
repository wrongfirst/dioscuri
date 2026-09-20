package main

import (
	"bytes"
	"io"
	"os"
)

type ReadOpener interface {
	Open(name string) (io.ReadCloser, error)
}

type MemoryOpener struct {
	files map[string][]byte
}

func NewMemoryOpener() *MemoryOpener {
	return &MemoryOpener{
		files: make(map[string][]byte),
	}
}

func (m *MemoryOpener) AddFile(name string, data []byte) {
	m.files[name] = data
}

func (m *MemoryOpener) Open(name string) (io.ReadCloser, error) {
	data, ok := m.files[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

// ReadFileString opens and reads file content to string.
func ReadFileString(ro ReadOpener, name string) (string, error) {
	r, err := ro.Open(name)
	if err != nil {
		return "", err
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
