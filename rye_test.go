package rye

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	tokens := tokenize("(+ 1 2)")
	assert.Equal(t, tokens[0].Type, tokenLParen)
	//assert.Equal(t, tokens[1], Token{"+", TokenSymbol})
	assert.Equal(t, tokens[2], token{1, tokenInt})
	assert.Equal(t, tokens[3], token{2, tokenInt})
	assert.Equal(t, tokens[4].Type, tokenRParen)
}

func TestEval(t *testing.T) {
	env := make(Env)
	obj := Parse("(+ 1 2 3)")
	result := Eval(obj, &env)
	assert.Equal(t, result.(int), 6)
}

func TestLambda(t *testing.T) {
	env := make(Env)
	obj := Parse("(define func (lambda (x) (+ x 1)))")
	Eval(obj, &env)
	obj = Parse("(func 1)")
	result := Eval(obj, &env)
	assert.Equal(t, result.(int), 2)
}
