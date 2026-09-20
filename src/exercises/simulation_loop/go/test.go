{
	// 500-step simulation runs cleanly without invariant failure
	rng := rand.New(rand.NewSource(42))
	res, err := RunSimulation(rng, 500)

	Tests.EqualCheck("simulation error is nil", nil, err)
	Tests.EqualCheck("executed 500 steps", 500, res.StepsExecuted)
	Tests.BoolCheck("multiple deploys executed", res.DeploysCount > 10)
	Tests.BoolCheck("multiple rollbacks executed", res.RollbackCount > 10)
	Tests.BoolCheck("multiple reads executed", res.ReadsCount > 10)
}
