package main

import (
	"errors"
	"os"
)

var (
	ErrPowerCut = errors.New("power cut")
	ErrDiskFull = errors.New("disk full")
)

type MemoryFS struct {
	Files map[string][]byte
}

func NewMemoryFS() *MemoryFS {
	return &MemoryFS{Files: make(map[string][]byte)}
}

type FaultyFS struct {
	mem        *MemoryFS
	diskFull   bool
	byteBudget int // -1 if disabled
}

func NewFaultyFS(mem *MemoryFS) *FaultyFS {
	return &FaultyFS{
		mem:        mem,
		diskFull:   false,
		byteBudget: -1,
	}
}

func (f *FaultyFS) SetDiskFull(full bool) {
	f.diskFull = full
}

func (f *FaultyFS) CrashAfter(bytes int) {
	f.byteBudget = bytes
}

func (f *FaultyFS) WriteFile(path string, data []byte) error {
	return nil
}

func (f *FaultyFS) CorruptByte(path string, offset int) error {
	return nil
}
