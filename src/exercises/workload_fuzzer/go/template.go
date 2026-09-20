package main

import (
	"fmt"
	"math/rand"
)

type Site map[string]string

// GenerateSite creates fileCount pages deterministically from rng.
func GenerateSite(rng *rand.Rand, fileCount int) Site {
	return nil
}

// MutateSite copies base and adds newFiles pages deterministically.
func MutateSite(rng *rand.Rand, base Site, newFiles int) Site {
	return nil
}
