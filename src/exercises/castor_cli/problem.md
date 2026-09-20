The deployment client orchestrates the deployment pipeline in 4 steps:

1. **Collect**: Extract all content hashes from the manifest.
2. **Delta Check**: Query the server for missing hashes.
3. **Upload**: Upload only the missing blobs.
4. **Deploy**: Commit the manifest and deploy ID.

## Task

Implement `RunDeploymentPipeline`:

```go
type ClientServerAPI interface {
    Check(hashes []string) (missing []string, err error)
    Upload(hash string, data []byte) error
    Deploy(deployID string, manifest map[string]string) error
}

type DeployStats struct {
    DeployID      string
    TotalFiles    int
    UploadedFiles int
}

func RunDeploymentPipeline(
    api ClientServerAPI,
    deployID string,
    localFiles map[string][]byte,
    manifest map[string]string,
) (DeployStats, error)
```

- Send all hashes from `manifest` to `api.Check`.
- For each hash in `missing`, locate its file bytes in `localFiles` and call `api.Upload`.
- Call `api.Deploy(deployID, manifest)` to activate the deployment.
- Return `DeployStats` with `DeployID`, `TotalFiles`, and `UploadedFiles`.
