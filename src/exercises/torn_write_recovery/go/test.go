{
	fs := NewMockFS()
	store := NewStore(fs, "/blobs")

	// 1. Successful staged save and atomic promotion
	payload := "reliable-content-bytes"
	h := sha256.Sum256([]byte(payload))
	validHash := hex.EncodeToString(h[:])

	err := SafeSaveBlob(fs, store, "/.staging", validHash, strings.NewReader(payload))
	Tests.EqualCheck("valid blob save error is nil", nil, err)
	Tests.EqualCheck("blob exists in destination", true, store.HasBlob(validHash))
	Tests.EqualCheck("staging path cleaned up", false, fs.Exists("/.staging/"+validHash+".tmp"))

	// 2. Corrupt payload leaves destination untouched
	badHash := "0000000000000000000000000000000000000000000000000000000000000000"
	err = SafeSaveBlob(fs, store, "/.staging", badHash, strings.NewReader("corrupt content"))
	Tests.EqualCheck("mismatch returns ErrCorruptPayload", ErrCorruptPayload, err)
	Tests.EqualCheck("no corrupt file in blobs directory", false, store.HasBlob(badHash))
	Tests.EqualCheck("staging cleaned up after corruption", false, fs.Exists("/.staging/"+badHash+".tmp"))

	// 3. Torn write simulation mid-stream
	faulty := NewFaultyFS(fs)
	faulty.CrashAfter(8)

	longPayload := "long stream exceeding byte budget"
	hLong := sha256.Sum256([]byte(longPayload))
	longHash := hex.EncodeToString(hLong[:])

	errCrash := SafeSaveBlob(faulty, store, "/.staging", longHash, strings.NewReader(longPayload))
	Tests.EqualCheck("crash returns ErrPowerCut", ErrPowerCut, errCrash)
	Tests.EqualCheck("incomplete blob not committed to destination", false, store.HasBlob(longHash))
	Tests.EqualCheck("staging cleaned up after power cut", false, fs.Exists("/.staging/"+longHash+".tmp"))
}
