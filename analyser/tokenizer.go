package analyser

import (
	"strings"
)

type Tokenizer struct{}

// Tokenize splits a string into an array of strings using
// white space characters, as defined by unicode.IsSpace
func (t Tokenizer) Tokenize(term string) (result []string) {
	result = strings.Fields(term)
	return result
}

// NewTokenizer returns a new Tokenizer struct
func NewTokenizer() Tokenizer {
	return Tokenizer{}
}
