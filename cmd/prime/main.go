package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/sean9999/hermeti"
	"github.com/sean9999/pear"
)

/**
 * A CLI that gives itself w (wait time) to compute the nth prime, using the least efficient and most naive algorithm possible
 **/

type exe struct {
	primes []int
	n      int
	w      time.Duration
	sync.RWMutex
}

func NewExe() *exe {
	//	keep a list of primes
	primes := make([]int, 2, 1024)
	primes[0] = 1
	primes[1] = 2
	return &exe{
		primes: primes,
	}
}

type result struct {
	prime int
	took  time.Duration
}

func calculatePrime(res chan result, exe *exe, start time.Time) {

	n := exe.n

	primes := exe.primes

	//	if our prime is already known, just return it
	if len(primes) > n {
		diff := time.Now().Sub(start)
		res <- result{primes[n], diff}
		return
	}

	var candidate int

	//	else start from the last known prime, incrementing by one, and trying to mod all lower primes
	//	if we find a prime, we add it, and keep going until len(primes) == n
outer:
	for candidate = primes[len(primes)-1] + 1; len(primes) < exe.n+1; candidate++ {
		for _, p := range exe.primes[1:] {
			if candidate%p == 0 {
				continue outer
			}
		}
		primes = append(primes, candidate)
	}
	nth_prime := primes[len(primes)-1]
	diff := time.Now().Sub(start)

	res <- result{nth_prime, diff}

}

func parseArgs(exe *exe, args []string) error {

	f := flag.NewFlagSet("flags", flag.PanicOnError)
	n := f.Int("n", 1, "the nth prime you want to calculate")
	w := f.Duration("w", time.Second*1, "how long to wait before timing out")
	f.Parse(args)

	exe.w = *w

	//fmt.Printf("%v\n", exe.w)

	if n == nil || *n < 2 {
		return pear.New("not possible")
	}
	exe.n = *n
	return nil

}

func (exe *exe) Run(env hermeti.Env) {
	ctx := context.WithValue(context.Background(), "start", time.Now())
	args := env.Args[1:]
	err := parseArgs(exe, args)

	ctx, cancel := context.WithTimeout(ctx, exe.w)
	defer cancel()

	if err != nil {
		fmt.Fprintf(env.ErrStream, "%v\n", err)
		return
	}
	ch := make(chan result)

	go calculatePrime(ch, exe, ctx.Value("start").(time.Time))

	select {
	case <-ctx.Done():
		fmt.Fprintln(env.ErrStream, ctx.Err())
	case res := <-ch:
		o := fmt.Sprintf("The %dth prime is %d, and it took %v.", exe.n, res.prime, res.took)
		fmt.Fprintln(env.OutStream, o)
	}

}

func (exe *exe) State() *exe {
	return exe
}

func main() {
	cli := hermeti.NewRealCli[*exe](NewExe())
	cli.Run(context.Background())
}
