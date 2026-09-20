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
}

func (m *MemoryOpener) Open(name string) (io.ReadCloser, error) {
	return nil, os.ErrNotExist
}

// ReadFileString opens and reads file content to string.
func ReadFileString(ro ReadOpener, name string) (string, error) {
	return "", nil
}
