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

	result = Eval(Parse("(* (+ 1 2 3) 2)"), &env)
	assert.Equal(t, result.(int), 12)
}

func TestLambda(t *testing.T) {
	env := make(Env)
	obj := Parse("(define func (lambda (x) (+ x 1)))")
	Eval(obj, &env)
	obj = Parse("(func 1)")
	result := Eval(obj, &env)
	assert.Equal(t, result.(int), 2)
}

func TestDef(t *testing.T) {
	env := make(Env)
	obj := Parse("(define hello 1)")
	Eval(obj, &env)
	assert.Equal(t, env["hello"].(int), 1)
	obj = Parse("(hello)")
	res := Eval(obj, &env)
	assert.Equal(t, res.(int), 1)

	res = Eval(Parse("(* 20 hello)"), &env)
	assert.Equal(t, res.(int), 20)
}
