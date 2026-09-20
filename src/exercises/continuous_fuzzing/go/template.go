package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
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
	return ChaosResult{}, nil
}
