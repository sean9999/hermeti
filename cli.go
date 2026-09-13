package hermeti

// A Runner takes an [Env] and runs some code against it.
type Runner interface {
	Run(*Env)
}

// A CLI is a command-line interface.
type CLI[T Runner] struct {
	Env *Env
	App T
}

func NewCLI[T Runner](env *Env, app T) *CLI[T] {
	return &CLI[T]{
		Env: env,
		App: app,
	}
}

// Run runs the Runner's Run method, passing in Env.
func (cli CLI[T]) Run() {
	cli.App.Run(cli.Env)
}

// NewRealCli produces a CLI suitable for main()
func NewRealCli[T Runner](exe T) CLI[T] {
	env := RealEnv(".")
	return CLI[T]{
		Env: &env,
		App: exe,
	}
}

// NewTestCli produces a CLI suitable for tests
func NewTestCli[T Runner](exe T, binaryName string) CLI[T] {
	env := TestEnv()
	env.Args = []string{binaryName}
	return CLI[T]{
		Env: &env,
		App: exe,
	}
}
