Deterministic simulation requires running thousands of filesystem operations per second entirely in memory.

A memory filesystem implements core file operations using in-memory byte slices.

## Task

Implement `MemoryFS`:

```go
type MemoryFS struct {
    files map[string][]byte
}

func NewMemoryFS() *MemoryFS
func (m *MemoryFS) WriteFile(path string, data []byte) error
func (m *MemoryFS) ReadFile(path string) ([]byte, error)
func (m *MemoryFS) Rename(oldPath, newPath string) error
func (m *MemoryFS) Remove(path string) error
func (m *MemoryFS) Exists(path string) bool
func (m *MemoryFS) ListDir(prefix string) []string
```

- `WriteFile`: Store data at `path`.
- `ReadFile`: Return data or `os.ErrNotExist` if missing.
- `Rename`: Move data from `oldPath` to `newPath`. Return `os.ErrNotExist` if `oldPath` does not exist.
- `Remove`: Delete file at `path`. Return `os.ErrNotExist` if missing.
- `Exists`: Return `true` if `path` exists.
- `ListDir`: Return all file paths starting with `prefix`.
