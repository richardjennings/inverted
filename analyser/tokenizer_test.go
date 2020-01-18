package analyser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// test tokenizer splits by whitespace
func TestTokenizer_Tokenize(t *testing.T) {
	tokenizer := NewTokenizer()
	have := tokenizer.Tokenize("1 2 3\n4\t5")
	want := []string{"1", "2", "3", "4", "5"}
	assert.Equal(t, want, have)
}
