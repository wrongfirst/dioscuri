Calling OS filesystem operations directly couples code to physical storage, making simulation and hardware fault injection difficult.

Abstracting file I/O behind an interface allows swapping between the real filesystem and an in-memory test double.

## Task

Implement the in-memory double and reader helper:

1. **`ReadOpener` interface**:
   ```go
   type ReadOpener interface {
       Open(name string) (io.ReadCloser, error)
   }
   ```

2. **`MemoryOpener` struct**:
   - `NewMemoryOpener() *MemoryOpener`: Initializes the struct.
   - `AddFile(name string, data []byte)`: Stores the byte slice in memory.
   - `Open(name string) (io.ReadCloser, error)`: Returns an `io.ReadCloser` over the data. If the file is missing, returns `os.ErrNotExist`.

3. **`ReadFileString(ro ReadOpener, name string) (string, error)`**:
   - Opens the file with `ro.Open(name)`.
   - Ensures the file is closed after reading (`defer`).
   - Reads the full content and returns it as a string.
