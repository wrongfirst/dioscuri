{
	// 1. Site generation determinism across identical seeds
	rng1 := rand.New(rand.NewSource(1337))
	rng2 := rand.New(rand.NewSource(1337))

	site1 := GenerateSite(rng1, 5)
	site2 := GenerateSite(rng2, 5)

	if site1 == nil {
		Tests.BoolCheck("GenerateSite returned non-nil map", false)
		return
	}

	Tests.EqualCheck("file count matches", 5, len(site1))
	for p, content := range site1 {
		Tests.EqualCheck("file content matches deterministically", content, site2[p])
	}

	// 2. Site mutation: Add 2, Modify 1, Delete 1
	rng := rand.New(rand.NewSource(99))
	base := GenerateSite(rng, 4) // /page_0.html through /page_3.html
	mutated := MutateSite(rng, base, 2, 1, 1)

	if mutated == nil {
		Tests.BoolCheck("MutateSite returned non-nil map", false)
		return
	}

	// Expected: 4 - 1 deleted + 2 added = 5 total files
	Tests.EqualCheck("mutated total count is 5", 5, len(mutated))

	// Verify /page_0.html was deleted
	_, hasDeleted := mutated["/page_0.html"]
	Tests.EqualCheck("deleted page is absent", false, hasDeleted)

	// Verify /page_1.html was modified
	Tests.BoolCheck("page_1 content was modified", mutated["/page_1.html"] != base["/page_1.html"])

	// Verify /page_2.html was untouched
	Tests.EqualCheck("page_2 content was preserved", base["/page_2.html"], mutated["/page_2.html"])

	// Verify new files were added
	Tests.BoolCheck("new_0.html is present", mutated["/new_0.html"] != "")
	Tests.BoolCheck("new_1.html is present", mutated["/new_1.html"] != "")
}
