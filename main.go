package hermeti

// func New(env Env, fn Command) CLI {
// 	return CLI{
// 		Env: env,
// 		Cmd: fn,
// 	}
// }

func NewRealCli[T any](exe Runner[T]) CLI[T] {
	env := RealEnv()
	return CLI[T]{
		Env: env,
		Cmd: exe,
	}
}
