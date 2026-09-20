{
	// Read existing file
	mem := NewMemoryOpener()
	mem.AddFile("/etc/app.json", []byte(`{"version":"1.0"}`))

	content, err := ReadFileString(mem, "/etc/app.json")
	Tests.EqualCheck("read existing file without error", nil, err)
	Tests.EqualCheck("content matches expected string", `{"version":"1.0"}`, content)
}

{
	// Missing file returns os.ErrNotExist
	mem := NewMemoryOpener()
	_, err := ReadFileString(mem, "/not-found.txt")
	Tests.BoolCheck("missing file returns error", err != nil)
	Tests.EqualCheck("missing file returns os.ErrNotExist", os.ErrNotExist, err)
}

{
	// Multiple files in memory
	mem := NewMemoryOpener()
	mem.AddFile("/index.html", []byte("<h1>Hello World</h1>"))
	mem.AddFile("/assets/app.js", []byte("console.log('ready');"))

	html, _ := ReadFileString(mem, "/index.html")
	js, _ := ReadFileString(mem, "/assets/app.js")

	Tests.EqualCheck("index.html matches", "<h1>Hello World</h1>", html)
	Tests.EqualCheck("app.js matches", "console.log('ready');", js)
}
