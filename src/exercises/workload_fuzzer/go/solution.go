package main

import (
	"fmt"
	"math/rand"
)

type Site map[string]string

// GenerateSite creates fileCount pages deterministically from rng.
func GenerateSite(rng *rand.Rand, fileCount int) Site {
	site := make(Site)
	for i := 0; i < fileCount; i++ {
		path := fmt.Sprintf("/page_%d.html", i)
		content := fmt.Sprintf("content_%d_%d", i, rng.Intn(1000))
		site[path] = content
	}
	return site
}

// MutateSite copies base and adds newFiles pages deterministically.
func MutateSite(rng *rand.Rand, base Site, newFiles int) Site {
	site := make(Site)
	for k, v := range base {
		site[k] = v
	}
	for i := 0; i < newFiles; i++ {
		path := fmt.Sprintf("/new_%d.html", i)
		content := fmt.Sprintf("new_content_%d", rng.Intn(1000))
		site[path] = content
	}
	return site
}
