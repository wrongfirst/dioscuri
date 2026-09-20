package main

import (
	"fmt"
	"math/rand"
	"sort"
)

type Site map[string]string

// GenerateSite creates fileCount pages deterministically from rng.
func GenerateSite(rng *rand.Rand, fileCount int) Site {
	return nil
}

// MutateSite clones base and applies deterministic additions, modifications, and deletions.
func MutateSite(rng *rand.Rand, base Site, toAdd, toModify, toDelete int) Site {
	return nil
}
