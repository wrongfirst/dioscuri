{
	// Successful serialized deployment
	lock := NewMockDeployLock()
	coord := NewSafeDeployCoordinator(lock)

	err := coord.Deploy("v1", func() error {
		return nil
	})
	Tests.EqualCheck("first deploy succeeds", nil, err)
	Tests.EqualCheck("active deploy set to v1", "v1", coord.GetActive())
}

{
	// Concurrent deploy rejected while lock is held
	lock := NewMockDeployLock()
	coord := NewSafeDeployCoordinator(lock)

	// Simulate lock held by another process
	lock.TryLock()

	err := coord.Deploy("v2", func() error {
		return nil
	})
	Tests.EqualCheck("second deploy returns ErrDeployInProgress", ErrDeployInProgress, err)
	Tests.EqualCheck("active deploy remains empty", "", coord.GetActive())

	lock.Unlock()

	// Once released, deploy succeeds
	errRetry := coord.Deploy("v2", func() error {
		return nil
	})
	Tests.EqualCheck("retry deploy succeeds", nil, errRetry)
	Tests.EqualCheck("active deploy set to v2", "v2", coord.GetActive())
}
