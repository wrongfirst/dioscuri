{
	// Setup test environment
	fs := NewMockFS()
	store := NewMockStore()

	// Active production deployment: dep-prod
	fs.Files["/production.json"] = []byte(`{"deploy_id":"dep-prod"}`)
	fs.Files["/manifests/dep-prod.json"] = []byte(`{"/index.html":"hash-html","/assets/style.css":"hash-css"}`)

	// Preview deployment: dep-preview
	fs.Files["/manifests/dep-preview.json"] = []byte(`{"/index.html":"hash-preview-html"}`)

	// Populate blobs
	store.Blobs["hash-html"] = []byte("<h1>Production Site</h1>")
	store.Blobs["hash-css"] = []byte("body { color: red; }")
	store.Blobs["hash-preview-html"] = []byte("<h1>Preview Site</h1>")

	server := NewWebServer(fs, store)

	// 1. Root route "/" resolves to production index.html
	reqRoot := httptest.NewRequest(http.MethodGet, "/", nil)
	recRoot := httptest.NewRecorder()
	server.ServeHTTP(recRoot, reqRoot)
	Tests.EqualCheck("root status is 200", http.StatusOK, recRoot.Code)
	Tests.EqualCheck("root content-type is text/html", "text/html; charset=utf-8", recRoot.Header().Get("Content-Type"))
	Tests.EqualCheck("root content matches", "<h1>Production Site</h1>", recRoot.Body.String())

	// 2. Production asset route
	reqCSS := httptest.NewRequest(http.MethodGet, "/assets/style.css", nil)
	recCSS := httptest.NewRecorder()
	server.ServeHTTP(recCSS, reqCSS)
	Tests.EqualCheck("css status is 200", http.StatusOK, recCSS.Code)
	Tests.EqualCheck("css content-type is text/css", "text/css; charset=utf-8", recCSS.Header().Get("Content-Type"))
	Tests.EqualCheck("css content matches", "body { color: red; }", recCSS.Body.String())

	// 3. Preview route resolves to preview manifest
	reqPrev := httptest.NewRequest(http.MethodGet, "/preview/dep-preview/", nil)
	recPrev := httptest.NewRecorder()
	server.ServeHTTP(recPrev, reqPrev)
	Tests.EqualCheck("preview status is 200", http.StatusOK, recPrev.Code)
	Tests.EqualCheck("preview content matches", "<h1>Preview Site</h1>", recPrev.Body.String())

	// 4. Missing file returns 404
	reqMissing := httptest.NewRequest(http.MethodGet, "/missing.png", nil)
	recMissing := httptest.NewRecorder()
	server.ServeHTTP(recMissing, reqMissing)
	Tests.EqualCheck("missing file returns 404", http.StatusNotFound, recMissing.Code)
}
