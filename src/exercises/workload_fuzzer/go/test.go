{
	// Site generation determinism
	rng1 := rand.New(rand.NewSource(1337))
	rng2 := rand.New(rand.NewSource(1337))

	site1 := GenerateSite(rng1, 5)
	site2 := GenerateSite(rng2, 5)

	Tests.EqualCheck("file count matches", 5, len(site1))
	for p, content := range site1 {
		Tests.EqualCheck("file content matches deterministically", content, site2[p])
	}
}

{
	// Site mutation preservation and addition
	rng := rand.New(rand.NewSource(99))
	base := GenerateSite(rng, 3)
	mutated := MutateSite(rng, base, 2)

	Tests.EqualCheck("mutated total count", 5, len(mutated))
	for p, content := range base {
		Tests.EqualCheck("base file preserved", content, mutated[p])
	}
	Tests.BoolCheck("new file present", mutated["/new_0.html"] != "")
	Tests.BoolCheck("second new file present", mutated["/new_1.html"] != "")
}
