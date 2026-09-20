package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
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
	sys := NewSimSystem()
	res := SimResult{}

	for step := 0; step < steps; step++ {
		action := 0
		if len(sys.History) > 0 {
			action = rng.Intn(3)
		}

		switch action {
		case 0:
			deployID := fmt.Sprintf("deploy_%d", len(sys.History)+1)
			manifest := make(map[string]string)
			for i := 0; i < 2; i++ {
				path := fmt.Sprintf("/page_%d.html", i)
				body := []byte(fmt.Sprintf("body_%s_%d", deployID, i))
				h := sha256.Sum256(body)
				hashStr := hex.EncodeToString(h[:])
				sys.Blobs[hashStr] = body
				manifest[path] = hashStr
			}
			sys.Manifests[deployID] = manifest
			sys.ActiveDeployID = deployID
			sys.History = append(sys.History, deployID)
			res.DeploysCount++

		case 1:
			target := sys.History[rng.Intn(len(sys.History))]
			sys.ActiveDeployID = target
			res.RollbackCount++

		case 2:
			activeManifest := sys.Manifests[sys.ActiveDeployID]
			for _, hash := range activeManifest {
				if _, ok := sys.Blobs[hash]; !ok {
					return res, ErrDanglingReference
				}
			}
			res.ReadsCount++
		}

		if err := sys.CheckInvariants(); err != nil {
			return res, err
		}
		res.StepsExecuted++
	}

	return res, nil
}
