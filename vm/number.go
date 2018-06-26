package vm

import (
	"fmt"

	"github.com/potix2/goscheme/ast"
)

func MakeInt(v int) ast.Uint10Expr {
	return ast.Uint10Expr{Lit: v}
}

func arithAdd(args []ast.Expr) (ast.Expr, error) {
	ret := MakeInt(0)
	for _, a := range args {
		if a0, ok := a.(ast.Uint10Expr); ok {
			ret = MakeInt(ret.Lit + a0.Lit)
		} else {
			return nil, &Error{Message: fmt.Sprintf("invalid number %#v", a)}
		}
	}
	return ret, nil
}

func arithSub(args []ast.Expr) (ast.Expr, error) {
	ret, ok := args[0].(ast.Uint10Expr)
	if !ok {
		return nil, &Error{Message: fmt.Sprintf("invalid number %#v", ret)}
	}

	for _, a := range args[1:] {
		if a0, ok := a.(ast.Uint10Expr); ok {
			ret = MakeInt(ret.Lit - a0.Lit)
		} else {
			return nil, &Error{Message: fmt.Sprintf("invalid number %#v", a)}
		}
	}
	return ret, nil
}

func arithMul(args []ast.Expr) (ast.Expr, error) {
	ret := MakeInt(1)
	for _, a := range args {
		if a0, ok := a.(ast.Uint10Expr); ok {
			ret = MakeInt(ret.Lit * a0.Lit)
		} else {
			return nil, &Error{Message: fmt.Sprintf("invalid number %#v", a)}
		}
	}
	return ret, nil
}

func arithDiv(args []ast.Expr) (ast.Expr, error) {
	return nil, &Error{Message: "not implemented"}
}
