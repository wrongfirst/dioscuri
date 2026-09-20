In a manifest-driven web server, static files are not read directly from the physical filesystem. Requests are routed through manifest lookups that resolve URLs to content hashes, which are then streamed from the blob store.

Production requests (`/*`) resolve using `/production.json`. Preview requests (`/preview/<deploy_id>/*`) resolve using `/manifests/<deploy_id>.json`.

## Router Execution Flow

1. **Route Resolution**:
   - If `r.URL.Path` starts with `/preview/`:
     Extract the `<deploy_id>` and the subpath (e.g. `/preview/dep-001/css/style.css` $\to$ deploy `dep-001`, path `/css/style.css`).
   - Otherwise:
     Read `/production.json` from `fs` to find the active deploy ID. If missing, return `404 Not Found`.
   - If the resolved path is `"/"` or ends with `"/"`, map to `"index.html"` (e.g. `/` $\to$ `/index.html`).
2. **Manifest Lookup**:
   - Open and decode `/manifests/<deploy_id>.json` from `fs`. If the manifest does not exist, return `404 Not Found`.
   - Check if the requested file path exists in the manifest. If not, return `404 Not Found`.
3. **Blob Delivery**:
   - Retrieve the readable stream via `store.OpenBlob(hash)`. Ensure closure (`defer r.Close()`).
   - Set `Content-Type` by extension (`.html` $\to$ `text/html`, `.css` $\to$ `text/css`, `.js` $\to$ `application/javascript`, default `application/octet-stream`).
   - Stream content to the client via `io.Copy(w, blobReader)`.

## Task

Implement `NewWebServer`:

```go
type FS interface {
    Exists(path string) bool
    Open(path string) (io.ReadCloser, error)
}

type BlobOpener interface {
    OpenBlob(hash string) (io.ReadCloser, error)
}

func NewWebServer(fs FS, store BlobOpener) http.Handler
```
