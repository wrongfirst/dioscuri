{
	// Deploy V1 and V2
	state := NewDeploymentState()
	checker := NewMockChecker("hashA", "hashB", "hashC")

	v1Manifest := map[string]string{"/index.html": "hashA", "/app.js": "hashB"}
	err1 := CommitDeploy(state, checker, "deploy-v1", v1Manifest)
	Tests.EqualCheck("commit deploy-v1 succeeds", nil, err1)
	Tests.EqualCheck("active deploy is deploy-v1", "deploy-v1", state.ActiveDeployID)

	v2Manifest := map[string]string{"/index.html": "hashC", "/app.js": "hashB"}
	err2 := CommitDeploy(state, checker, "deploy-v2", v2Manifest)
	Tests.EqualCheck("commit deploy-v2 succeeds", nil, err2)
	Tests.EqualCheck("active deploy is deploy-v2", "deploy-v2", state.ActiveDeployID)

	// Rollback to V1
	errRollback := RollbackDeploy(state, "deploy-v1")
	Tests.EqualCheck("rollback to deploy-v1 succeeds", nil, errRollback)
	Tests.EqualCheck("active deploy reverted to deploy-v1", "deploy-v1", state.ActiveDeployID)
}

{
	// Commit with missing blob fails
	state := NewDeploymentState()
	checker := NewMockChecker("hashA")

	incompleteManifest := map[string]string{"/index.html": "hashA", "/missing.png": "hashMissing"}
	err := CommitDeploy(state, checker, "bad-deploy", incompleteManifest)
	Tests.EqualCheck("missing blob returns ErrMissingBlobs", ErrMissingBlobs, err)
	Tests.EqualCheck("active deploy unchanged", "", state.ActiveDeployID)

	// Rollback to nonexistent deploy fails
	errNotFound := RollbackDeploy(state, "nonexistent")
	Tests.EqualCheck("rollback nonexistent returns ErrDeployNotFound", ErrDeployNotFound, errNotFound)
}
