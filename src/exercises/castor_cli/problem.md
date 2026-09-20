In this assembly capstone, you bring together the storage engine, HTTP handlers, and client into a fully functional Content-Addressable deployment platform.

## Code Self-Review

The starter template includes reference implementations of all units you built throughout Chapter 1 (`HashFile`, `BuildManifest`, `Store`, `DeltaCheckHandler`, `UploadBlobHandler`, `CommitDeploy`, `NewWebServer`).

> [!NOTE]
> Before continuing, inspect the pre-populated code in the editor and compare it with the solutions you wrote in units 1.1 – 1.8. Check for differences in error handling, interface boundaries, and stream closures.

---

## Assembly Architecture & Flow

### 1. Server Routing (`NewServer`)
Mounts the HTTP handlers onto a unified `http.NewServeMux()`:
- `POST /api/check` $\to$ `DeltaCheckHandler(store)`
- `PUT/POST /api/blob/` $\to$ `UploadBlobHandler(store)`
- `POST /api/deploy` $\to$ Decodes `{"deploy_id":"...", "manifest":{...}}` and commits via `CommitDeploy(fs, store, deployID, manifest)`.
- Catch-all `/` $\to$ `NewWebServer(fs, store)`

### 2. Client Deployment Pipeline (`DeployPipeline`)
Coordinates the end-to-end synchronization workflow:
1. **Manifest Generation**: Call `BuildManifest(fs, localDir)` to obtain the file-to-hash mapping.
2. **Delta Query**: Collect all hashes and call `client.Check(hashes)` to determine which blobs are missing on the remote server.
3. **Selective Streaming Upload**: For each hash in `missing`, open the local file from `fs` and stream upload it via `client.UploadBlob(hash, r)`. Unchanged files are skipped.
4. **Atomic Commit**: Call `client.Deploy(deployID, manifest)` to activate the new deployment.
5. **Return Stats**: Return `&DeployStats` with `DeployID`, `TotalFiles`, and `UploadedFiles`.

---

## Task

Implement `NewServer` and `DeployPipeline`:

```go
type DeployStats struct {
    DeployID      string
    TotalFiles    int
    UploadedFiles int
}

func NewServer(fs FS, store *Store) http.Handler
func DeployPipeline(client *Client, fs FS, localDir string, deployID string) (*DeployStats, error)
```

---

## Running Locally

In your local `dioscuri` workspace, this same code runs natively:

```bash
# 1. Start the Castor server on port 8080
go run ./cmd/castor-server

# 2. In another terminal, deploy the test site
go run ./cmd/castor-cli deploy ./testdata/site-v1 --server http://localhost:8080

# 3. View your deployed site live
curl http://localhost:8080/
```
