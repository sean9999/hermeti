package hermeti

import (
	"bytes"
	"io"

	"github.com/sean9999/pear"
)

// a Runner takes an [Env] and runs some code against it
type Runner interface {
	Run(Env)
}

// an Initializer initializes itself in preparation of running
type Initializer interface {
	Init(Env) error
}

type PassthroughInit struct{}

func (p PassthroughInit) Init(_ Env) error {
	return nil
}

type InitRunner interface {
	Runner
	Initializer
}

// a CLI is a command line interface.
type CLI struct {
	Env         Env
	Cmd         InitRunner
	initialized bool
}

// Run() runs the Runners Run method, passing in Env.
// It's simly a convenience function.
func (cli CLI) Run() {
	if !cli.initialized {
		cli.Cmd.Init(cli.Env)
		cli.initialized = true
	}
	cli.Cmd.Run(cli.Env)
}

// Obj exposes our application object
// func (cli CLI) Obj() T {
// 	return cli.Cmd.(T)
// }

var ErrOutputNotReadable = pear.Defer("output stream is not readable")

// OutStream returns an io.Reader representing the stuff you put in StdOut.
//
//	This will not work in a real CLI because os.StdOut is not readable
func (cli CLI) OutStream() (*bytes.Buffer, error) {

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
