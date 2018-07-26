package vm

import (
	"fmt"

	"github.com/potix2/goscheme/scm"
)

var interactionEnvironment *scm.Env

func SetInteractionEnvironment(env *scm.Env) {
	interactionEnvironment = env
}

func getInteractionEnvironment() *scm.Env {
	return interactionEnvironment
}

func NewEnv() *scm.Env {
	return &scm.Env{Values: map[string]scm.Object{}}
}

func Lookup(env *scm.Env, name string) (scm.Object, error) {
	e := env
	for e != nil {
		if expr, ok := e.Values[name]; ok {
			return expr, nil
		}
		e = e.Parent
	}

	return nil, &Error{Message: fmt.Sprintf("Unbound Variable: %s", name)}
}

func Extend(env *scm.Env, vals map[string]scm.Object) *scm.Env {
	return &scm.Env{Values: vals, Parent: env}
}

func envInteractionEnvironment(args []scm.Object) (scm.Object, error) {
	if len(args) != 0 {
		return nil, &Error{Message: fmt.Sprintf("required no args, but got %d", len(args))}
	}
	return *getInteractionEnvironment(), nil
}

func envEval(args []scm.Object) (scm.Object, error) {
	if len(args) != 2 {
		return nil, &Error{Message: fmt.Sprintf("required 2, but got %d", len(args))}
	}

	expr := args[0]
	if env, ok := args[1].(scm.Env); ok {
		return Eval(expr, &env)
	} else {
		return nil, &Error{Message: fmt.Sprintf("expected env, but got %s", scm.TypeString(args[1]))}
	}
}

//(apply proc arg1 ... args)
//  => (proc (append (list arg1 ...) args))
func envApply(args []scm.Object) (scm.Object, error) {
	if len(args) < 2 {
		return nil, &Error{Message: fmt.Sprintf("required at lescm 2, but got %d", len(args))}
	}

	op := args[0]
	tail := args[len(args)-1]
	if !isList(tail) {
		return nil, &Error{Message: fmt.Sprintf("expected list, but got %s", scm.TypeString(tail))}
	}

	var vals []scm.Object
	if len(args) >= 3 {
		vals = append(args[1:len(args)-1], listToSlice(tail)...)
	} else {
		vals = listToSlice(tail)
	}
	return apply(op, vals)
}
