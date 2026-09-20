{
	// Build manifest from directory tree
	fs := NewMockFS()
	fs.Dirs = append(fs.Dirs, "/app/dist", "/app/dist/assets")
	fs.Files["/app/dist/index.html"] = []byte("hello world")
	fs.Files["/app/dist/assets/style.css"] = []byte("body { margin: 0; }")

	manifest, err := BuildManifest(fs, "/app/dist")
	Tests.EqualCheck("manifest build error is nil", nil, err)
	Tests.EqualCheck("manifest entry count", 2, len(manifest))

	helloHash := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	Tests.EqualCheck("index.html key present", helloHash, manifest["/index.html"])
	Tests.BoolCheck("style.css key present", manifest["/assets/style.css"] != "")
	Tests.BoolCheck("directory not present as key", manifest["/assets"] == "")
}
