When two deployments are triggered simultaneously, uncoordinated operations can interleave writes or create conflicting production states.

A deployment coordinator serializes cutovers using mutual exclusion, rejecting concurrent attempts while an operation is in progress.

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

func NewSafeDeployCoordinator(lock DeployLock) *SafeDeployCoordinator
func (c *SafeDeployCoordinator) GetActive() string
func (c *SafeDeployCoordinator) Deploy(deployID string, applyFn func() error) error
```

1. Call `c.lock.TryLock()`. If locked, return `ErrDeployInProgress`.
2. Ensure the lock is released on exit (`defer c.lock.Unlock()`).
3. Call `applyFn()`. If it errors, return the error without modifying `active`.
4. If `applyFn()` succeeds, set `active = deployID`.
