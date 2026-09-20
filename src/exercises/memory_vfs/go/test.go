{
	var fs FS = NewMemoryFS()

	// 1. Initial existence check
	Tests.EqualCheck("file does not exist initially", false, fs.Exists("/index.html"))

	// 2. Streamed write via Create
	w, err := fs.Create("/site/index.html")
	Tests.EqualCheck("create error is nil", nil, err)
	if w == nil {
		Tests.BoolCheck("Create returned valid io.WriteCloser", false)
		return
	}
	_, err = io.WriteString(w, "<h1>Pollux Memory VFS</h1>")
	Tests.EqualCheck("write error is nil", nil, err)
	err = w.Close()
	Tests.EqualCheck("close error is nil", nil, err)

	// 3. Existence after close
	Tests.EqualCheck("file exists after write", true, fs.Exists("/site/index.html"))

	// 4. Streamed read via Open
	r, err := fs.Open("/site/index.html")
	Tests.EqualCheck("open error is nil", nil, err)
	if r == nil {
		Tests.BoolCheck("Open returned valid io.ReadCloser", false)
		return
	}
	data, err := io.ReadAll(r)
	r.Close()
	Tests.EqualCheck("read error is nil", nil, err)
	Tests.EqualCheck("content matches", "<h1>Pollux Memory VFS</h1>", string(data))

	// 5. Rename file
	err = fs.Rename("/site/index.html", "/site/home.html")
	Tests.EqualCheck("rename error is nil", nil, err)
	Tests.EqualCheck("old path removed", false, fs.Exists("/site/index.html"))
	Tests.EqualCheck("new path exists", true, fs.Exists("/site/home.html"))

	// 6. ReadDir direct children
	w2, _ := fs.Create("/site/about.html")
	io.WriteString(w2, "about")
	w2.Close()

	children, err := fs.ReadDir("/site")
	Tests.EqualCheck("readDir error is nil", nil, err)
	Tests.EqualCheck("readDir child count is 2", 2, len(children))

	// 7. Remove file
	err = fs.Remove("/site/about.html")
	Tests.EqualCheck("remove error is nil", nil, err)
	Tests.EqualCheck("removed file does not exist", false, fs.Exists("/site/about.html"))

	// 8. WalkDir verification
	walkCount := 0
	err = fs.WalkDir("/site", func(p string, isDir bool) error {
		if !isDir {
			walkCount++
		}
		return nil
	})
	Tests.EqualCheck("walkDir error is nil", nil, err)
	Tests.EqualCheck("walkDir found exactly 1 remaining file", 1, walkCount)
}
