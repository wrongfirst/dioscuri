Calling OS filesystem functions directly hardwires code to physical disk, preventing in-memory testing and deterministic simulation. Defining an interface contract decouples storage callers from the underlying filesystem.

## Architecture & Flow

1. **Interface Contract (`ReaderFS`)**:
   Requires only opening a path to obtain a readable, closeable stream:
   ```go
   type ReaderFS interface {
       Open(name string) (io.ReadCloser, error)
   }
   ```

2. **In-Memory Double (`MockFS`)**:
   - Holds an internal `map[string][]byte`.
   - `AddFile(name, data)` stores the byte slice.
   - `Open(name)` checks for the key. If found, wraps the data in `io.NopCloser(bytes.NewReader(data))`. If missing, returns `os.ErrNotExist`.

3. **Reader Helper (`ReadFileString`)**:
   - Opens `name` via `fs.Open(name)`.
   - Schedules closure immediately: `defer r.Close()`.
   - Streams bytes into a `strings.Builder` using `io.Copy(&sb, r)` to avoid intermediate slice allocations.
   - Returns the assembled string.

## Task

Implement `MockFS`, its methods, and `ReadFileString`:

```go
func NewMockFS() *MockFS
func (m *MockFS) AddFile(name string, data []byte)
func (m *MockFS) Open(name string) (io.ReadCloser, error)
func ReadFileString(fs ReaderFS, name string) (string, error)
```
