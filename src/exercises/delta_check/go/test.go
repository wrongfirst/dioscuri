{
	// Partial missing hashes
	store := NewMockBlobChecker("hashA", "hashC")
	req := CheckRequest{Hashes: []string{"hashA", "hashB", "hashC", "hashD"}}

	resp := DeltaCheck(store, req)
	Tests.EqualCheck("missing count is 2", 2, len(resp.Missing))
	Tests.EqualCheck("first missing is hashB", "hashB", resp.Missing[0])
	Tests.EqualCheck("second missing is hashD", "hashD", resp.Missing[1])
}

{
	// All present returns empty list
	store := NewMockBlobChecker("hash1", "hash2")
	req := CheckRequest{Hashes: []string{"hash1", "hash2"}}

	resp := DeltaCheck(store, req)
	Tests.EqualCheck("missing list is empty", 0, len(resp.Missing))
	Tests.BoolCheck("missing slice is non-nil", resp.Missing != nil)
}
