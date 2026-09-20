package main

import (
	"errors"
)

var (
	ErrMissingBlobs   = errors.New("missing blobs in store")
	ErrDeployNotFound = errors.New("deployment not found")
)

type BlobChecker interface {
	HasBlob(hash string) bool
}

type DeploymentState struct {
	ActiveDeployID string
	Manifests      map[string]map[string]string
}

func NewDeploymentState() *DeploymentState {
	return &DeploymentState{
		Manifests: make(map[string]map[string]string),
	}
}

type MockChecker struct {
	blobs map[string]bool
}

func NewMockChecker(hashes ...string) *MockChecker {
	m := &MockChecker{blobs: make(map[string]bool)}
	for _, h := range hashes {
		m.blobs[h] = true
	}
	return m
}

func (m *MockChecker) HasBlob(h string) bool {
	return m.blobs[h]
}

// CommitDeploy verifies manifest hashes, stores manifest, and points ActiveDeployID to deployID.
func CommitDeploy(state *DeploymentState, checker BlobChecker, deployID string, manifest map[string]string) error {
	for _, hash := range manifest {
		if !checker.HasBlob(hash) {
			return ErrMissingBlobs
		}
	}
	state.Manifests[deployID] = manifest
	state.ActiveDeployID = deployID
	return nil
}

// RollbackDeploy switches ActiveDeployID to targetDeployID if it exists.
func RollbackDeploy(state *DeploymentState, targetDeployID string) error {
	if _, ok := state.Manifests[targetDeployID]; !ok {
		return ErrDeployNotFound
	}
	state.ActiveDeployID = targetDeployID
	return nil
}
