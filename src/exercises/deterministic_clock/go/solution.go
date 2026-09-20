package main

import (
	"math/rand"
	"sync"
	"time"
)

type SimHarness struct {
	mu  sync.Mutex
	now time.Time
	rng *rand.Rand
}

func NewSimHarness(seed int64, start time.Time) *SimHarness {
	return &SimHarness{
		now: start,
		rng: rand.New(rand.NewSource(seed)),
	}
}

func (h *SimHarness) Now() time.Time {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.now
}

func (h *SimHarness) Advance(d time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.now = h.now.Add(d)
}

func (h *SimHarness) RandomInt(max int) int {
	h.mu.Lock()
	defer h.mu.Unlock()
	if max <= 0 {
		return 0
	}
	return h.rng.Intn(max)
}
