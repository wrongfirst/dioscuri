{
	// Production request with default index.html and content-type
	store := NewMockBlobOpener()
	store.Blobs["hashHTML"] = []byte("<h1>Welcome</h1>")
	store.Blobs["hashCSS"] = []byte("body { color: red; }")

	state := &RouterState{
		ActiveDeployID: "v1",
		Manifests: map[string]map[string]string{
			"v1": {
				"/index.html":       "hashHTML",
				"/assets/style.css": "hashCSS",
			},
		},
	}

	resRoot := RouteRequest(state, store, "/")
	Tests.EqualCheck("root path status 200", 200, resRoot.StatusCode)
	Tests.EqualCheck("root path content type text/html", "text/html", resRoot.ContentType)
	Tests.EqualCheck("root body matches", "<h1>Welcome</h1>", string(resRoot.Body))

	resCSS := RouteRequest(state, store, "/assets/style.css")
	Tests.EqualCheck("css status 200", 200, resCSS.StatusCode)
	Tests.EqualCheck("css content type text/css", "text/css", resCSS.ContentType)
}

{
	// Preview routing and 404 handling
	store := NewMockBlobOpener()
	store.Blobs["hashV2"] = []byte("<h1>Preview V2</h1>")

	state := &RouterState{
		ActiveDeployID: "v1",
		Manifests: map[string]map[string]string{
			"v1": {"/index.html": "hashV1"},
			"v2": {"/index.html": "hashV2"},
		},
	}

	resPreview := RouteRequest(state, store, "/preview/v2/index.html")
	Tests.EqualCheck("preview status 200", 200, resPreview.StatusCode)
	Tests.EqualCheck("preview body matches v2", "<h1>Preview V2</h1>", string(resPreview.Body))

	res404 := RouteRequest(state, store, "/nonexistent.txt")
	Tests.EqualCheck("missing path returns 404", 404, res404.StatusCode)
}
