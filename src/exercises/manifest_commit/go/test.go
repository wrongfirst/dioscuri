{
	// Commit succeeds when all blobs exist
	fs := NewMockFS()
	checker := NewMockBlobChecker("hash1", "hash2")
	manifest := Manifest{
		"/index.html": "hash1",
		"/style.css":  "hash2",
	}

	err := CommitDeploy(fs, checker, "dep-001", manifest)
	Tests.EqualCheck("commit error is nil", nil, err)

	Tests.BoolCheck("manifest file created", fs.Exists("/manifests/dep-001.json"))
	Tests.BoolCheck("production.json exists", fs.Exists("/production.json"))
	Tests.BoolCheck("tmp file was cleaned up via rename", !fs.Exists("/production.json.tmp"))

	r, _ := fs.Open("/production.json")
	var active ActivePointer
	json.NewDecoder(r).Decode(&active)
	r.Close()
	Tests.EqualCheck("active deploy ID is dep-001", "dep-001", active.DeployID)
}

{
	// Commit fails if any blob is missing
	fs := NewMockFS()
	checker := NewMockBlobChecker("hash1") // hash2 is missing
	manifest := Manifest{
		"/index.html": "hash1",
		"/style.css":  "hash2",
	}

	err := CommitDeploy(fs, checker, "dep-002", manifest)
	Tests.EqualCheck("commit returns ErrMissingBlobs", ErrMissingBlobs, err)
	Tests.BoolCheck("manifest was not created", !fs.Exists("/manifests/dep-002.json"))
}

{
	// Rollback succeeds for existing deployment
	fs := NewMockFS()
	checker := NewMockBlobChecker("h1", "h2")
	CommitDeploy(fs, checker, "dep-001", Manifest{"/a": "h1"})
	CommitDeploy(fs, checker, "dep-002", Manifest{"/b": "h2"})

	err := RollbackDeploy(fs, "dep-001")
	Tests.EqualCheck("rollback error is nil", nil, err)

	r, _ := fs.Open("/production.json")
	var active ActivePointer
	json.NewDecoder(r).Decode(&active)
	r.Close()
	Tests.EqualCheck("active pointer restored to dep-001", "dep-001", active.DeployID)

	// Rollback fails for non-existent deployment
	errMissing := RollbackDeploy(fs, "dep-999")
	Tests.EqualCheck("rollback returns ErrDeployNotFound", ErrDeployNotFound, errMissing)
}
