package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
)

var ErrHashMismatch = errors.New("hash mismatch")

type BlobSaver interface {
	SaveBlob(hash string, r io.Reader) error
}

type MockBlobSaver struct {
	Saved map[string][]byte
}

func NewMockBlobSaver() *MockBlobSaver {
	return &MockBlobSaver{Saved: make(map[string][]byte)}
}

func (m *MockBlobSaver) SaveBlob(hash string, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	m.Saved[hash] = data
	return nil
}

// IngestBlob verifies that payload matches claimedHash before persisting.
func IngestBlob(saver BlobSaver, claimedHash string, r io.Reader) error {
	var buf bytes.Buffer
	hasher := sha256.New()
	w := io.MultiWriter(&buf, hasher)

	if _, err := io.Copy(w, r); err != nil {
		return err
	}

	digest := hex.EncodeToString(hasher.Sum(nil))
	if digest != claimedHash {
		return ErrHashMismatch
	}

	return saver.SaveBlob(claimedHash, &buf)
}
