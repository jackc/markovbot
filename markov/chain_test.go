package markov

import (
	"bytes"
	"testing"
)

func TestGenerate(t *testing.T) {
	sourceText := "This is a test. A test is a terror to students. The students were a terror to their teacher."
	buf := bytes.NewBufferString(sourceText)
	chain, err := NewChain(buf, 2)
	if err != nil {
		t.Fatal(err)
	}

	utterance := chain.Generate(50)
	if utterance == "" {
		t.Fatal("Generate returned empty string")
	}
}

func TestGenerateWithAllLowerCase(t *testing.T) {
	sourceText := "this is a test of a bot"
	buf := bytes.NewBufferString(sourceText)
	chain, err := NewChain(buf, 2)
	if err != nil {
		t.Fatal(err)
	}

	utterance := chain.Generate(3)
	if utterance == "" {
		t.Fatal("Generate returned empty string")
	}
}

func TestGenerateWhenReachingEndOfChain(t *testing.T) {
	sourceText := "this is a test of a bot"
	buf := bytes.NewBufferString(sourceText)
	chain, err := NewChain(buf, 2)
	if err != nil {
		t.Fatal(err)
	}

	utterance := chain.Generate(999999)
	if utterance == "" {
		t.Fatal("Generate returned empty string")
	}
}
