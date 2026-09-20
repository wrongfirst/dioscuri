When multiple deployments are triggered concurrently against the same environment, uncoordinated operations can interleave file writes or overwrite the active production pointer out of order.

A deployment coordinator serializes cutovers using non-blocking mutual exclusion (`TryLock`), rejecting concurrent deployment attempts while an existing operation is in progress.

## Execution Flow

1. **Acquire Lock**:
   - Attempt to acquire the mutex via `c.lock.TryLock()`.
   - If `TryLock()` returns `false`, return `ErrDeployInProgress` immediately without blocking.
2. **Release Guarantee**:
   - Ensure the lock is released on return via `defer c.lock.Unlock()`.
3. **Execute Deployment Hook**:
   - Execute `applyFn()`.
   - If `applyFn()` returns an error, return that error immediately without updating `c.active`.
4. **Update State**:
   - If `applyFn()` succeeds, set `c.active = deployID`.
   - Return `nil`.

## Task

Implement `SafeDeployCoordinator`:

```go
var ErrDeployInProgress = errors.New("deployment currently locked")

type DeployLock interface {
    TryLock() bool
    Unlock()
}

type SafeDeployCoordinator struct {
    lock   DeployLock
    active string
}

func (c *SafeDeployCoordinator) Deploy(deployID string, applyFn func() error) error
```
