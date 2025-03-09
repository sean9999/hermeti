package hermeti

import (
	"context"
	"io"

	"github.com/sean9999/pear"
)

// a Runner takes an [Env] and runs some code against it
type Runner[T any] interface {
	Run(Env)
}

// a CLI is a command line interface.
type CLI[T any] struct {
	Env Env
	Cmd Runner[T]
}

// Run() runs the Runners Run method, passing in Env.
// It's simly a convenience function.
func (cli CLI[T]) Run(ctx context.Context) {
	cli.Cmd.Run(cli.Env)
}

// Obj exposes our application object
func (cli CLI[T]) Obj() T {
	return cli.Cmd.(T)
}

var ErrOutputNotReadable = pear.Defer("output stream is not readable")

// OutStream returns an io.Reader representing the stuff you put in StdOut.
//
//	This will not work in a real CLI because os.StdOut is not readable
func (cli CLI[T]) OutStream() (io.Reader, error) {

	o, ok := cli.Env.OutStream.(io.Reader)
	if !ok {
		return nil, ErrOutputNotReadable
	}

	return o, nil

}
