In a manifest-driven web server, requests are served via manifest lookup rather than direct disk traversal:

1. **Production requests** (`/*`): Look up the file in the manifest of `ActiveDeployID`.
2. **Preview requests** (`/preview/<deploy_id>/*`): Look up the file in the manifest of the specified `deploy_id`.

A request for `"/"` resolves to `"/index.html"`.

## Task

Implement `RouteRequest`:

```go
type BlobOpener interface {
    OpenBlob(hash string) (io.ReadCloser, error)
}

type RouterState struct {
    ActiveDeployID string
    Manifests      map[string]map[string]string
}

type Response struct {
    StatusCode  int
    ContentType string
    Body        []byte
}

func RouteRequest(state *RouterState, store BlobOpener, urlPath string) Response
```

- Resolve root path `"/"` to `"/index.html"`.
- If `urlPath` begins with `"/preview/<deploy_id>/"`, route using that deployment's manifest. Otherwise, use `state.ActiveDeployID`.
- Return `StatusCode: 404` if the deployment or file path is not in the manifest.
- Return `StatusCode: 200` with Content-Type based on extension (`.html` $\rightarrow$ `"text/html"`, `.css` $\rightarrow$ `"text/css"`, `.js` $\rightarrow$ `"application/javascript"`, default `"application/octet-stream"`).
