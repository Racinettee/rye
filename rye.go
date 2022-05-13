package rye

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Racinettee/generics"
)

type Symbol string
type Object interface{}
type Void struct{}
type Env map[Symbol]Object
type Lambda struct {
	args []Symbol
	body []Object
}
type tokenType byte
type token struct {
	Symbol interface{}
	Type   tokenType
}

const (
	tokenInt tokenType = iota
	tokenSymbol
	tokenLParen
	tokenRParen
)

func tokenize(program string) []token {
	words := strings.Fields(strings.ReplaceAll(strings.ReplaceAll(program, "(", " ( "), ")", " ) "))
	var result []token
	for _, word := range words {
		switch word {
		case "(":
			result = append(result, token{"(", tokenLParen})
		case ")":
			result = append(result, token{")", tokenRParen})
		default:
			if i, err := strconv.Atoi(word); err != nil {
				result = append(result, token{Symbol(word), tokenSymbol})
			} else {
				result = append(result, token{i, tokenInt})
			}
		}
	}
	return result
}

func parseTokens(tokens *generics.Queue[token]) ([]Object, error) {
	var result generics.List[Object]
	token := tokens.Pop()
	if token.Type != tokenLParen {
		return result, fmt.Errorf("expected ( but found %+v", token)
	}
	for len(*tokens) != 0 {
		token = tokens.Front()
		switch token.Type {
		case tokenInt, tokenSymbol:
			tokens.Pop()
			result.Push(token.Symbol)
		case tokenLParen:
			subList, err := parseTokens(tokens)
			if err != nil {
				return result, err
			}
			result.Push(subList)
		case tokenRParen:
			tokens.Pop()
			return result, nil
		}
	}
	return result, nil
}
func Parse(program string) Object {
	tokens := generics.Queue[token](tokenize(program))
	object, err := parseTokens(&tokens)
	if err != nil {
		return err
	}
	return object
}

func (env Env) Clone() Env {
	result := make(Env)
	for k, v := range env {
		result[k] = v
	}
	return result
}

func Eval(obj Object, env *Env) Object {
	switch obj := obj.(type) {
	case error, int, bool:
		return obj
	case Symbol:
		return evalSym(obj, env)
	case []Object:
		return evalList(obj, env)
	case Void:
		return Void{}
	}
	return nil
}
func evalSym(sym Symbol, env *Env) Object {
	if obj, ok := (*env)[sym]; ok {
		return obj
	}
	return nil
}
func evalList(list []Object, env *Env) Object {
	first := list[0]
	switch first := first.(type) {
	case Symbol:
		switch first {
		case "+", "-", "*", "/", "<", ">", "=", "!=":
			return evalBinop(list, env)
		case "define":
			return evalDefine(list, env)
		case "if":
			return evalIf(list, env)
		case "lambda":
			return evalFnDefine(list, env)
		default:
			return evalFnCall(first, list, env)
		}
	}
	var result []Object
	for _, obj := range list {
		subResult := Eval(obj, env)
		switch subResult.(type) {
		case Void:
			continue
		default:
			result = append(result, subResult)
		}
	}
	return result
}
func evalBinop(list []Object, env *Env) Object {
	if len(list) < 3 {
		return fmt.Errorf("invalid number of arguments")
	}
	switch operator := list[0]; operator.(type) {
	case Symbol:
		switch operator.(Symbol) {
		case "+":
			sum := 0
			for _, obj := range list[1:] {
				res, ok := Eval(obj, env).(int)
				if !ok {
					return fmt.Errorf("error value didnt evaluate to int")
				}
				sum += res
			}
			return sum
		case "-":
			first, ok := Eval(list[1], env).(int)
			if !ok {
				return fmt.Errorf("error element in - was not an int")
			}
			for _, obj := range list[2:] {
				res, ok := Eval(obj, env).(int)
				if !ok {
					return fmt.Errorf("error value in list evaluating - was not int")
				}
				first -= res
			}
			return first
		case "*":
			product := 1
			for _, obj := range list[1:] {
				res, ok := Eval(obj, env).(int)
				if !ok {
					return fmt.Errorf("error evaluating * expected int")
				}
				product *= res
			}
			return product
		case "/":
			divisor, ok := Eval(list[1], env).(int)
			if !ok {
				return fmt.Errorf("expected int evaluating /")
			}
			for _, obj := range list[2:] {
				res, ok := Eval(obj, env).(int)
				if !ok {
					return fmt.Errorf("expected int evaluating /")
				}
				divisor /= res
			}
			return divisor
		}
		left, ok1 := Eval(list[1], env).(int)
		right, ok2 := Eval(list[2], env).(int)
		if !(ok1 && ok2) {
			return fmt.Errorf("error evaluating comparitor expected int")
		}
		switch operator.(Symbol) {
		case "<":
			return left < right
		case ">":
			return left > right
		case "=":
			return left == right
		case "!=":
			return left != right
		}
	default:
	}
	return fmt.Errorf("operator must be symbol")
}

func evalDefine(list []Object, env *Env) Object {
	if len(list) != 3 {
		return fmt.Errorf("invalid number of arguments supplied to define")
	}
	if sym, ok := list[1].(Symbol); ok {
		(*env)[sym] = Eval(list[2], env)
	}
	return Void{}
}

func evalIf(list []Object, env *Env) Object {
	if len(list) != 4 {
		return fmt.Errorf("invalid number of arguments for if")
	}
	cond, ok := Eval(list[1], env).(bool)
	if !ok {
		return fmt.Errorf("conditional does not evaluate to bool")
	} else if cond {
		return Eval(list[2], env)
	}
	return Eval(list[3], env)
}

func evalFnDefine(list []Object, env *Env) Object {
	parameters, ok := list[1].([]Object)
	if !ok {
		return fmt.Errorf("invalid function parameters expected list")
	}
	var params []Symbol
	for _, param := range parameters {
		if s, ok := param.(Symbol); ok {
			params = append(params, s)
		} else {
			return fmt.Errorf("arguments for lamba must all be symbol")
		}
	}
	if body, ok := list[2].([]Object); ok {
		return Lambda{params, body}
		// not sure this is correct
	}
	return fmt.Errorf("expected list for lambda body")
}

func evalFnCall(fnname Symbol, list []Object, env *Env) Object {
	switch val := (*env)[fnname].(type) {
	case Lambda:
		nestedEnv := env.Clone()
		for i, param := range val.args {
			res := Eval(list[i+1], env)
			nestedEnv[param] = res
		}
		return Eval(val.body, &nestedEnv)
	default:
		return Eval(val, env)
	}
}
