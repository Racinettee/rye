package rye

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	tokens := Tokenize("(+ 1 2)")
	assert.Equal(t, tokens[0].Type, TokenLParen)
	//assert.Equal(t, tokens[1], Token{"+", TokenSymbol})
	assert.Equal(t, tokens[2], Token{1, TokenInt})
	assert.Equal(t, tokens[3], Token{2, TokenInt})
	assert.Equal(t, tokens[4].Type, TokenRParen)
}

func TestEval(t *testing.T) {
	env := make(Env)
	obj := Parse("(+ 1 2 3)")
	result := Eval(obj, &env)
	assert.Equal(t, result.(int), 6)
}
