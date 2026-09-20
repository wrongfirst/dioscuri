{
	// Setup filesystem with 2 manifests and 5 blobs
	fs := NewMockFS()

	// Manifests
	fs.Files["/manifests/dep-001.json"] = []byte(`{"/index.html":"live1","/style.css":"live2"}`)
	fs.Files["/manifests/dep-002.json"] = []byte(`{"/index.html":"live2","/app.js":"live3"}`)

	// Blobs in storage
	fs.Files["/blobs/live1"] = []byte("blob 1")
	fs.Files["/blobs/live2"] = []byte("blob 2")
	fs.Files["/blobs/live3"] = []byte("blob 3")
	fs.Files["/blobs/orphanA"] = []byte("orphaned blob A")
	fs.Files["/blobs/orphanB"] = []byte("orphaned blob B")

	deleted, err := RunGC(fs, "/manifests", "/blobs")
	Tests.EqualCheck("GC error is nil", nil, err)
	Tests.EqualCheck("deleted count is 2", 2, deleted)

	// Verify live blobs retained
	Tests.BoolCheck("live1 preserved", fs.Files["/blobs/live1"] != nil)
	Tests.BoolCheck("live2 preserved", fs.Files["/blobs/live2"] != nil)
	Tests.BoolCheck("live3 preserved", fs.Files["/blobs/live3"] != nil)

	// Verify orphaned blobs deleted
	Tests.BoolCheck("orphanA deleted", fs.Files["/blobs/orphanA"] == nil)
	Tests.BoolCheck("orphanB deleted", fs.Files["/blobs/orphanB"] == nil)
}

{
	// No orphaned blobs to delete
	fs := NewMockFS()
	fs.Files["/manifests/dep-001.json"] = []byte(`{"/a":"h1"}`)
	fs.Files["/blobs/h1"] = []byte("content 1")

	deleted, err := RunGC(fs, "/manifests", "/blobs")
	Tests.EqualCheck("GC error is nil", nil, err)
	Tests.EqualCheck("zero deleted when all live", 0, deleted)
}
