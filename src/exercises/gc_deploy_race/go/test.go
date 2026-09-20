{
	// In-flight upload within quarantine is preserved
	store := NewMockGCStore()
	t0 := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	// Old unreferenced blob (uploaded 30m ago)
	store.AddBlob("hashOld", t0.Add(-30*time.Minute))

	// In-flight blob (uploaded 2m ago, not yet in committed manifest)
	store.AddBlob("hashInFlight", t0.Add(-2*time.Minute))

	// Committed blob (referenced by manifest)
	store.AddBlob("hashCommitted", t0.Add(-10*time.Minute))

	manifests := []map[string]string{
		{"/index.html": "hashCommitted"},
	}

	quarantine := 15 * time.Minute
	deleted, err := RunQuarantinedGC(store, manifests, t0, quarantine)

	Tests.EqualCheck("GC error is nil", nil, err)
	Tests.EqualCheck("only 1 blob deleted", 1, deleted)

	_, hasOld := store.blobs["hashOld"]
	Tests.EqualCheck("old blob was deleted", false, hasOld)

	_, hasInFlight := store.blobs["hashInFlight"]
	Tests.EqualCheck("in-flight blob was preserved by quarantine", true, hasInFlight)

	_, hasCommitted := store.blobs["hashCommitted"]
	Tests.EqualCheck("committed blob was preserved", true, hasCommitted)
}
