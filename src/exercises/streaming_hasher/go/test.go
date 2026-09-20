{
	// Empty file matches SHA-256 standard
	fs := NewMockFS()
	fs.Files["empty.txt"] = []byte("")

	emptyHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	h, err := HashFile(fs, "empty.txt")
	Tests.EqualCheck("empty file error is nil", nil, err)
	Tests.EqualCheck("empty file matches SHA-256 standard", emptyHash, h)
}

{
	// Standard file content
	fs := NewMockFS()
	fs.Files["hello.txt"] = []byte("hello world")

	helloHash := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	h, err := HashFile(fs, "hello.txt")
	Tests.EqualCheck("hello.txt error is nil", nil, err)
	Tests.EqualCheck("hello.txt hash matches", helloHash, h)
}

{
	// Missing file returns error
	fs := NewMockFS()
	_, err := HashFile(fs, "missing.txt")
	Tests.BoolCheck("missing file returns error", err != nil)
}
