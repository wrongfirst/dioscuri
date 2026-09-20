{
	// 100-step simulation runs cleanly without invariant failure
	rng := rand.New(rand.NewSource(42))
	res, err := RunSimulation(rng, 100)

	if res.StepsExecuted == 0 && err == nil {
		Tests.BoolCheck("RunSimulation implemented and executed steps", false)
		return
	}

	Tests.EqualCheck("simulation error is nil", nil, err)
	Tests.EqualCheck("executed 100 steps", 100, res.StepsExecuted)
	Tests.BoolCheck("multiple deploys executed", res.DeploysCount > 2)
	Tests.BoolCheck("multiple rollbacks executed", res.RollbackCount > 2)
	Tests.BoolCheck("multiple reads executed", res.ReadsCount > 2)
}
