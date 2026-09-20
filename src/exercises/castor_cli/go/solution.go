package main

import (
	"errors"
)

var ErrFileNotFound = errors.New("file not found in localFiles")

type ClientServerAPI interface {
	Check(hashes []string) (missing []string, err error)
	Upload(hash string, data []byte) error
	Deploy(deployID string, manifest map[string]string) error
}

type DeployStats struct {
	DeployID      string
	TotalFiles    int
	UploadedFiles int
}

type MockServerAPI struct {
	KnownBlobs   map[string]bool
	Uploaded     map[string][]byte
	DeployedID   string
	LastManifest map[string]string
}

func NewMockServerAPI(known ...string) *MockServerAPI {
	m := &MockServerAPI{
		KnownBlobs: make(map[string]bool),
		Uploaded:   make(map[string][]byte),
	}
	for _, k := range known {
		m.KnownBlobs[k] = true
	}
	return m
}

func (m *MockServerAPI) Check(hashes []string) ([]string, error) {
	missing := make([]string, 0)
	for _, h := range hashes {
		if !m.KnownBlobs[h] {
			missing = append(missing, h)
		}
	}
	return missing, nil
}

func (m *MockServerAPI) Upload(hash string, data []byte) error {
	m.Uploaded[hash] = data
	m.KnownBlobs[hash] = true
	return nil
}

func (m *MockServerAPI) Deploy(deployID string, manifest map[string]string) error {
	m.DeployedID = deployID
	m.LastManifest = manifest
	return nil
}

// RunDeploymentPipeline coordinates the check, upload, and deploy sequence.
func RunDeploymentPipeline(api ClientServerAPI, deployID string, localFiles map[string][]byte, manifest map[string]string) (DeployStats, error) {
	hashes := make([]string, 0, len(manifest))
	hashToPath := make(map[string]string)
	for path, hash := range manifest {
		hashes = append(hashes, hash)
		hashToPath[hash] = path
	}

	missing, err := api.Check(hashes)
	if err != nil {
		return DeployStats{}, err
	}

	for _, hash := range missing {
		p, ok := hashToPath[hash]
		if !ok {
			return DeployStats{}, ErrFileNotFound
		}
		data, ok := localFiles[p]
		if !ok {
			return DeployStats{}, ErrFileNotFound
		}
		if err := api.Upload(hash, data); err != nil {
			return DeployStats{}, err
		}
	}

	if err := api.Deploy(deployID, manifest); err != nil {
		return DeployStats{}, err
	}

	return DeployStats{
		DeployID:      deployID,
		TotalFiles:    len(manifest),
		UploadedFiles: len(missing),
	}, nil
}
