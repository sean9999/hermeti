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

func (g *globalState) Run(env hermeti.Env) {

	//	we have no use for the first argument in [os.Args]
	args := env.Args[1:]
	args = g.parseArgs(args)
	if g.DieNow {
		panic("die now")
	}

	//	here's how we might choose to share state with the next [SubCommand]
	ctx := context.WithValue(context.Background(), "globalState", g)

	fmt.Fprintf(env.OutStream, "verbosity now is:\t%d.\nArgs passed in are:\t%v\n", g.Verbosity, args)

	if len(args) > 0 {
		if args[0] == "hello" {
			args, err := hello(ctx, env, args)
			fmt.Fprintln(env.OutStream, args, err)
		}
	}

}

func (exe *globalState) State() *globalState {
	return exe
}

func main() {

	//	run the cli. If panic, do so nicely
	defer func() {
		if r := recover(); r != nil {
			pear.NicePanic(os.Stdout)
		}
	}()

	env := hermeti.RealEnv()
	ctx := context.Background()
	exe := new(globalState)
	cli := &hermeti.CLI[*globalState]{Env: env, Cmd: exe}
	cli.Run(ctx)

}

func (g *globalState) parseArgs(args []string) []string {
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

func hello(ctx context.Context, env hermeti.Env, args []string) ([]string, error) {

	gs, ok := ctx.Value("globalState").(*globalState)
	if !ok {
		return args, pear.New("no global state")
	}

	fmt.Fprintf(env.OutStream, "the remainig args are:\t%s%v%s\n", Blue, args, Reset)
	fmt.Fprintf(env.OutStream, "verbosity still is:\t%s%v%s\n", Blue, gs.Verbosity, Reset)
	return args, nil
}

func info(ctx context.Context, env hermeti.Env, args []string) []string {
	fmt.Fprintln(env.OutStream, ctx, env, args)
	return nil
}
