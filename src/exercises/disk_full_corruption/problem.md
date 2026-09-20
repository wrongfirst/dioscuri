Writing critical state files directly in-place risks corruption: if the disk runs out of space mid-write, the file is truncated or corrupted, bringing down production.

To protect pointer updates against disk-full errors, write to a temporary file first and atomically swap it into place.

## Task

Implement `UpdateProductionPointer`:

```go
type AtomicFS interface {
    WriteFile(path string, data []byte) error
    Rename(oldPath, newPath string) error
    Remove(path string) error
    ReadFile(path string) ([]byte, error)
}

func UpdateProductionPointer(fs AtomicFS, deployID string) error
```

1. Format the pointer as JSON: `fmt.Sprintf("{\"active\":\"%s\"}", deployID)`.
2. Write to `production.json.tmp`.
3. If writing fails, return the error. Existing `production.json` must remain intact.
4. If writing succeeds, atomically rename `production.json.tmp` to `production.json`.
