package vm

import (
	"github.com/potix2/goscheme/scm"
	"testing"
)

func TestSimple(t *testing.T) {
	env := NewEnv()
	//{}
	e, err := Lookup(env, "a")
	if err == nil && e != nil {
		t.Fatalf("error: %#v", e)
	}

	//{a: 10, x: 100}
	env.Bind("a", scm.IntNum(10))
	env.Bind("x", scm.IntNum(100))
	e, err = Lookup(env, "a")
	if ue, ok := e.(scm.IntNum); ok {
		if ue != 10 {
			t.Fatalf("found unexpected variable: %#v\n", ue)
		}
	}

	//{a: 20, {a: 10, x: 100}}
	env = Extend(env, map[string]scm.Object{"a": scm.IntNum(20)})
	e, err = Lookup(env, "a")
	if ue, ok := e.(scm.IntNum); ok {
		if ue != 20 {
			t.Fatalf("found unexpected variable: %#v\n", ue)
		}
	}
	//{b: 30, {a: 20, {a: 10, x: 100}}}
	env = Extend(env, map[string]scm.Object{"b": scm.IntNum(30)})
	e, err = Lookup(env, "a")
	if ue, ok := e.(scm.IntNum); ok {
		if ue != 20 {
			t.Fatalf("found unexpected variable: %#v\n", ue)
		}
	}
	e, err = Lookup(env, "b")
	if ue, ok := e.(scm.IntNum); ok {
		if ue != 30 {
			t.Fatalf("found unexpected variable: %#v\n", ue)
		}
	}
	e, err = Lookup(env, "x")
	if ue, ok := e.(scm.IntNum); ok {
		if ue != 100 {
			t.Fatalf("found unexpected variable: %#v\n", ue)
		}
	}
}
