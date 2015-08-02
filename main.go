package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/jackc/markovbot/markov"
)

const Version = "0.0.2"

var options struct {
	prefixSize    int
	maxOutputSize int
	seed          int64
	filePath      string
	httpAddr      string
	version       bool
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.IntVar(&options.prefixSize, "prefix", 2, "prefix size")
	flag.IntVar(&options.maxOutputSize, "output", 200, "max output size in words")
	flag.Int64Var(&options.seed, "seed", -1, "seed for random number generator")
	flag.StringVar(&options.filePath, "file", "", "source file (if not provided will read from stdin)")
	flag.StringVar(&options.httpAddr, "http", "", "HTTP listen address (e.g. 127.0.0.1:3000)")
	flag.BoolVar(&options.version, "version", false, "print version and exit")

	flag.Parse()

	if options.version {
		fmt.Printf("markovbot v%v\n", Version)
		os.Exit(0)
	}

	if options.seed < 0 {
		options.seed = time.Now().UnixNano()
		fmt.Fprintln(os.Stderr, "seed:", options.seed)
	}

	rand.Seed(options.seed)

	var err error
	var file *os.File
	var in io.Reader
	if options.filePath != "" {
		file, err = os.Open(options.filePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		in = file
	} else {
		in = os.Stdin
	}

	chain, err := markov.NewChain(in, options.prefixSize)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if file != nil {
		file.Close()
	}

	if options.httpAddr != "" {
		fmt.Fprintln(os.Stderr, "Listening on:", options.httpAddr)

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var jsonResp struct {
				Text string `json:"text"`
			}
			jsonResp.Text = chain.Generate(options.maxOutputSize)
			js, err := json.Marshal(jsonResp)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			w.Write(js)
		})

		err = http.ListenAndServe(options.httpAddr, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		fmt.Println(chain.Generate(options.maxOutputSize))
	}

}
