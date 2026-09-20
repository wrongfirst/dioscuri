Activating a deployment requires verifying that all referenced blobs exist in storage before publishing. To prevent partial updates or downtime, production pointers must be updated atomically on disk.

Manifests are saved at `/manifests/<deploy_id>.json`. The active site pointer is stored at `/production.json`.

## Execution Flow

1. **Integrity Pre-check (`CommitDeploy`)**:
   - Verify every hash in `manifest` via `checker.HasBlob(hash)`.
   - If any hash is missing, abort and return `ErrMissingBlobs`.
2. **Persist Manifest**:
   - Open `/manifests/<deployID>.json` via `fs.Create`.
   - Encode `manifest` as JSON and close the stream.
3. **Atomic Pointer Swap**:
   - Write `ActivePointer{DeployID: deployID}` to `/production.json.tmp` via `fs.Create`.
   - Atomically rename `/production.json.tmp` to `/production.json` using `fs.Rename`.
4. **Rollback (`RollbackDeploy`)**:
   - Check if `/manifests/<targetDeployID>.json` exists in `fs`. If not, return `ErrDeployNotFound`.
   - Atomically write `ActivePointer{DeployID: targetDeployID}` to `/production.json.tmp` and rename to `/production.json`.

## Task

Implement `CommitDeploy` and `RollbackDeploy`:

```go
var (
    ErrMissingBlobs   = errors.New("missing blobs in store")
    ErrDeployNotFound = errors.New("deployment not found")
)

type FS interface {
    Exists(path string) bool
    Create(path string) (io.WriteCloser, error)
    Open(path string) (io.ReadCloser, error)
    Rename(oldPath, newPath string) error
}

type BlobChecker interface {
    HasBlob(hash string) bool
}

type Manifest map[string]string

type ActivePointer struct {
    DeployID string `json:"deploy_id"`
}

func CommitDeploy(fs FS, checker BlobChecker, deployID string, manifest Manifest) error
func RollbackDeploy(fs FS, targetDeployID string) error
```
