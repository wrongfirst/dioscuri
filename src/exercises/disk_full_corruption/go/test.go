{
	// Clean pointer update via atomic swap
	fs := NewMockAtomicFS()
	fs.Files["production.json"] = []byte(`{"active":"v1"}`)

	err := UpdateProductionPointer(fs, "v2")
	Tests.EqualCheck("pointer update error is nil", nil, err)

	data, _ := fs.ReadFile("production.json")
	Tests.EqualCheck("active pointer updated to v2", `{"active":"v2"}`, string(data))
	Tests.EqualCheck("temp file cleaned up", false, fs.Files["production.json.tmp"] != nil)
}

{
	// Disk full leaves original production pointer intact
	fs := NewMockAtomicFS()
	fs.Files["production.json"] = []byte(`{"active":"v1"}`)
	fs.DiskFull = true

	err := UpdateProductionPointer(fs, "v2")
	Tests.EqualCheck("disk full returns error", ErrDiskFull, err)

	data, _ := fs.ReadFile("production.json")
	Tests.EqualCheck("original pointer preserved intact", `{"active":"v1"}`, string(data))
}
