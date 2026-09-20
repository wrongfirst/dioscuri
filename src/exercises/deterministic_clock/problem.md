Distributed system simulation requires complete control over time progression and pseudo-random decisions so race conditions can be reproduced identically from a seed.

`SimHarness` provides a thread-safe deterministic clock and pseudo-random number generator.

## Execution Flow

1. **State Initialization**: Initialize virtual time to `start` and instantiate an isolated PRNG using `rand.New(rand.NewSource(seed))`.
2. **Virtual Time Progression**:
   - `Now()` returns the current simulated `time.Time`.
   - `Advance(d)` moves virtual time forward by `d`. Simulated time never advances spontaneously.
3. **Reproducible Randomness**:
   - `RandomInt(max)` produces an integer in $[0, max)$ using the internal PRNG. If `max <= 0`, return `0`.
   - Two harnesses created with the same seed and start time must produce identical sequences.
4. **Concurrency Safety**: Protect shared state using `sync.Mutex` so concurrent workers (such as simultaneous deployment pipelines and garbage collection workers) can query time and randomness safely.

## Task

Implement `SimHarness`:

```go
type SimHarness struct {
    mu  sync.Mutex
    now time.Time
    rng *rand.Rand
}

func NewSimHarness(seed int64, start time.Time) *SimHarness
func (h *SimHarness) Now() time.Time
func (h *SimHarness) Advance(d time.Duration)
func (h *SimHarness) RandomInt(max int) int
```
