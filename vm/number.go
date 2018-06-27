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

func arithGreaterThan(args []ast.Expr) (ast.Expr, error) {
	l := args[0].(ast.Uint10Expr)
	r := args[1].(ast.Uint10Expr)
	return ast.BooleanExpr{l.Lit > r.Lit}, nil
}

func arithLessThan(args []ast.Expr) (ast.Expr, error) {
	l := args[0].(ast.Uint10Expr)
	r := args[1].(ast.Uint10Expr)
	return ast.BooleanExpr{l.Lit < r.Lit}, nil
}

func arithIsNumber(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}
	switch args[0].(type) {
	case ast.Uint10Expr:
		return ast.BooleanExpr{true}, nil
	default:
		return ast.BooleanExpr{false}, nil
	}
}
