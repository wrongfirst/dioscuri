The simulation loop exercises the system by executing randomized operations (deploy, rollback, read) and asserting all correctness invariants after every step.

Running thousands of operations in memory exposes regressions before software reaches production.

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

In each step:
1. If no deployment exists yet, execute **Deploy**. Otherwise, select between **Deploy** (`0`), **Rollback** (`1`), and **Read** (`2`).
2. **Deploy**: Create a new deployment ID, register its manifest and blobs, set `ActiveDeployID`, and append to history.
3. **Rollback**: Pick a deployment from history and update `ActiveDeployID`.
4. **Read**: Verify that all files in the active manifest exist in blob storage.
5. After each step, verify that all 3 system invariants hold. If an invariant is violated, return the error.
