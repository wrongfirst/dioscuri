{
	// Empty stream
	src := strings.NewReader("")
	var dst bytes.Buffer
	n, err := CopyAndCount(&dst, src)
	Tests.EqualCheck("empty stream error is nil", nil, err)
	Tests.EqualCheck("empty stream bytes copied", int64(0), n)
	Tests.EqualCheck("empty stream dst content", "", dst.String())
}

{
	// Standard copy
	message := "streaming payload data"
	src := strings.NewReader(message)
	var dst bytes.Buffer
	n, err := CopyAndCount(&dst, src)
	Tests.EqualCheck("standard copy error is nil", nil, err)
	Tests.EqualCheck("standard copy bytes copied", int64(len(message)), n)
	Tests.EqualCheck("standard copy dst content matches", message, dst.String())
}

{
	// Large payload (100 KB)
	largeData := strings.Repeat("chunk-stream-bytes-", 5000)
	src := strings.NewReader(largeData)
	var dst bytes.Buffer
	n, err := CopyAndCount(&dst, src)
	Tests.EqualCheck("large copy error is nil", nil, err)
	Tests.EqualCheck("large copy byte count", int64(len(largeData)), n)
	Tests.EqualCheck("large copy content matches", largeData, dst.String())
}
