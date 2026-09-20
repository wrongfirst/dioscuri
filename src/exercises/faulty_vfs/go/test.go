{
	// 1. Torn write simulation via CrashAfter
	mem := NewMockFS()
	faulty := NewFaultyFS(mem)
	faulty.CrashAfter(5)

	w, err := faulty.Create("/var/blobs/hash1")
	Tests.EqualCheck("create error is nil", nil, err)
	if w == nil {
		Tests.BoolCheck("Create returned valid io.WriteCloser", false)
		return
	}

	payload := []byte("hello world") // 11 bytes
	n, writeErr := w.Write(payload)
	Tests.EqualCheck("write returns ErrPowerCut", ErrPowerCut, writeErr)
	Tests.EqualCheck("wrote exactly 5 bytes before power cut", 5, n)
	w.Close()
	Tests.EqualCheck("disk contains partial 5 bytes", "hello", string(mem.Files["/var/blobs/hash1"]))

	// 2. Disk full simulation via SetDiskFull
	faulty.SetDiskFull(true)
	_, fullErr := faulty.Create("/test.txt")
	Tests.EqualCheck("disk full returns ErrDiskFull", ErrDiskFull, fullErr)

	// 3. Byte corruption simulation
	mem.Files["/file.bin"] = []byte{0x00, 0x00}
	faulty.SetDiskFull(false)
	faulty.CrashAfter(-1)
	err = faulty.CorruptByte("/file.bin", 1)
	Tests.EqualCheck("corruption error is nil", nil, err)
	Tests.EqualCheck("byte inverted", byte(0xFF), mem.Files["/file.bin"][1])
}
