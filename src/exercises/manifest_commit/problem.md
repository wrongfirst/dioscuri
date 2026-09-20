Deploying a manifest requires verifying that every referenced content hash exists in storage before activating the deployment.

Once verified, the manifest is saved under its deployment ID and the active production pointer is updated. A rollback switches the active pointer back to an earlier deployment ID.

## Task

Implement `CommitDeploy` and `RollbackDeploy`:

```go
var (
    ErrMissingBlobs   = errors.New("missing blobs in store")
    ErrDeployNotFound = errors.New("deployment not found")
)

type BlobChecker interface {
    HasBlob(hash string) bool
}

type DeploymentState struct {
    ActiveDeployID string
    Manifests      map[string]map[string]string
}

func CommitDeploy(state *DeploymentState, checker BlobChecker, deployID string, manifest map[string]string) error
func RollbackDeploy(state *DeploymentState, targetDeployID string) error
```

- `CommitDeploy`:
  - Verify every hash in `manifest` exists in `checker`. Return `ErrMissingBlobs` if any is missing.
  - Store `manifest` in `state.Manifests[deployID]`.
  - Set `state.ActiveDeployID = deployID`.
- `RollbackDeploy`:
  - If `targetDeployID` does not exist in `state.Manifests`, return `ErrDeployNotFound`.
  - Set `state.ActiveDeployID = targetDeployID`.
