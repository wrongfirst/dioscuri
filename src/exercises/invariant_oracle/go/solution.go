package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var (
	ErrIntegrityViolation = errors.New("blob content does not match hash")
	ErrDanglingReference  = errors.New("manifest references non-existent blob")
	ErrProductionInvalid  = errors.New("production pointer invalid")
)

type SystemState struct {
	ActiveDeployID string
	Manifests      map[string]map[string]string
	Blobs          map[string][]byte
}

func VerifyIntegrity(blobs map[string][]byte) error {
	for key, data := range blobs {
		h := sha256.Sum256(data)
		digest := hex.EncodeToString(h[:])
		if digest != key {
			return ErrIntegrityViolation
		}
	}
	return nil
}

func VerifyNoDangling(manifests map[string]map[string]string, blobs map[string][]byte) error {
	for _, m := range manifests {
		for _, hash := range m {
			if _, ok := blobs[hash]; !ok {
				return ErrDanglingReference
			}
		}
	}
	return nil
}

func VerifyProduction(state *SystemState) error {
	if state.ActiveDeployID == "" {
		return ErrProductionInvalid
	}
	if _, ok := state.Manifests[state.ActiveDeployID]; !ok {
		return ErrProductionInvalid
	}
	return nil
}
