Writing directly to `/blobs/<hash>` creates a vulnerability: if the process crashes mid-write, a truncated file remains in storage. Later delta checks see the file exists and falsely assume it is intact.

To guarantee crash resilience, write to a temporary staging file first, verify the content digest, and atomically rename it into place.

## Task

Implement `SafeSaveBlob`:

```go
var ErrCorruptPayload = errors.New("corrupt payload")

type StagedFS interface {
    WriteFile(path string, data []byte) error
    Rename(oldPath, newPath string) error
    Remove(path string) error
    Exists(path string) bool
    ReadFile(path string) ([]byte, error)
}

func SafeSaveBlob(fs StagedFS, claimedHash string, r io.Reader) error
```

1. Read all bytes from `r` and compute their SHA-256 hex digest.
2. If reading fails or the digest does not match `claimedHash`, return `ErrCorruptPayload` without leaving any file at `/blobs/<claimedHash>`.
3. Write the verified bytes to `/.staging/<claimedHash>.tmp`.
4. Atomically rename the staging path to `/blobs/<claimedHash>`.
