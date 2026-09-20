If garbage collection executes while a client is in the middle of uploading files, newly arrived blobs are not yet referenced by any committed manifest. A naive mark-and-sweep GC immediately purges them as orphaned, causing the upcoming deployment commit to fail with dangling references.

To protect in-flight uploads, the garbage collector must enforce a quarantine period: unreferenced blobs created within the quarantine window are preserved.

## Execution Flow

1. **Mark Phase**:
   - Traverse all provided manifests.
   - Insert every referenced hash into a `live` set (`map[string]bool`).
2. **Sweep Phase with Quarantine**:
   - Query all blobs and creation timestamps via `store.ListBlobs()`.
   - For each blob:
     - If referenced in `live`: skip deletion.
     - If unreferenced, compute its age: `now.Sub(blob.CreatedAt)`.
     - If `age > quarantine`: delete it via `store.DeleteBlob(blob.Hash)` and increment the deleted counter.
     - If `age <= quarantine`: preserve the blob (in-flight upload protection).
3. Return `(deletedCount, nil)`.

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
