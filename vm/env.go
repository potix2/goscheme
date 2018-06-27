package vm

import (
	"fmt"

	"github.com/potix2/goscheme/ast"
)

func NewEnv() *ast.Env {
	return &ast.Env{Values: map[string]ast.Expr{}}
}

func Lookup(env *ast.Env, name string) (ast.Expr, error) {
	e := env
	for e != nil {
		if expr, ok := e.Values[name]; ok {
			return expr, nil
		}
		e = env.Parent
	}

	return nil, &Error{Message: fmt.Sprintf("Unbound Variable: %s", name)}
}

func Extend(env *ast.Env, vals map[string]ast.Expr) *ast.Env {
	return &ast.Env{Values: vals, Parent: env}
}
