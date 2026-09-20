A Content-Addressable Storage (CAS) store organizes immutable files identified solely by their content hash. Blobs are persisted at `<blobDir>/<hash>`.

Because objects are immutable and content-keyed, storing a blob that already exists is an idempotent no-op.

## Architecture & Flow

1. **Path Construction (`BlobPath`)**:
   - Joins `blobDir` and `hash` using `path.Join`.
2. **Existence Query (`HasBlob`)**:
   - Returns `s.fs.Exists(s.BlobPath(hash))`.
3. **Idempotent Save (`SaveBlob`)**:
   - If `s.HasBlob(hash)` is true, return `nil` immediately without performing any I/O.
   - Call `s.fs.Create(s.BlobPath(hash))` to open a writable stream.
   - Ensure the stream is closed on return (`defer w.Close()`).
   - Stream content from `r` to `w` using `io.Copy(w, r)`.
4. **Retrieval (`OpenBlob`)**:
   - Open and return the readable stream using `s.fs.Open(s.BlobPath(hash))`.

## Task

Implement `Store`:

```go
type FS interface {
    Exists(path string) bool
    Create(path string) (io.WriteCloser, error)
    Open(path string) (io.ReadCloser, error)
}

type Store struct {
    fs      FS
    blobDir string
}

func NewStore(fs FS, blobDir string) *Store
func (s *Store) BlobPath(hash string) string
func (s *Store) HasBlob(hash string) bool
func (s *Store) SaveBlob(hash string, r io.Reader) error
func (s *Store) OpenBlob(hash string) (io.ReadCloser, error)
```
