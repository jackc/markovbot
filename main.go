package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

var options struct {
	prefixSize    int
	maxOutputSize int
	seed          int64
}

type Prefix []string

func (p Prefix) Key() string {
	return strings.Join(p, string([]byte{0}))
}

type Chain map[string][]string

func (c Chain) Add(prefix Prefix, word string) {
	key := prefix.Key()
	c[key] = append(c[key], word)
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

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	words := strings.Fields(string(input))

	prefixes := make([]Prefix, 0)
	chain := make(Chain)

	for i := 0; i < len(words)-options.prefixSize; i++ {
		prefix := Prefix(words[i : i+options.prefixSize])
		chain.Add(prefix, words[i+options.prefixSize])
		prefixes = append(prefixes, prefix)
	}

	var prefix Prefix
	for {
		prefix = prefixes[rand.Int63n(int64(len(prefixes)))]
		if unicode.IsUpper(rune(prefix[0][0])) {
			break
		}
	}

	outputWords := make([]string, len(prefix))
	copy(outputWords, []string(prefix))

	var stopIndexes []int
	for len(outputWords) < options.maxOutputSize {
		suffixes := chain[prefix.Key()]
		next := suffixes[rand.Intn(len(suffixes))]
		outputWords = append(outputWords, next)
		copy(prefix, prefix[1:])
		prefix[options.prefixSize-1] = next
		if strings.HasSuffix(next, ".") || strings.HasSuffix(next, "?") || strings.HasSuffix(next, "!") {
			stopIndexes = append(stopIndexes, len(outputWords))
		}
	}

	var stopIdx int
	if len(stopIndexes) > 0 {
		stopIdx = stopIndexes[rand.Intn(len(stopIndexes))]
	} else {
		stopIdx = len(outputWords)
	}

	fmt.Println(strings.Join(outputWords[:stopIdx], " "))
}
