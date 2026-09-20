{
	// Successful streaming upload with matching hash
	store := NewMockBlobStore()
	handler := UploadBlobHandler(store)

	content := "hello world"
	hash := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

	req := httptest.NewRequest(http.MethodPut, "/api/blob/"+hash, strings.NewReader(content))
	rec := httptest.NewRecorder()

	handler(rec, req)

	Tests.EqualCheck("upload returns 201 Created", http.StatusCreated, rec.Code)
	Tests.EqualCheck("stored payload matches", content, string(store.Blobs[hash]))
}

{
	// Corrupt upload rejected and cleaned up
	store := NewMockBlobStore()
	handler := UploadBlobHandler(store)

	corruptedContent := "tampered content payload"
	claimedHash := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

	req := httptest.NewRequest(http.MethodPut, "/api/blob/"+claimedHash, strings.NewReader(corruptedContent))
	rec := httptest.NewRecorder()

	handler(rec, req)

	Tests.EqualCheck("mismatch returns 400 Bad Request", http.StatusBadRequest, rec.Code)
	Tests.BoolCheck("corrupted blob was not kept in store", store.Blobs[claimedHash] == nil)
}

{
	// Disallow GET method
	store := NewMockBlobStore()
	handler := UploadBlobHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/api/blob/somehash", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	Tests.EqualCheck("GET returns 405 Method Not Allowed", http.StatusMethodNotAllowed, rec.Code)
}
