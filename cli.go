package hermeti

import "context"

// CLI command line interface.
// It is a [Command] that takes in arguments and operates on an [Env]
type CLI struct {
	Env Env
	Cmd Command
}

func (cli CLI) Run(ctx context.Context, args []string) ([]string, error) {
	return cli.Cmd(ctx, cli.Env, args)
}
