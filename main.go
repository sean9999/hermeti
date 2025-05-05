package hermeti

func NewRealCli(exe InitRunner) CLI {
	env := RealEnv()
	return CLI{
		Env: env,
		Cmd: exe,
	}
}

func NewTestCli(exe InitRunner) CLI {
	env := TestEnv()
	return CLI{
		Env: env,
		Cmd: exe,
	}
}
