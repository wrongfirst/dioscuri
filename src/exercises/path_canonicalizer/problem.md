Uploaded file paths can contain directory traversal attacks (e.g. `../../secret.env`), platform-specific backslashes (`\`), or redundant path elements. Before storing files in a manifest or serving them over HTTP, paths must be validated and canonicalized into deterministic POSIX format.

## Canonicalization Flow

1. **POSIX Conversion**: Replace all Windows backslashes (`\`) with forward slashes (`/`).
2. **Drive Letter Stripping**: If prefixed with a Windows drive letter (e.g. `C:` or `d:`), strip it so the remaining path is evaluated from the virtual root.
3. **Relative Path Resolution**:
   - Strip leading slashes to evaluate the relative path: `path.Clean(strings.TrimLeft(posix, "/"))`.
   - If the cleaned path equals `..` or begins with `../`, it attempts to escape above the root directory. Return `"", ErrPathEscapesRoot`.
4. **Absolute Formatting**:
   - If the path is empty, `"/"`, or `"."`, return `"/"`.
   - Otherwise, return `"/" + cleanRel` (ensuring a leading `/` and no trailing `/`).

## Task

Implement `CanonicalizePath`:

```go
func CanonicalizePath(rawPath string) (string, error)
```

- Return the canonicalized POSIX path starting with `/`.
- Reject any path that escapes the root with `ErrPathEscapesRoot`.
