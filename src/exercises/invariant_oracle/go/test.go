{
	fs := NewMockFS()
	store := NewStore(fs, "/blobs")
	oracle := NewOracle(fs, store, "/manifests")

	// 1. Setup valid system state
	content := []byte("<h1>Dioscuri</h1>")
	h := sha256.Sum256(content)
	validHash := hex.EncodeToString(h[:])

	fs.Files["/blobs/"+validHash] = content
	fs.Files["/manifests/dep-001.json"] = []byte(`{"/index.html":"` + validHash + `"}`)
	fs.Files["/production.json"] = []byte(`{"deploy_id":"dep-001"}`)

	Tests.EqualCheck("valid state passes AssertAll", nil, oracle.AssertAll())

	// 2. Corrupt blob content (integrity violation)
	fs.Files["/blobs/"+validHash] = []byte("tampered content")
	Tests.EqualCheck("tampered blob triggers ErrIntegrityViolation", ErrIntegrityViolation, oracle.AssertBlobIntegrity())
	fs.Files["/blobs/"+validHash] = content // restore

	// 3. Dangling reference in a manifest
	fs.Files["/manifests/dep-002.json"] = []byte(`{"/about.html":"non_existent_blob_hash"}`)
	Tests.EqualCheck("dangling reference triggers ErrDanglingReference", ErrDanglingReference, oracle.AssertNoDanglingReferences())
	delete(fs.Files, "/manifests/dep-002.json") // restore

	// 4. Invalid production pointer
	fs.Files["/production.json"] = []byte(`{"deploy_id":"unknown_deploy_id"}`)
	Tests.EqualCheck("missing production manifest triggers ErrProductionInvalid", ErrProductionInvalid, oracle.AssertProductionSafety())
}
