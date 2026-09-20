{
	// 1. Clean pointer update via atomic swap
	fs := NewMockFS()
	fs.Files["/production.json"] = []byte(`{"deploy_id":"dep-001"}`)

	err := UpdateProductionPointer(fs, "dep-002")
	Tests.EqualCheck("pointer update error is nil", nil, err)

	data := string(fs.Files["/production.json"])
	Tests.BoolCheck("active pointer updated to dep-002", strings.Contains(data, "dep-002"))
	Tests.EqualCheck("temp file cleaned up", false, fs.Exists("/production.json.tmp"))

	// 2. Storage full fault leaves existing production pointer untouched
	faulty := NewFaultyFS(fs)
	faulty.SetDiskFull(true)

	errFull := UpdateProductionPointer(faulty, "dep-003")
	Tests.EqualCheck("disk full returns ErrDiskFull", ErrDiskFull, errFull)

	preservedData := string(fs.Files["/production.json"])
	Tests.BoolCheck("production pointer preserved intact at dep-002", strings.Contains(preservedData, "dep-002"))
}
