{
	// Empty input matches SHA-256 standard
	emptyHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	h, err := HashReader(strings.NewReader(""))
	Tests.EqualCheck("empty stream error is nil", nil, err)
	Tests.EqualCheck("empty stream matches SHA-256 standard", emptyHash, h)
}

{
	// Known string digest
	helloHash := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	h, err := HashReader(strings.NewReader("hello world"))
	Tests.EqualCheck("hello world error is nil", nil, err)
	Tests.EqualCheck("hello world hash matches", helloHash, h)
}

{
	// Determinism
	content := "deterministic-content-payload-42"
	h1, _ := HashReader(strings.NewReader(content))
	h2, _ := HashReader(strings.NewReader(content))
	Tests.EqualCheck("identical content produces identical hash", h1, h2)

	// Sensitivity
	h3, _ := HashReader(strings.NewReader(content + "!"))
	Tests.BoolCheck("different content produces different hash", h1 != h3)
}
