Testing crash safety and recovery requires simulating hardware failures: mid-write power loss, storage capacity exhaustion, and bit-level corruption.

`FaultyFS` is a decorator around the canonical `FS` interface that injects faults into streaming writes.

## Execution Flow

1. **Decorator Pattern**: Wrap an underlying `FS` instance.
2. **Fault Triggers**:
   - `SetDiskFull(full bool)`: When enabled, `Create()` and `Write()` calls return `ErrDiskFull`.
   - `CrashAfter(bytes int)`: Sets a write budget. When a `Write(p)` call exceeds the remaining budget, write only the allowable bytes, set the budget to `0`, and return `ErrPowerCut`.
   - `CorruptByte(path, offset)`: Open the file, invert the byte at `offset` (`^= 0xFF`), and save it back. Return an error if out of bounds.
3. **Stream Interception**:
   - Return a `faultyWriter` wrapping `io.WriteCloser` from `base.Create(p)`.
   - Intercept each `Write(p)` to enforce budget limits and disk-full conditions.
4. **Delegation**: Forward `Open`, `Rename`, `Remove`, `Exists`, `MkdirAll`, `ReadDir`, and `WalkDir` directly to `base`.

## Task

Implement `FaultyFS` and `faultyWriter`:

```go
var (
    ErrPowerCut = errors.New("power cut")
    ErrDiskFull = errors.New("disk full")
)

type FaultyFS struct {
    base       FS
    diskFull   bool
    byteBudget int // -1 if disabled
}

func NewFaultyFS(base FS) *FaultyFS
func (f *FaultyFS) SetDiskFull(full bool)
func (f *FaultyFS) CrashAfter(bytes int)
func (f *FaultyFS) CorruptByte(path string, offset int) error
```
