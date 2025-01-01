package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/sean9999/hermeti"
	"github.com/sean9999/pear"
)

// we want this struct to be globally available
type globalState struct {
	Verbosity uint
	DieNow    bool
}

func (g *globalState) parse(args []string) []string {
	var dieNow bool
	var verbosity uint
	fset := flag.NewFlagSet("global", flag.ExitOnError)
	fset.BoolVar(&dieNow, "die", false, "should we die?")
	fset.UintVar(&verbosity, "verbosity", 0, "verbosity")
	fset.Parse(args)
	g.DieNow = dieNow
	g.Verbosity = verbosity
	return fset.Args()
}

func (g *globalState) Run(ctx context.Context, env hermeti.Env, args []string) []string {
	remainders := g.parse(args)

	if g.DieNow {
		panic("die now")
	}

	//	here's how we might choose to share global state
	ctx = context.WithValue(ctx, "globalState", g)

	fmt.Fprintf(env.OutStream, "verbosity now is:\t%d.\nArgs passed in are:\t%v\n", g.Verbosity, args)
	return hello(ctx, env, remainders)
}

func hello(ctx context.Context, env hermeti.Env, args []string) []string {
	gs, ok := ctx.Value("globalState").(*globalState)
	if !ok {
		panic("no global state")
	}
	fmt.Fprintf(env.OutStream, "the remainig args are:\t%s%v%s\n", Blue, args, Reset)
	fmt.Fprintf(env.OutStream, "verbosity still is:\t%s%v%s\n", Blue, gs.Verbosity, Reset)
	return args
}

func info(ctx context.Context, env hermeti.Env, args []string) []string {
	fmt.Fprintln(env.OutStream, ctx, env, args)
	return nil
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			pear.NicePanic(os.Stdout)
		}
	}()

	env := hermeti.RealEnv()

	ctx := context.Background()
	s := new(globalState)

	cli := &hermeti.CLI{Env: env, Cmd: s.Run}
	remainders := cli.Run(ctx, os.Args[1:])

	if len(remainders) > 0 {
		subcommand := remainders[0]

		switch subcommand {
		case "info":
			info(ctx, env, remainders[1:])
		default:
			panic("unknown subcommand")
		}

	}

}
