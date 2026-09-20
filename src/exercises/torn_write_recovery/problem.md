Writing directly to `/blobs/<hash>` creates a catastrophic failure mode: if the process crashes mid-stream, a truncated or corrupted file remains in storage. Subsequent delta checks see the file exists and falsely assume it is intact.

To guarantee crash resilience and storage integrity, stream incoming payloads to a staging path first, verify the content digest on the fly, clean up staging on any failure, and atomically rename the verified file into place.

## Execution Flow

1. **Idempotency**: If `store.HasBlob(claimedHash)` is already true, return `nil` immediately.
2. **Staging File**: Create a temporary file at `path.Join(stagingDir, claimedHash+".tmp")`.
3. **On-the-Fly Hashing**: Use `io.TeeReader` to compute SHA-256 while streaming the payload from `r` directly into the staging file.
4. **Fault Cleanup**:
   - If `io.Copy` or `w.Close()` fails (e.g., mid-stream power cut `ErrPowerCut` or disk full `ErrDiskFull`), call `fs.Remove(stagingPath)` and return the error.
5. **Digest Verification**:
   - If the computed digest does not match `claimedHash`, call `fs.Remove(stagingPath)` and return `ErrCorruptPayload`.
6. **Atomic Promotion**:
   - Call `fs.Rename(stagingPath, store.BlobPath(claimedHash))`. If rename fails, remove the staging file and return the error.

## Task

Implement `SafeSaveBlob`:

```go
var ErrCorruptPayload = errors.New("corrupt payload")

func SafeSaveBlob(fs FS, store *Store, stagingDir string, claimedHash string, r io.Reader) error
```
