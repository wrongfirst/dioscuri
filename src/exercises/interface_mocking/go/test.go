{
	// Read existing file
	mock := NewMockFS()
	mock.AddFile("/etc/app.json", []byte(`{"version":"1.0"}`))

	content, err := ReadFileString(mock, "/etc/app.json")
	Tests.EqualCheck("read existing file without error", nil, err)
	Tests.EqualCheck("content matches expected string", `{"version":"1.0"}`, content)
}

{
	// Missing file returns os.ErrNotExist
	mock := NewMockFS()
	_, err := ReadFileString(mock, "/not-found.txt")
	Tests.BoolCheck("missing file returns error", err != nil)
	Tests.EqualCheck("missing file returns os.ErrNotExist", os.ErrNotExist, err)
}

{
	// Multiple files in memory
	mock := NewMockFS()
	mock.AddFile("/index.html", []byte("<h1>Hello World</h1>"))
	mock.AddFile("/assets/app.js", []byte("console.log('ready');"))

	html, _ := ReadFileString(mock, "/index.html")
	js, _ := ReadFileString(mock, "/assets/app.js")

	Tests.EqualCheck("index.html matches", "<h1>Hello World</h1>", html)
	Tests.EqualCheck("app.js matches", "console.log('ready');", js)
}
