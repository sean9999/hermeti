package fmt

import (
	"fmt"
	"io"
)

var stdout io.Writer

func SetOutput(w io.Writer) {
	stdout = w
}

func Println(things ...any) {
	fmt.Fprintln(stdout, things...)
}

func Printf(format string, a ...any) {
	fmt.Fprintf(stdout, format, a...)
}

func Errorf(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

//go:generate pkgalias fmt .
