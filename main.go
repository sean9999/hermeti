package hermeti

func NewRealCli(exe Runner) CLI {
	env := RealEnv()
	return CLI{
		Env: env,
		Cmd: exe,
	}
}

func NewTestCli[T any](exe Runner) CLI {
	env := TestEnv()
	return CLI{
		Env: env,
		Cmd: exe,
	}
}
