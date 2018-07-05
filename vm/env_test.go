package vm

import (
	"github.com/potix2/goscheme/ast"
	"testing"
)

func TestSimple(t *testing.T) {
	env := NewEnv()
	e, err := Lookup(env, "a")
	if err == nil && e != nil {
		t.Fatalf("error: %#v", e)
	}
	env.Bind("a", ast.IntNum(10))
	e, err = Lookup(env, "a")
	if ue, ok := e.(ast.IntNum); ok {
		if ue != 10 {
			t.Fatalf("found unexpected variable: %#v\n", ue)
		}
	}

	env = Extend(env, map[string]ast.Expr{"a": ast.IntNum(20)})
	e, err = Lookup(env, "a")
	if ue, ok := e.(ast.IntNum); ok {
		if ue != 20 {
			t.Fatalf("found unexpected variable: %#v\n", ue)
		}
	}
}
