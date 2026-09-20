{
	// Empty stream
	src := strings.NewReader("")
	var dst bytes.Buffer
	n, err := StreamCopy(&dst, src, 16)
	Tests.EqualCheck("empty stream error is nil", nil, err)
	Tests.EqualCheck("empty stream bytes copied", int64(0), n)
	Tests.EqualCheck("empty stream dst content", "", dst.String())
}

{
	// Chunked copy with small buffer (4 bytes)
	message := "streaming payload data with chunking"
	src := strings.NewReader(message)
	var dst bytes.Buffer
	n, err := StreamCopy(&dst, src, 4)
	Tests.EqualCheck("chunked copy error is nil", nil, err)
	Tests.EqualCheck("chunked copy bytes copied", int64(len(message)), n)
	Tests.EqualCheck("chunked copy dst content matches", message, dst.String())
}

{
	// Large payload (80 KB) with 1024-byte buffer
	largeData := strings.Repeat("chunk-stream-bytes-", 4000)
	src := strings.NewReader(largeData)
	var dst bytes.Buffer
	n, err := StreamCopy(&dst, src, 1024)
	Tests.EqualCheck("large copy error is nil", nil, err)
	Tests.EqualCheck("large copy byte count", int64(len(largeData)), n)
	Tests.EqualCheck("large copy content matches", largeData, dst.String())
}

{
	// Invalid buffer size handling
	var dst bytes.Buffer
	src := strings.NewReader("some data")
	_, errZero := StreamCopy(&dst, src, 0)
	Tests.EqualCheck("zero buffer size returns error", ErrInvalidBufferSize, errZero)

	_, errNeg := StreamCopy(&dst, src, -8)
	Tests.EqualCheck("negative buffer size returns error", ErrInvalidBufferSize, errNeg)
}
