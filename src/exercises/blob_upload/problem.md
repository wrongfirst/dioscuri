When uploading a blob, the client transmits raw bytes while claiming a specific content hash. The server must verify that the streamed bytes match the claimed hash before finalizing storage.

To prevent memory exhaustion on large file uploads, the server streams the incoming request body directly to disk while concurrently hashing bytes on the fly via `io.TeeReader`.

## Handler Execution Flow

1. **Method Check**: Accept `PUT` or `POST`. Reject other methods with `http.StatusMethodNotAllowed` (`405`).
2. **Hash Extraction**: Extract `claimedHash` from URL path prefix `/api/blob/` (e.g. `/api/blob/<hash>`). If missing or invalid, respond with `http.StatusBadRequest` (`400`).
3. **Concurrent Stream & Hash**:
   - Initialize `hasher := sha256.New()`.
   - Wrap `r.Body` with `tee := io.TeeReader(r.Body, hasher)`.
4. **Persist Payload**: Stream through `store.SaveBlob(claimedHash, tee)`.
5. **Integrity Validation**:
   - Compute `actualHash := hex.EncodeToString(hasher.Sum(nil))`.
   - If `actualHash != claimedHash`: remove corrupt data via `store.RemoveBlob(claimedHash)` and respond with `http.StatusBadRequest` (`400`).
   - If valid: respond with `http.StatusCreated` (`201`).

## Task

Implement `UploadBlobHandler`:

```go
type BlobStore interface {
    SaveBlob(hash string, r io.Reader) error
    RemoveBlob(hash string) error
}

func UploadBlobHandler(store BlobStore) http.HandlerFunc
```
