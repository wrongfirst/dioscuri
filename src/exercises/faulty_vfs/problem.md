Testing crash-safety requires simulating hardware faults: sudden power cuts mid-write, disk capacity limits, and storage corruption.

A fault-injecting filesystem wrapper intercepts operations to simulate torn writes and disk failures.

## Task

Implement `FaultyFS`:

```go
var (
    ErrPowerCut = errors.New("power cut")
    ErrDiskFull = errors.New("disk full")
)

type FaultyFS struct {
    mem         *MemoryFS
    diskFull    bool
    byteBudget  int // -1 if disabled
}

func NewFaultyFS(mem *MemoryFS) *FaultyFS
func (f *FaultyFS) SetDiskFull(full bool)
func (f *FaultyFS) CrashAfter(bytes int)
func (f *FaultyFS) WriteFile(path string, data []byte) error
func (f *FaultyFS) CorruptByte(path string, offset int) error
```

- `SetDiskFull(true)`: Future writes return `ErrDiskFull` without saving data.
- `CrashAfter(n)`: Sets a budget of `n` total bytes. If a write exceeds the remaining budget, write only the allowable bytes to `mem` and return `ErrPowerCut`.
- `CorruptByte(path, offset)`: Inverts the byte at `offset` (`^= 0xFF`) to simulate corruption.
