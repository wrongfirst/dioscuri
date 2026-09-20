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
	return h.now
}

func (h *SimHarness) Advance(d time.Duration) {
}

func (h *SimHarness) RandomInt(max int) int {
	return 0
}
