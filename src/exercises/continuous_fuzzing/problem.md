Continuous chaos testing subjects the system to randomized permutations of concurrent uploads, simulated mid-stream power cuts, quarantined garbage collection cycles, and deployment cutovers.

By driving the hardened engine through hundreds of fault-injected operations in memory, we verify mathematically that all system invariants hold with zero data loss, zero corrupted blobs, and zero broken production pointers.

## Execution Flow

In `RunChaosHarness(seed int64, iterations int) (ChaosResult, error)`:

1. **Setup**:
   - Seed a deterministic PRNG using `seed`.
   - Initialize the virtual filesystem, content-addressed store, and invariant oracle.
   - Advance virtual time by 1 minute on each iteration.
2. **Chaos State Transitions**:
   - Randomly select one of four actions:
     - **Action 0 (Safe Upload)**: Stream a newly generated payload through staging with SHA-256 verification and commit to blob storage. Record its timestamp and increment `res.Uploads`.
     - **Action 1 (Torn-Write Crash)**: Inject a mid-stream power cut fault during staging. Verify that incomplete data is cleaned up, the corrupt blob is never committed, and increment `res.Crashes`.
     - **Action 2 (Quarantined GC)**: Run garbage collection. Sweep unreferenced blobs older than the quarantine threshold while strictly preserving in-flight blobs and all active manifests. Increment `res.CyclesGC`.
     - **Action 3 (Atomic Deploy)**: Commit a new manifest and atomically update `/production.json`.
3. **Invariant Assertion**:
   - Assert all invariants after every transition:
     - Content-addressed integrity: Every blob matches its SHA-256 hash.
     - Production safety: `/production.json` points to a valid, existing manifest.
     - Referential integrity: All manifests reference existing blobs.
   - Increment `res.StepsCompleted`.
4. Return `(res, nil)`.

## Task

Implement `RunChaosHarness`:

```go
type ChaosResult struct {
    StepsCompleted int
    Uploads        int
    Crashes        int
    CyclesGC       int
}

func RunChaosHarness(seed int64, iterations int) (ChaosResult, error)
```

## Run Locally: Native Go Tooling

In your local repository scaffold (`dioscuri.zip`), you can leverage Go's native toolchain for exhaustive runtime verification:

```bash
# 1. Concurrency Race Detection (ThreadSanitizer)
go test -race ./pkg/server/...

# 2. Native Coverage-Guided Fuzzing (libFuzzer engine)
go test -fuzz=FuzzPathNormalization -fuzztime=30s ./pkg/manifest
```
