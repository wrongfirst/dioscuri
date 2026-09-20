{
	// Pipeline with 2 total files and 1 missing file
	api := NewMockServerAPI("hash1") // hash1 already on server

	localFiles := map[string][]byte{
		"/index.html":       []byte("<h1>Hello</h1>"),
		"/assets/style.css": []byte("body { margin: 0; }"),
	}

	manifest := map[string]string{
		"/index.html":       "hash1",
		"/assets/style.css": "hash2", // missing, must be uploaded
	}

	stats, err := RunDeploymentPipeline(api, "release-101", localFiles, manifest)
	Tests.EqualCheck("pipeline succeeds without error", nil, err)
	Tests.EqualCheck("deployed ID matches", "release-101", stats.DeployID)
	Tests.EqualCheck("total files count", 2, stats.TotalFiles)
	Tests.EqualCheck("uploaded files count", 1, stats.UploadedFiles)

	Tests.EqualCheck("server activated release-101", "release-101", api.DeployedID)
	Tests.EqualCheck("server received uploaded hash2", "body { margin: 0; }", string(api.Uploaded["hash2"]))
}
