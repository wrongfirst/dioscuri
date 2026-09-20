Before uploading files, the client sends a list of content hashes to the server to check which blobs already exist. The server inspects its blob store and returns only the subset of hashes that are missing, avoiding redundant network uploads.

## Handler Execution Flow

1. **Method Check**: Validate that `r.Method == http.MethodPost`. If not, respond with `http.StatusMethodNotAllowed` (`405`).
2. **Decode Request**: Parse JSON from `r.Body` into `CheckRequest`:
   ```json
   { "hashes": ["hash1", "hash2"] }
   ```
   If decoding fails, respond with `http.StatusBadRequest` (`400`).
3. **Inspect Store**: Iterate through `req.Hashes`. For each hash where `store.HasBlob(h)` returns `false`, append it to a `missing` slice (initialized non-nil: `make([]string, 0)`).
4. **Write Response**: Set `Content-Type: application/json`, write status `http.StatusOK` (`200`), and encode:
   ```json
   { "missing": ["hash2"] }
   ```

## Task

Implement `DeltaCheckHandler`:

```go
type BlobChecker interface {
    HasBlob(hash string) bool
}

type CheckRequest struct {
    Hashes []string `json:"hashes"`
}

type CheckResponse struct {
    Missing []string `json:"missing"`
}

func DeltaCheckHandler(store BlobChecker) http.HandlerFunc
```
