{
	// Torn write after byte budget exceeded
	mem := NewMemoryFS()
	faulty := NewFaultyFS(mem)
	faulty.CrashAfter(5)

	payload := []byte("hello world") // 11 bytes
	err := faulty.WriteFile("/var/blobs/hash1", payload)

	Tests.EqualCheck("write returns ErrPowerCut", ErrPowerCut, err)
	Tests.EqualCheck("disk contains partial 5 bytes", "hello", string(mem.Files["/var/blobs/hash1"]))
}

{
	// Disk full simulation
	mem := NewMemoryFS()
	faulty := NewFaultyFS(mem)
	faulty.SetDiskFull(true)

	err := faulty.WriteFile("/test.txt", []byte("data"))
	Tests.EqualCheck("disk full returns ErrDiskFull", ErrDiskFull, err)
	Tests.EqualCheck("no file written to memory", 0, len(mem.Files))
}

{
	// Byte corruption
	mem := NewMemoryFS()
	mem.Files["/file.bin"] = []byte{0x00, 0x00}

	faulty := NewFaultyFS(mem)
	err := faulty.CorruptByte("/file.bin", 1)
	Tests.EqualCheck("corruption error is nil", nil, err)
	Tests.EqualCheck("byte inverted", byte(0xFF), mem.Files["/file.bin"][1])
}
