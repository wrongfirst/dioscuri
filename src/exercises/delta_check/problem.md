Before uploading files, a client sends its manifest hashes to the server to check which blobs are already present in storage.

This delta check prevents transferring files that already exist.

## Task

Implement `DeltaCheck`:

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

func DeltaCheck(store BlobChecker, req CheckRequest) CheckResponse
```

- For each hash in `req.Hashes`, check if `store.HasBlob(hash)` is false.
- Append missing hashes to `Missing`.
- Return an initialized (non-nil) slice in `CheckResponse`.
