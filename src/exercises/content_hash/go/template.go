package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// HashReader returns the hex SHA-256 digest of r.
func HashReader(r io.Reader) (string, error) {
	return "", nil
}
