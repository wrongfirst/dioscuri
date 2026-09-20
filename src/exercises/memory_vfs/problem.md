A thread-safe in-memory virtual filesystem enables thousands of simulated filesystem operations per second without touching physical disk.

`MemoryFS` implements the canonical `FS` interface backed by in-memory byte buffers and protected by `sync.RWMutex`.

## Execution Flow

1. **Concurrency Control**: Protect internal file and directory maps using `sync.RWMutex` (`RLock`/`RUnlock` for read queries, `Lock`/`Unlock` for mutations).
2. **Streamed Writing (`Create`)**:
   - Return an `io.WriteCloser`.
   - On `Close()`, lock the filesystem, store the buffered bytes at `path`, and mark ancestor directory paths.
3. **Streamed Reading (`Open`)**:
   - Under read lock, retrieve byte slice.
   - Return `os.ErrNotExist` if missing, otherwise return `io.NopCloser(bytes.NewReader(data))`.
4. **Filesystem Operations**:
   - **`Exists(path)`**: Return `true` if path exists as a file or directory.
   - **`Rename(oldPath, newPath)`**: Move file data under lock. Return `os.ErrNotExist` if `oldPath` is not present.
   - **`Remove(path)`**: Delete file under lock. Return `os.ErrNotExist` if not present.
   - **`MkdirAll(path)`**: Register directory and parent paths.
   - **`ReadDir(dir)`**: Return sorted direct child names within `dir`.
   - **`WalkDir(root, fn)`**: Walk matching directories (`isDir: true`) and files (`isDir: false`) under `root`.

## Task

Implement `MemoryFS`:

```go
type FS interface {
    Open(path string) (io.ReadCloser, error)
    Create(path string) (io.WriteCloser, error)
    Rename(oldPath, newPath string) error
    Remove(path string) error
    Exists(path string) bool
    MkdirAll(path string) error
    ReadDir(path string) ([]string, error)
    WalkDir(root string, fn func(path string, isDir bool) error) error
}

func NewMemoryFS() *MemoryFS
```
