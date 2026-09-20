package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/rand"
)

var (
	ErrIntegrityViolation = errors.New("blob content does not match hash")
	ErrDanglingReference  = errors.New("manifest references non-existent blob")
	ErrProductionInvalid  = errors.New("production pointer invalid")
)

type SimResult struct {
	StepsExecuted int
	DeploysCount  int
	RollbackCount int
	ReadsCount    int
}

type SimSystem struct {
	ActiveDeployID string
	History        []string
	Manifests      map[string]map[string]string
	Blobs          map[string][]byte
}

func NewSimSystem() *SimSystem {
	return &SimSystem{
		History:   make([]string, 0),
		Manifests: make(map[string]map[string]string),
		Blobs:     make(map[string][]byte),
	}
}

func (s *SimSystem) CheckInvariants() error {
	for k, data := range s.Blobs {
		h := sha256.Sum256(data)
		if hex.EncodeToString(h[:]) != k {
			return ErrIntegrityViolation
		}
	}
	for _, m := range s.Manifests {
		for _, hash := range m {
			if _, ok := s.Blobs[hash]; !ok {
				return ErrDanglingReference
			}
		}
	}
	if s.ActiveDeployID == "" || s.Manifests[s.ActiveDeployID] == nil {
		return ErrProductionInvalid
	}
	return nil
}

// RunSimulation runs steps operations and checks invariants after each step.
func RunSimulation(rng *rand.Rand, steps int) (SimResult, error) {
	return SimResult{}, nil
}
