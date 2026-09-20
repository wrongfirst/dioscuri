package main

import (
	"errors"
	"sync"
)

var ErrDeployInProgress = errors.New("deployment currently locked")

type DeployLock interface {
	TryLock() bool
	Unlock()
}

type MockDeployLock struct {
	mu     sync.Mutex
	locked bool
}

func NewMockDeployLock() *MockDeployLock {
	return &MockDeployLock{}
}

func (m *MockDeployLock) TryLock() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.locked {
		return false
	}
	m.locked = true
	return true
}

func (m *MockDeployLock) Unlock() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.locked = false
}

type SafeDeployCoordinator struct {
	lock   DeployLock
	active string
}

func NewSafeDeployCoordinator(lock DeployLock) *SafeDeployCoordinator {
	return &SafeDeployCoordinator{lock: lock}
}

func (c *SafeDeployCoordinator) GetActive() string {
	return c.active
}

func (c *SafeDeployCoordinator) Deploy(deployID string, applyFn func() error) error {
	return nil
}
