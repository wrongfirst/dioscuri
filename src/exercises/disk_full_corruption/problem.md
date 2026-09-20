Overwriting critical state files directly in-place introduces severe vulnerability: if the underlying storage fills up mid-write (`ErrDiskFull`), the existing file is truncated or corrupted, bringing down the entire production routing layer.

To guarantee crash resilience against disk exhaustion, write the updated state to a temporary staging file first and atomically swap it into place via `fs.Rename`.

## Execution Flow

1. Define the staging path `tmpPath := "/production.json.tmp"` and target path `targetPath := "/production.json"`.
2. Create the staging file via `w, err := fs.Create(tmpPath)`. If this fails (e.g., storage full), return the error immediately without modifying the existing `/production.json`.
3. Encode the active pointer JSON `{"deploy_id": deployID}` into `w` and close the writer. If encoding or closing fails, call `fs.Remove(tmpPath)` and return the error.
4. Atomically swap the new state into place via `fs.Rename(tmpPath, targetPath)`.
5. Return `nil`.

## Task

Implement `UpdateProductionPointer`:

```go
var ErrDiskFull = errors.New("disk full")

type ActivePointer struct {
    DeployID string `json:"deploy_id"`
}

func UpdateProductionPointer(fs FS, deployID string) error
```
