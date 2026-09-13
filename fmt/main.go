/*
Package fmt is like fmt from the standard library,
but suitable for use in a hermeti.CLI. It "shadows"
that package, but writes go to the writer you specify,
rather than assuming os.StdOut.
This keeps your CLI hermetic while allowing the use of familiar functions.

Here's how you can use it:

	func (a *myApp) Run(env *hermeti.Env) {

		fmt.SetOutput(env.OutputStream)

		fmt.Println("Hello, world.")
	}
*/
package fmt

import (
	"fmt"
	"io"
)

var stdout io.Writer
var stdin io.Reader

// SetOutput sets the default [io.Writer] for all functions
// that would normally assume [os.StdOut].
func SetOutput(w io.Writer) {
	stdout = w
}

// SetInput sets the default [io.Reader] for all functions
// that would normally assume [os.StdIn].
func SetInput(r io.Reader) {
	stdin = r
}

func Print(things ...any) (int, error) {
	return fmt.Fprint(stdout, things...)
}

func Println(things ...any) (int, error) {
	return fmt.Fprintln(stdout, things...)
}

func Printf(format string, a ...any) (int, error) {
	return fmt.Fprintf(stdout, format, a...)
}

func Scan(a ...any) (int, error) {
	return fmt.Fscan(stdin, a...)
}

func Scanf(format string, a ...any) (int, error) {
	return fmt.Fscanf(stdin, format, a...)
}

func Scanln(a ...any) (int, error) {
	return fmt.Fscanln(stdin, a...)
}

//go:generate go tool pkgalias fmt .
