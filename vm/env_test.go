package vm

import (
	"github.com/potix2/goscheme/ast"
	"testing"
)

func TestSimple(t *testing.T) {
	env := NewEnv()
	e, err := env.Lookup("a")
	if err == nil && e != nil {
		t.Fatalf("error: %#v", e)
	}
	env.SetVariable("a", ast.Uint10Expr{Lit: 10})
	e, err = env.Lookup("a")
	if ue, ok := e.(ast.Uint10Expr); ok {
		if ue.Lit != 10 {
			t.Fatalf("found unexpected variable: %#v\n", ue)
		}
	}

	env.Extend(map[string]ast.Expr{"a": ast.Uint10Expr{Lit: 20}})
	e, err = env.Lookup("a")
	if ue, ok := e.(ast.Uint10Expr); ok {
		if ue.Lit != 20 {
			t.Fatalf("found unexpected variable: %#v\n", ue)
		}
	}
}
