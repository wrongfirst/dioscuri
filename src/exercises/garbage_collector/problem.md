When deployments are removed, some blobs in storage are no longer referenced by any manifest.

A mark-and-sweep garbage collector deletes these orphaned blobs while preserving blobs shared by active deployments.

## Task

Implement `RunGC`:

```go
type GCFS interface {
    ListBlobs() ([]string, error)
    DeleteBlob(hash string) error
}

func RunGC(fs GCFS, activeManifests []map[string]string) (int, error)
```

1. **Mark**: Iterate through all `activeManifests` and add every referenced hash to a live set.
2. **Sweep**: List all existing blobs via `fs.ListBlobs()`. For any blob not in the live set, call `fs.DeleteBlob(hash)`.
3. Return the total count of deleted blobs.
