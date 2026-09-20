If garbage collection runs while a client is uploading files, newly uploaded blobs are not yet referenced by any committed manifest. A naive GC deletes them immediately, causing the upcoming deployment to fail.

To prevent this race condition, GC must enforce a quarantine period: unreferenced blobs created recently must be preserved.

## Task

Implement `RunQuarantinedGC`:

```go
type StoredBlob struct {
    Hash      string
    CreatedAt time.Time
}

type GCStore interface {
    ListBlobs() ([]StoredBlob, error)
    DeleteBlob(hash string) error
}

func RunQuarantinedGC(
    store GCStore,
    manifests []map[string]string,
    now time.Time,
    quarantine time.Duration,
) (int, error)
```

1. **Mark**: Collect all hashes referenced in `manifests` into a live set.
2. **Sweep with Quarantine**: For each blob:
   - If present in the live set, preserve it.
   - If unreferenced, compute its age (`now.Sub(b.CreatedAt)`).
   - Only delete the blob if its age exceeds `quarantine`.
3. Return the count of deleted blobs.
