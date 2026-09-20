{
	// Matching hash saves successfully
	saver := NewMockBlobSaver()
	helloHash := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	err := IngestBlob(saver, helloHash, strings.NewReader("hello world"))

	Tests.EqualCheck("valid upload succeeds", nil, err)
	Tests.EqualCheck("saved content matches", "hello world", string(saver.Saved[helloHash]))
}

{
	// Mismatched hash rejects with ErrHashMismatch
	saver := NewMockBlobSaver()
	claimedWrong := "0000000000000000000000000000000000000000000000000000000000000000"
	err := IngestBlob(saver, claimedWrong, strings.NewReader("hello world"))

	Tests.EqualCheck("mismatched upload returns ErrHashMismatch", ErrHashMismatch, err)
	Tests.EqualCheck("corrupt payload not saved", 0, len(saver.Saved))
}
