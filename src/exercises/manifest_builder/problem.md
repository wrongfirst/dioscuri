A deployment manifest maps canonical web paths to their cryptographic content hashes:

```go
type Manifest map[string]string // e.g. "/assets/style.css" -> "sha256_hash"
```

To build a manifest, the deployment client walks a local directory tree, ignores directories, computes the SHA-256 hash for each file, and records the relative canonical path.

## Execution Flow

1. **Path Preparation**: Normalize `rootDir` by converting backslashes to forward slashes and stripping any trailing slash.
2. **Directory Walk**: Traverse the directory using `fs.WalkDir(rootDir, fn)`.
3. **Node Processing**:
   - **Directories**: If `isDir` is true, skip by returning `nil`.
   - **Relative Path**: Strip the `rootDir` prefix from the file path and ensure it begins with a single leading `/`.
   - **Hash File**: Stream and hash the file using `HashFile(fs, p)`.
   - **Record**: Set `manifest[relPath] = hash`.
4. **Return**: Return the populated `Manifest` and any walk error encountered.

## Task

Implement `BuildManifest`:

```go
type FS interface {
    Open(name string) (io.ReadCloser, error)
    WalkDir(root string, fn func(path string, isDir bool) error) error
}

type Manifest map[string]string

func BuildManifest(fs FS, rootDir string) (Manifest, error)
```
