package hermeti

import (
	"bytes"
	"io"
	"strings"

	"github.com/sean9999/pear"
)

// A Runner takes an [Env] and runs some code against it.
type Runner interface {
	Run(*Env)
}

// A CLI is a command. It runs an app against an environment.
// It does not need to be a CLI in the strict sense of the term.
// It can be a utility or daemon. It probably shouldn't be a TUI though
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

// Invoke a cli by using this string and Run ing
func (cli *CLI[T]) Invoke(str string) {
	args := strings.Split(str, " ")
	cli.Env.Args = args
	cli.Run()
}

// Run runs the Runners Run method, passing in Env.
// It's simply a convenience function.
func (cli CLI[T]) Run() {
	cli.App.Run(cli.Env)
}

var ErrOutputNotReadable = pear.Defer("output stream is not readable")

// OutStream returns an io.Reader representing the stuff you put in StdOut.
//
//	This will not work in a real CLI because os.StdOut is not readable
func (cli CLI[T]) OutStream() (*bytes.Buffer, error) {

	o, ok := cli.Env.OutStream.(io.Reader)
	if !ok {
		return nil, ErrOutputNotReadable
	}

	b, err := io.ReadAll(o)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(b)

	return buff, nil

}

// ErrStream returns an io.Reader representing the stuff you put in StdErr.
//
//	This will not work in a real CLI because os.StdOut is not readable
func (cli CLI[T]) ErrStream() (*bytes.Buffer, error) {

	o, ok := cli.Env.ErrStream.(io.Reader)
	if !ok {
		return nil, ErrOutputNotReadable
	}

	b, err := io.ReadAll(o)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(b)

	return buff, nil

}
