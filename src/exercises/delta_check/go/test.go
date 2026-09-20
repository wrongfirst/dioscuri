{
	// Partial missing hashes tested via httptest
	store := NewMockBlobChecker("hashA", "hashC")
	handler := DeltaCheckHandler(store)

	body := strings.NewReader(`{"hashes":["hashA","hashB","hashC","hashD"]}`)
	req := httptest.NewRequest(http.MethodPost, "/api/check", body)
	rec := httptest.NewRecorder()

	handler(rec, req)

	Tests.EqualCheck("status code is 200 OK", http.StatusOK, rec.Code)
	Tests.EqualCheck("content-type is application/json", "application/json", rec.Header().Get("Content-Type"))

	var resp CheckResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	Tests.EqualCheck("json decode error is nil", nil, err)
	Tests.EqualCheck("missing count is 2", 2, len(resp.Missing))
	Tests.EqualCheck("first missing is hashB", "hashB", resp.Missing[0])
	Tests.EqualCheck("second missing is hashD", "hashD", resp.Missing[1])
}

{
	// All present returns empty non-nil missing slice
	store := NewMockBlobChecker("hash1", "hash2")
	handler := DeltaCheckHandler(store)

	body := strings.NewReader(`{"hashes":["hash1","hash2"]}`)
	req := httptest.NewRequest(http.MethodPost, "/api/check", body)
	rec := httptest.NewRecorder()

	handler(rec, req)

	Tests.EqualCheck("status code is 200 OK", http.StatusOK, rec.Code)
	var resp CheckResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	Tests.EqualCheck("missing count is 0", 0, len(resp.Missing))
	Tests.BoolCheck("missing slice is non-nil", resp.Missing != nil)
}

{
	// Disallow GET method
	store := NewMockBlobChecker()
	handler := DeltaCheckHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/api/check", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	Tests.EqualCheck("GET method rejected with 405", http.StatusMethodNotAllowed, rec.Code)
}
