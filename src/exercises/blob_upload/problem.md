When uploading a blob, a client specifies the claimed hash. The server must verify that the payload matches this hash before persisting it.

If the incoming bytes do not match the claimed hash, the upload must be rejected to prevent corrupt data in storage.

## Task

Implement `IngestBlob`:

```go
var ErrHashMismatch = errors.New("hash mismatch")

type BlobSaver interface {
    SaveBlob(hash string, r io.Reader) error
}

func IngestBlob(saver BlobSaver, claimedHash string, r io.Reader) error
```

- Read bytes from `r` and compute their SHA-256 hex digest.
- If the digest does not match `claimedHash`, return `ErrHashMismatch` without saving.
- If the digest matches, save the payload via `saver.SaveBlob(claimedHash, ...)`.
