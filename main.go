package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pop_markov/markov"
)

var options struct {
	prefixSize    int
	maxOutputSize int
	seed          int64
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.IntVar(&options.prefixSize, "prefix", 2, "prefix size")
	flag.IntVar(&options.maxOutputSize, "output", 200, "max output size")
	flag.Int64Var(&options.seed, "seed", -1, "seed for random number generator")

	flag.Parse()

	if options.seed < 0 {
		options.seed = time.Now().UnixNano()
		fmt.Fprintln(os.Stderr, "seed:", options.seed)
	}

	rand.Seed(options.seed)

	chain, err := markov.NewChain(os.Stdin, options.prefixSize)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(chain.Generate(options.maxOutputSize))
}
