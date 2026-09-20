Filesystem paths vary across operating systems (Windows `\` vs POSIX `/`) and may contain redundant elements (`./`, `../`) or directory traversal attempts.

Storage manifests require deterministic, normalized POSIX paths.

## Task

Implement `CanonicalizePath`:

```go
func CanonicalizePath(rawPath string) (string, error)
```

- Convert all backslashes (`\`) to forward slashes (`/`).
- Resolve all relative components using `path.Clean`.
- Ensure the path begins with a leading `/`.
- If the path attempts to traverse outside the root directory (e.g. `../secret`), return `ErrPathEscapesRoot`.
