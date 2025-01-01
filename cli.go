package hermeti

import "context"

// CLI command line interface.
// It is a [Command] that takes in arguments and operates on an [Env]
type CLI struct {
	Env Env
	Cmd Command
}

func (fli CLI) Run(ctx context.Context, args []string) []string {
	return fli.Cmd(ctx, fli.Env, args)
}
