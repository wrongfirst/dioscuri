{
	// Orphaned blob deleted while shared blobs are preserved
	fs := NewMockGCFS("hashA", "hashB", "hashC", "hashD")

	manifest1 := map[string]string{"/index.html": "hashA", "/app.js": "hashB"}
	manifest2 := map[string]string{"/index.html": "hashA", "/style.css": "hashC"}

	// hashD is not referenced by any manifest
	deleted, err := RunGC(fs, []map[string]string{manifest1, manifest2})
	Tests.EqualCheck("GC error is nil", nil, err)
	Tests.EqualCheck("deleted count is 1", 1, deleted)

	Tests.EqualCheck("hashA still exists", true, fs.Blobs["hashA"])
	Tests.EqualCheck("hashB still exists", true, fs.Blobs["hashB"])
	Tests.EqualCheck("hashC still exists", true, fs.Blobs["hashC"])
	Tests.EqualCheck("hashD was removed", false, fs.Blobs["hashD"])
}
