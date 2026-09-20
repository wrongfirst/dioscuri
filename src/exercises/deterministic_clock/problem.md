Physical system time and concurrency are non-deterministic. A simulation framework requires full control over time and randomness from a single seed.

A deterministic clock advances only when stepped programmatically, ensuring two runs with the same seed follow the exact same trace.

## Task

Implement `SimHarness`:

```go
type SimHarness struct {
    now time.Time
    rng *rand.Rand
}

func NewSimHarness(seed int64, start time.Time) *SimHarness
func (h *SimHarness) Now() time.Time
func (h *SimHarness) Advance(d time.Duration)
func (h *SimHarness) RandomInt(max int) int
```

- `Now`: Return the current virtual time.
- `Advance(d)`: Step virtual time forward by duration `d`.
- `RandomInt(max)`: Return a pseudo-random integer in `[0, max)` using the internal PRNG.
- Two harnesses initialized with the same seed and start time must yield identical sequences.
