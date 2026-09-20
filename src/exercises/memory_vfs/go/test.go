{
	// Write, read, exists
	fs := NewMemoryFS()
	Tests.EqualCheck("file does not exist initially", false, fs.Exists("/index.html"))

	err := fs.WriteFile("/index.html", []byte("<h1>Dioscuri</h1>"))
	Tests.EqualCheck("writeFile error is nil", nil, err)
	Tests.EqualCheck("file exists after write", true, fs.Exists("/index.html"))

	data, err := fs.ReadFile("/index.html")
	Tests.EqualCheck("readFile error is nil", nil, err)
	Tests.EqualCheck("readFile content matches", "<h1>Dioscuri</h1>", string(data))
}

{
	// Rename, remove, and list
	fs := NewMemoryFS()
	fs.WriteFile("/tmp/file.txt", []byte("temporary"))
	fs.WriteFile("/tmp/other.txt", []byte("other"))

	err := fs.Rename("/tmp/file.txt", "/data/permanent.txt")
	Tests.EqualCheck("rename error is nil", nil, err)
	Tests.EqualCheck("old path removed", false, fs.Exists("/tmp/file.txt"))
	Tests.EqualCheck("new path exists", true, fs.Exists("/data/permanent.txt"))

	err = fs.Remove("/tmp/other.txt")
	Tests.EqualCheck("remove error is nil", nil, err)
	Tests.EqualCheck("other removed", false, fs.Exists("/tmp/other.txt"))

	list := fs.ListDir("/data")
	Tests.EqualCheck("list count", 1, len(list))
	Tests.EqualCheck("listed item path", "/data/permanent.txt", list[0])
}
