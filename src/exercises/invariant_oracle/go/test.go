{
	// Integrity check with valid and invalid blob
	emptyHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	validBlobs := map[string][]byte{
		emptyHash: []byte(""),
	}
	Tests.EqualCheck("valid integrity returns nil", nil, VerifyIntegrity(validBlobs))

	corruptBlobs := map[string][]byte{
		emptyHash: []byte("corrupted"),
	}
	Tests.EqualCheck("corrupt integrity returns ErrIntegrityViolation", ErrIntegrityViolation, VerifyIntegrity(corruptBlobs))
}

{
	// Dangling reference detection
	manifests := map[string]map[string]string{
		"v1": {"/index.html": "hashA", "/style.css": "hashB"},
	}
	completeBlobs := map[string][]byte{"hashA": []byte("A"), "hashB": []byte("B")}
	Tests.EqualCheck("complete store returns nil", nil, VerifyNoDangling(manifests, completeBlobs))

	missingBlobs := map[string][]byte{"hashA": []byte("A")}
	Tests.EqualCheck("missing blob returns ErrDanglingReference", ErrDanglingReference, VerifyNoDangling(manifests, missingBlobs))
}

{
	// Production validity
	state := &SystemState{
		ActiveDeployID: "v1",
		Manifests: map[string]map[string]string{
			"v1": {"/index.html": "hashA"},
		},
	}
	Tests.EqualCheck("valid production state returns nil", nil, VerifyProduction(state))

	state.ActiveDeployID = "vNonExistent"
	Tests.EqualCheck("nonexistent deploy returns ErrProductionInvalid", ErrProductionInvalid, VerifyProduction(state))
}
