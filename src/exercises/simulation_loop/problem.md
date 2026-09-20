The simulation loop drives the system through hundreds of randomized operations (deployments, rollbacks, reads) against an in-memory virtual filesystem, asserting all system invariants after every state transition.

Running thousands of operations entirely in RAM exposes data corruption, broken pointers, and state leakage instantly.

## Code Self-Review: Concurrency in `MemoryFS`

Before assembling the simulation loop, review the pre-populated `MemoryFS` below:
- Where could race conditions occur if concurrent simulated workers (e.g., background garbage collection and simultaneous deployments) access the virtual filesystem at the same time?
- Verify that every inspection method (`Exists`, `Open`, `ReadDir`, `WalkDir`) acquires `RLock()` and releases it via `RUnlock()`, and every mutating method (`Create`, `Rename`, `Remove`, `MkdirAll`) holds `Lock()`.

## Execution Flow

In `RunSimulation(rng *rand.Rand, steps int) (SimResult, error)`:

1. **Environment Setup**:
   - Instantiate `fs := NewMemoryFS()`.
   - Instantiate `store := NewStore(fs, "/blobs")`.
   - Instantiate `oracle := NewOracle(fs, store, "/manifests")`.
   - Initialize `history := make([]string, 0)`.
2. **Step Execution**:
   - For each step from `0` to `steps - 1`:
     - If `len(history) == 0`: Action is **Deploy** (`0`).
     - Otherwise: Select `action := rng.Intn(3)` (0: Deploy, 1: Rollback, 2: Read).
   - **Action 0 (Deploy)**:
     - Generate deployment ID: `fmt.Sprintf("dep-%04d", len(history)+1)`.
     - Create 2 pages with deterministic content, compute SHA-256 hashes, store each blob via `store.SaveBlob(hash, reader)`, and construct `manifest[path] = hash`.
     - Write manifest JSON to `/manifests/<deployID>.json`.
     - Atomically update `/production.json` via temporary file rename (`/production.json.tmp` $\to$ `/production.json`).
     - Append deployment ID to `history` and increment `res.DeploysCount`.
   - **Action 1 (Rollback)**:
     - Select target deployment from `history`: `targetID := history[rng.Intn(len(history))]`.
     - Atomically update `/production.json` to point to `targetID`.
     - Increment `res.RollbackCount`.
   - **Action 2 (Read)**:
     - Read active deployment pointer from `/production.json`.
     - Open and decode active manifest from `/manifests/<activeID>.json`.
     - Stream each referenced blob via `store.OpenBlob(hash)` and verify content matches.
     - Increment `res.ReadsCount`.
3. **Invariant Assertion**:
   - After each step, execute `oracle.AssertAll()`. If an invariant fails, return `(res, err)` immediately.
   - Increment `res.StepsExecuted`.
4. Return `(res, nil)` upon completing all steps.

## Task

Implement `RunSimulation`:

```go
type SimResult struct {
    StepsExecuted int
    DeploysCount  int
    RollbackCount int
    ReadsCount    int
}

func RunSimulation(rng *rand.Rand, steps int) (SimResult, error)
```

## Run Locally

In your local repository scaffold (`dioscuri.zip`), you can run an aggressive 50,000-step simulation in RAM in less than a second:

```bash
go test -v -run TestSimulationLoop ./pkg/crucible
```
