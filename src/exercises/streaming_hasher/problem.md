Content-Addressable Storage relies on computing content identities directly from files without loading entire files into memory. Reading through a filesystem interface abstraction allows the same hashing logic to run against both OS disks and virtual test filesystems.

## Execution Flow

1. **Open File**: Call `fs.Open(path)` to obtain an `io.ReadCloser`.
2. **Handle Open Errors**: If opening fails (e.g. file does not exist), return `"", err`.
3. **Schedule Closure**: Ensure the file descriptor is released upon function exit (`defer f.Close()`).
4. **Stream Hash**: Initialize a SHA-256 hasher via `sha256.New()` and stream bytes using `io.Copy(hasher, f)`.
5. **Format Hex Digest**: Compute `hasher.Sum(nil)` and return the 64-character lowercase hexadecimal string using `hex.EncodeToString`.

## Task

Implement `HashFile`:

```go
type FS interface {
    Open(path string) (io.ReadCloser, error)
}

func HashFile(fs FS, path string) (string, error)
```

- Open `path`, defer closure, and stream its content through SHA-256.
- Return the hexadecimal digest and any error encountered.
