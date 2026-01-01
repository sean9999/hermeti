package hermeti

import (
	"bytes"
	"io"

	"strings"
	"github.com/sean9999/pear"
)

// A Runner takes an [Env] and runs some code against it.
// It cannot modify the Env.
type Runner interface {
	Run(Env)
}

// An Initializer initializes itself in preparation of running.
// It can modify its [Env]
type Initializer interface {
	Init(*Env) error
}

// PassthroughInit is an Initializer that does nothing
type PassthroughInit struct{}

func (p PassthroughInit) Init(_ *Env) error {
	return nil
}

type InitRunner interface {
	Runner
	Initializer
}

// A CLI is a command line interface. It runs an app against an environment
type CLI[T InitRunner] struct {
	Env         Env
	App         T
	initialized bool
}

func NewCLI[T InitRunner](env *Env, app T) *CLI[T] {
	return &CLI[T]{
		Env: *env,
		App: app,
	}
}


//	Invoke a cli by using this string and Run ing
func (c *CLI[T]) Invoke(str string) {
	args := strings.Split(str, " ")
	c.Env.Args = args
	c.Run()
}

// Run runs the Runners Run method, passing in Env.
// It's simply a convenience function.
func (cli CLI[T]) Run() {
	if !cli.initialized {
		err := cli.App.Init(&cli.Env)
		if err != nil {
			panic(err)
		}
		cli.initialized = true
	}
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
