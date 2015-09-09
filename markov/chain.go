package markov

import (
	"io"
	"io/ioutil"
	"math/rand"
	"strings"
	"unicode"
)

type Prefix []string

func (p Prefix) Key() string {
	return strings.Join(p, string([]byte{0}))
}

type Chain struct {
	chain    map[string][]string
	prefixes []Prefix
	starters []Prefix
}

func (c *Chain) Add(prefix Prefix, word string) {
	key := prefix.Key()
	c.chain[key] = append(c.chain[key], word)
}

func (c *Chain) Generate(maxWords int) string {
	prefix := c.starters[rand.Int63n(int64(len(c.starters)))]

	outputWords := make([]string, len(prefix))
	copy(outputWords, []string(prefix))

	var stopIndexes []int
	for len(outputWords) < maxWords {
		suffixes := c.chain[prefix.Key()]
		if len(suffixes) == 0 {
			break
		}
		next := suffixes[rand.Intn(len(suffixes))]
		outputWords = append(outputWords, next)
		copy(prefix, prefix[1:])
		prefix[len(prefix)-1] = next
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

	return strings.Join(outputWords[:stopIdx], " ")
}

func NewChain(r io.Reader, prefixSize int) (c *Chain, err error) {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	words := strings.Fields(string(input))

	c = &Chain{}
	c.prefixes = make([]Prefix, 0)
	c.chain = make(map[string][]string)

	for i := 0; i < len(words)-prefixSize; i++ {
		prefix := Prefix(words[i : i+prefixSize])
		c.Add(prefix, words[i+prefixSize])
		c.prefixes = append(c.prefixes, prefix)
		if unicode.IsUpper(rune(prefix[0][0])) {
			c.starters = append(c.starters, prefix)
		}
	}

	if len(c.starters) == 0 {
		c.starters = c.prefixes
	}

	return c, nil
}
