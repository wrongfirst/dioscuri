{
	// 1. Initialize server environment and spin up in-memory HTTP server
	serverFS := NewMockFS()
	store := NewStore(serverFS, "/blobs")
	server := NewServer(serverFS, store)
	if server == nil {
		Tests.BoolCheck("NewServer returned a valid http.Handler", false)
		return
	}

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := NewClient(ts.URL)

	// 2. Prepare Client V1 Static Site
	clientFS := NewMockFS()
	clientFS.Dirs = append(clientFS.Dirs, "/dist")
	clientFS.Files["/dist/index.html"] = []byte("<h1>Hello V1</h1>")
	clientFS.Files["/dist/style.css"] = []byte("body { margin: 0; }")

	// 3. Deploy V1: Initial deployment (both files uploaded)
	statsV1, errV1 := DeployPipeline(client, clientFS, "/dist", "dep-001")
	Tests.EqualCheck("v1 deploy error is nil", nil, errV1)
	if statsV1 == nil {
		Tests.BoolCheck("DeployPipeline returned non-nil stats (v1)", false)
		return
	}
	Tests.EqualCheck("v1 total files is 2", 2, statsV1.TotalFiles)
	Tests.EqualCheck("v1 uploaded files is 2", 2, statsV1.UploadedFiles)

	// Verify server serves live V1 content
	respV1, errGetV1 := http.Get(ts.URL + "/")
	Tests.EqualCheck("get v1 root error is nil", nil, errGetV1)
	bodyV1, _ := io.ReadAll(respV1.Body)
	respV1.Body.Close()
	Tests.EqualCheck("get v1 status is 200", http.StatusOK, respV1.StatusCode)
	Tests.EqualCheck("live content is V1", "<h1>Hello V1</h1>", string(bodyV1))

	// 4. Update index.html for V2, leave style.css identical (Delta Sync Test)
	clientFS.Files["/dist/index.html"] = []byte("<h1>Hello V2</h1>")

	// 5. Deploy V2: Only 1 file should be uploaded across the wire
	statsV2, errV2 := DeployPipeline(client, clientFS, "/dist", "dep-002")
	Tests.EqualCheck("v2 deploy error is nil", nil, errV2)
	if statsV2 == nil {
		Tests.BoolCheck("DeployPipeline returned non-nil stats (v2)", false)
		return
	}
	Tests.EqualCheck("v2 total files is 2", 2, statsV2.TotalFiles)
	Tests.EqualCheck("v2 uploaded files is exactly 1 (delta upload)", 1, statsV2.UploadedFiles)

	// Verify server atomically serves updated V2 content
	respV2, errGetV2 := http.Get(ts.URL + "/")
	Tests.EqualCheck("get v2 root error is nil", nil, errGetV2)
	bodyV2, _ := io.ReadAll(respV2.Body)
	respV2.Body.Close()
	Tests.EqualCheck("get v2 status is 200", http.StatusOK, respV2.StatusCode)
	Tests.EqualCheck("live content atomically switched to V2", "<h1>Hello V2</h1>", string(bodyV2))
}
