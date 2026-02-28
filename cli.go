package hermeti

import (
	"bytes"
	"io"

	"github.com/sean9999/pear"
)

// A Runner takes an [Env] and runs some code against it.
// It cannot modify the Env.
type Runner interface {
	Run(*Env)
}

// A CLI is a command line interface. It runs an app against an environment
type CLI struct {
	Env         *Env
	Cmd         Runner
	initialized bool
}

// Run runs the Runners Run method, passing in Env.
// It's simply a convenience function.
func (cli CLI) Run() {
	cli.Cmd.Run(cli.Env)
}

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

// ErrStream returns an io.Reader representing the stuff you put in StdErr.
//
//	This will not work in a real CLI because os.StdOut is not readable
func (cli CLI) ErrStream() (*bytes.Buffer, error) {

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
