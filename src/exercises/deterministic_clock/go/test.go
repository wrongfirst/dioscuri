{
	// Deterministic sequence reproduction
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	h1 := NewSimHarness(42, start)
	h2 := NewSimHarness(42, start)

	for i := 0; i < 20; i++ {
		r1 := h1.RandomInt(1000)
		r2 := h2.RandomInt(1000)
		Tests.EqualCheck("random numbers match across identical seeds", r1, r2)
	}
}

{
	// Virtual time advancement
	start := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	h := NewSimHarness(100, start)
	Tests.EqualCheck("initial time matches start", start, h.Now())

	h.Advance(5 * time.Minute)
	expected := time.Date(2026, 1, 1, 12, 5, 0, 0, time.UTC)
	Tests.EqualCheck("time advances by 5 minutes", expected, h.Now())
}
