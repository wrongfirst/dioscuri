Decoupling filesystem operations behind an interface allows reading from both OS storage and virtual memory filesystems.

Computing the content identity of a file requires streaming its contents through SHA-256 and closing the file handle when finished.

## Task

Implement `HashFile`:

```go
type FS interface {
    Open(path string) (io.ReadCloser, error)
}

func HashFile(fs FS, path string) (string, error)
```

- Open `path` using `fs.Open(path)`.
- Ensure the opened file is closed when finished (`defer`).
- Stream the file contents through a SHA-256 hasher without buffering the entire file into memory.
- Return the 64-character lowercase hex digest string.
- If opening or reading fails, return an empty string and the error.
