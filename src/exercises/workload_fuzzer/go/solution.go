package main

import (
	"fmt"
	"math/rand"
	"sort"
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

// MutateSite clones base and applies deterministic additions, modifications, and deletions.
func MutateSite(rng *rand.Rand, base Site, toAdd, toModify, toDelete int) Site {
	mutated := make(Site)
	for k, v := range base {
		mutated[k] = v
	}

	keys := make([]string, 0, len(base))
	for k := range base {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	deleted := 0
	for i := 0; i < toDelete && i < len(keys); i++ {
		delete(mutated, keys[i])
		deleted++
	}

	for i := deleted; i < deleted+toModify && i < len(keys); i++ {
		mutated[keys[i]] = fmt.Sprintf("modified_%d", rng.Intn(1000))
	}

	for i := 0; i < toAdd; i++ {
		path := fmt.Sprintf("/new_%d.html", i)
		content := fmt.Sprintf("new_content_%d", rng.Intn(1000))
		mutated[path] = content
	}

	return mutated
}
