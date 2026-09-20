package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// HashReader returns the hex SHA-256 digest of r.
func HashReader(r io.Reader) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
