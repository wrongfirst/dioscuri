A deployment manifest maps canonical web paths to their content hashes:

```go
type Manifest map[string]string // "/assets/style.css" -> "sha256_hash"
```

Building a manifest requires walking a directory, filtering out directory nodes, computing the SHA-256 hash for each file, and storing each normalized path.

## Task

Implement `BuildManifest`:

```go
type WalkableFS interface {
    Walk(root string, fn func(path string, isDir bool) error) error
    Open(path string) (io.ReadCloser, error)
}

type Manifest map[string]string

func BuildManifest(fs WalkableFS, rootDir string) (Manifest, error)
```

- Walk `rootDir` using `fs.Walk`.
- For each non-directory file:
  - Strip the `rootDir` prefix and ensure the relative path starts with a leading `/`.
  - Convert any backslashes to forward slashes.
  - Compute the 64-character SHA-256 hex hash of the file.
  - Insert `manifest[path] = hash`.
- Directories must not be included in the returned manifest.
