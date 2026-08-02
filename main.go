package hermeti

func NewRealCli[T Runner](exe T) CLI[T] {
	env := RealEnv()
	return CLI[T]{
		Env: env,
		App: exe,
	}
}

func NewTestCli[T Runner](exe T, binaryName string) CLI[T] {
	env := TestEnv()
	env.Args = []string{binaryName}
	return CLI[T]{
		Env: env,
		App: exe,
	}
}
