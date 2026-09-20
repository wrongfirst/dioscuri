A content-addressable blob store manages immutable files keyed solely by their content hash. Blobs are stored at `<blobDir>/<hash>`.

Because files are addressed by content hash, saving a blob that already exists is an idempotent no-op.

## Task

Implement `Store`:

```go
type BlobFS interface {
    Exists(path string) bool
    Create(path string) (io.WriteCloser, error)
    Open(path string) (io.ReadCloser, error)
}

type Store struct {
    fs      BlobFS
    blobDir string
}

func NewStore(fs BlobFS, blobDir string) *Store
func (s *Store) HasBlob(hash string) bool
func (s *Store) SaveBlob(hash string, r io.Reader) error
func (s *Store) OpenBlob(hash string) (io.ReadCloser, error)
```

- `HasBlob(hash)`: Return `true` if the blob exists in `fs`.
- `SaveBlob(hash, r)`: If the blob already exists, return `nil` without writing. Otherwise, create `<blobDir>/<hash>`, stream from `r`, and close the file handle.
- `OpenBlob(hash)`: Open and return the blob reader.
