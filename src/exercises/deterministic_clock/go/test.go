{
	// 1. Deterministic sequence reproduction across identical seeds
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	h1 := NewSimHarness(42, start)
	h2 := NewSimHarness(42, start)

	for i := 0; i < 20; i++ {
		r1 := h1.RandomInt(1000)
		r2 := h2.RandomInt(1000)
		Tests.EqualCheck("random numbers match across identical seeds", r1, r2)
	}

	// 2. Virtual time progression via Advance
	timeHarness := NewSimHarness(100, start)
	Tests.BoolCheck("initial time matches start", start.Equal(timeHarness.Now()))

	timeHarness.Advance(5 * time.Minute)
	expected := time.Date(2026, 1, 1, 0, 5, 0, 0, time.UTC)
	Tests.BoolCheck("time advances by 5 minutes", expected.Equal(timeHarness.Now()))
}
