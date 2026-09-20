package main

import (
	"math/rand"
	"time"
)

type SimHarness struct {
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
	h.now = h.now.Add(d)
}

func (h *SimHarness) RandomInt(max int) int {
	if max <= 0 {
		return 0
	}
	return h.rng.Intn(max)
}
