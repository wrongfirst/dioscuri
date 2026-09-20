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
	if f.diskFull {
		return ErrDiskFull
	}

	if f.byteBudget >= 0 {
		if f.byteBudget == 0 {
			return ErrPowerCut
		}
		if len(data) > f.byteBudget {
			partial := data[:f.byteBudget]
			cp := make([]byte, len(partial))
			copy(cp, partial)
			f.mem.Files[path] = cp
			f.byteBudget = 0
			return ErrPowerCut
		}
		f.byteBudget -= len(data)
	}

	cp := make([]byte, len(data))
	copy(cp, data)
	f.mem.Files[path] = cp
	return nil
}

func (f *FaultyFS) CorruptByte(path string, offset int) error {
	data, ok := f.mem.Files[path]
	if !ok {
		return os.ErrNotExist
	}
	if offset < 0 || offset >= len(data) {
		return errors.New("offset out of bounds")
	}
	data[offset] ^= 0xFF
	return nil
}
