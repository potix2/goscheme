package vm

import (
	"fmt"

	"github.com/potix2/goscheme/ast"
)

type frame struct {
	Values map[string]ast.Expr
	next   *frame
	prev   *frame
}

type Env struct {
	head *frame
	tail *frame
}

func NewEnv() *Env {
	f := makeFrame(map[string]ast.Expr{})
	return &Env{head: &f, tail: &f}
}

func (e *Env) Lookup(name string) (ast.Expr, error) {
	f := e.head
	for f != nil {
		if expr, ok := f.Values[name]; ok {
			return expr, nil
		}
		f = f.next
	}
	return nil, &Error{Message: fmt.Sprintf("Unbound Variable: %s", name)}
}

func (e *Env) Extend(vals map[string]ast.Expr) error {
	f := makeFrame(vals)
	f.next = e.head
	e.head.prev = &f
	e.head = &f
	return nil
}

func (e *Env) SetVariable(name string, value ast.Expr) {
	e.head.Values[name] = value
}

func makeFrame(vals map[string]ast.Expr) frame {
	return frame{Values: vals}
}
