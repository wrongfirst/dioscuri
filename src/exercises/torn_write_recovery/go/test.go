{
	// Successful staged save and rename
	fs := NewMockStagedFS()
	payload := "reliable-content-bytes"
	h := sha256.Sum256([]byte(payload))
	validHash := hex.EncodeToString(h[:])

	err := SafeSaveBlob(fs, validHash, strings.NewReader(payload))
	Tests.EqualCheck("valid blob save error is nil", nil, err)

	blobPath := fmt.Sprintf("/blobs/%s", validHash)
	stagingPath := fmt.Sprintf("/.staging/%s.tmp", validHash)

	Tests.EqualCheck("blob exists in destination", true, fs.Exists(blobPath))
	Tests.EqualCheck("staging path cleaned up", false, fs.Exists(stagingPath))
}

{
	// Corrupt payload leaves destination untouched
	fs := NewMockStagedFS()
	badClaimedHash := "0000000000000000000000000000000000000000000000000000000000000000"

	err := SafeSaveBlob(fs, badClaimedHash, strings.NewReader("incomplete or altered data"))
	Tests.EqualCheck("mismatch returns ErrCorruptPayload", ErrCorruptPayload, err)

	badBlobPath := fmt.Sprintf("/blobs/%s", badClaimedHash)
	Tests.EqualCheck("no corrupt file in blobs directory", false, fs.Exists(badBlobPath))
}
