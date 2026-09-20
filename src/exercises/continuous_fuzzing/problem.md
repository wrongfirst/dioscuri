Continuous chaos testing combines fault injection, concurrent actions, and garbage collection into a randomized harness to verify that the hardened system never violates invariants.

Under power cuts, in-flight uploads, and GC cycles, all system invariants must hold.

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

In each step:
1. Select a randomized action:
   - **Upload (`0`)**: Compute SHA-256, write safely via staging, and commit blob.
   - **Torn Write Crash (`1`)**: Reject corrupted payload and ensure no partial blob enters storage.
   - **Quarantined GC (`2`)**: Run GC with quarantine; verify active and recent blobs are preserved.
   - **Deploy (`3`)**: Commit manifest using atomic pointer replacement.
2. Assert that all blobs match their SHA-256 keys and no dangling references exist.
3. Return `ChaosResult`.
