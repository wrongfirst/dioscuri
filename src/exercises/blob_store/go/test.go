{
	// Save and retrieve blob
	fs := NewMemBlobFS()
	store := NewStore(fs, "/var/blobs")

	hash := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	content := "hello world"

	Tests.EqualCheck("hasBlob is false initially", false, store.HasBlob(hash))

	err := store.SaveBlob(hash, strings.NewReader(content))
	Tests.EqualCheck("saveBlob error is nil", nil, err)
	Tests.EqualCheck("hasBlob is true after save", true, store.HasBlob(hash))

	r, err := store.OpenBlob(hash)
	Tests.EqualCheck("openBlob error is nil", nil, err)
	data, _ := io.ReadAll(r)
	r.Close()
	Tests.EqualCheck("retrieved content matches", content, string(data))
}

{
	// Idempotency: saving existing blob is a no-op
	fs := NewMemBlobFS()
	store := NewStore(fs, "/var/blobs")

	hash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	err1 := store.SaveBlob(hash, strings.NewReader(""))
	Tests.EqualCheck("first save succeeds", nil, err1)

	err2 := store.SaveBlob(hash, strings.NewReader(""))
	Tests.EqualCheck("duplicate save succeeds idempotently", nil, err2)
}
