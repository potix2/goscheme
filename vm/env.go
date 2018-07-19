package vm

import (
	"fmt"

	"github.com/potix2/goscheme/ast"
)

var interactionEnvironment *ast.Env

func SetInteractionEnvironment(env *ast.Env) {
	interactionEnvironment = env
}

func getInteractionEnvironment() *ast.Env {
	return interactionEnvironment
}

func NewEnv() *ast.Env {
	return &ast.Env{Values: map[string]ast.Expr{}}
}

func Lookup(env *ast.Env, name string) (ast.Expr, error) {
	e := env
	for e != nil {
		if expr, ok := e.Values[name]; ok {
			return expr, nil
		}
		e = e.Parent
	}

	return nil, &Error{Message: fmt.Sprintf("Unbound Variable: %s", name)}
}

func Extend(env *ast.Env, vals map[string]ast.Expr) *ast.Env {
	return &ast.Env{Values: vals, Parent: env}
}

func envInteractionEnvironment(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 0 {
		return nil, &Error{Message: fmt.Sprintf("required no args, but got %d", len(args))}
	}
	return *getInteractionEnvironment(), nil
}

func envEval(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 2 {
		return nil, &Error{Message: fmt.Sprintf("required 2, but got %d", len(args))}
	}

	expr := args[0]
	if env, ok := args[1].(ast.Env); ok {
		return Eval(expr, &env)
	} else {
		return nil, &Error{Message: fmt.Sprintf("expected env, but got %s", ast.TypeString(args[1]))}
	}
}

//(apply proc arg1 ... args)
//  => (proc (append (list arg1 ...) args))
func envApply(args []ast.Expr) (ast.Expr, error) {
	if len(args) < 2 {
		return nil, &Error{Message: fmt.Sprintf("required at least 2, but got %d", len(args))}
	}

	op := args[0]
	tail := args[len(args)-1]
	if !isList(tail) {
		return nil, &Error{Message: fmt.Sprintf("expected list, but got %s", ast.TypeString(tail))}
	}

	var vals []ast.Expr
	if len(args) >= 3 {
		vals = append(args[1:len(args)-1], listToSlice(tail)...)
	} else {
		vals = listToSlice(tail)
	}
	return apply(op, vals)
}
