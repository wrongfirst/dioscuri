When deployments are replaced or retired, some stored blobs are no longer referenced by any active manifest. A mark-and-sweep garbage collector scans the filesystem to reclaim storage by deleting unreferenced blobs while preserving shared assets.

## Collector Execution Flow

1. **Mark Phase**:
   - List all files in `manifestDir` via `fs.ReadDir(manifestDir)`.
   - For every entry ending with `.json`:
     - Open and decode the manifest file into `map[string]string`.
     - Add every referenced hash value to a `live` set (`map[string]bool`).
2. **Sweep Phase**:
   - List all files in `blobDir` via `fs.ReadDir(blobDir)`.
   - For each blob:
     - Check if it is present in the `live` set.
     - If not live: remove it using `fs.Remove(path.Join(blobDir, hash))` and increment `deletedCount`.
3. **Return**:
   - Return `deletedCount` and any error encountered.

## Task

Implement `RunGC`:

```go
type FS interface {
    ReadDir(dir string) ([]string, error)
    Open(path string) (io.ReadCloser, error)
    Remove(path string) error
}

func RunGC(fs FS, manifestDir string, blobDir string) (int, error)
```
