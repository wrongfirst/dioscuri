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
	return nil
}

func VerifyNoDangling(manifests map[string]map[string]string, blobs map[string][]byte) error {
	return nil
}

func VerifyProduction(state *SystemState) error {
	return nil
}
