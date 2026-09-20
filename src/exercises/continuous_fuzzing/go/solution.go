package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var (
	ErrInvariantFailed = errors.New("system invariant violated")
)

type ChaosResult struct {
	StepsCompleted int
	Uploads        int
	Crashes        int
	CyclesGC       int
}

type HardenedSystem struct {
	blobs     map[string][]byte
	blobTimes map[string]time.Time
	manifests map[string]map[string]string
	activeID  string
	now       time.Time
}

func NewHardenedSystem(start time.Time) *HardenedSystem {
	return &HardenedSystem{
		blobs:     make(map[string][]byte),
		blobTimes: make(map[string]time.Time),
		manifests: make(map[string]map[string]string),
		now:       start,
	}
}

func (s *HardenedSystem) CheckInvariants() error {
	for hash, data := range s.blobs {
		h := sha256.Sum256(data)
		if hex.EncodeToString(h[:]) != hash {
			return ErrInvariantFailed
		}
	}
	for _, m := range s.manifests {
		for _, hash := range m {
			if _, ok := s.blobs[hash]; !ok {
				return ErrInvariantFailed
			}
		}
	}
	return nil
}

// RunChaosHarness drives randomized hardened chaos steps and verifies invariants.
func RunChaosHarness(seed int64, iterations int) (ChaosResult, error) {
	rng := rand.New(rand.NewSource(seed))
	startTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	sys := NewHardenedSystem(startTime)
	res := ChaosResult{}

	for step := 0; step < iterations; step++ {
		sys.now = sys.now.Add(time.Minute)
		action := rng.Intn(4)

		switch action {
		case 0:
			data := []byte(fmt.Sprintf("payload_%d_%d", step, rng.Intn(1000)))
			h := sha256.Sum256(data)
			hashStr := hex.EncodeToString(h[:])
			sys.blobs[hashStr] = data
			sys.blobTimes[hashStr] = sys.now
			res.Uploads++

		case 1:
			res.Crashes++

		case 2:
			quarantine := 10 * time.Minute
			live := make(map[string]bool)
			for _, m := range sys.manifests {
				for _, h := range m {
					live[h] = true
				}
			}
			for h, t := range sys.blobTimes {
				if live[h] {
					continue
				}
				if sys.now.Sub(t) > quarantine {
					delete(sys.blobs, h)
					delete(sys.blobTimes, h)
				}
			}
			res.CyclesGC++

		case 3:
			deployID := fmt.Sprintf("deploy_%d", step)
			manifest := make(map[string]string)
			data := []byte(fmt.Sprintf("page_%s", deployID))
			h := sha256.Sum256(data)
			hashStr := hex.EncodeToString(h[:])
			sys.blobs[hashStr] = data
			sys.blobTimes[hashStr] = sys.now
			manifest["/index.html"] = hashStr

			sys.manifests[deployID] = manifest
			sys.activeID = deployID
		}

		if err := sys.CheckInvariants(); err != nil {
			return res, err
		}
		res.StepsCompleted++
	}

	return res, nil
}
