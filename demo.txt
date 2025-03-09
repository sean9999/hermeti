package main

import (
	"context"
	"fmt"
)

type IMac interface {
	String() string
}

func Dork(ctx context.Context, name string) {

	greeting, ok := ctx.Value("greeting").(string)
	if !ok {
		greeting = "hello, "
	}

	fmt.Println(greeting, " ", name)
}
