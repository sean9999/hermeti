package hermeti

// func New(env Env, fn Command) CLI {
// 	return CLI{
// 		Env: env,
// 		Cmd: fn,
// 	}
// }

func NewRealCli(fn Command) CLI {
	env := RealEnv()
	return CLI{
		Env: env,
		Cmd: fn,
	}
}
