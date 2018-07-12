package vm

import (
	"bytes"
	"fmt"

	"github.com/potix2/goscheme/ast"
)

func implIsString(expr ast.Expr) bool {
	_, ok := expr.(ast.StringExpr)
	return ok
}

func strIsString(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("requires 1, but got %d", len(args))}
	}
	return ast.BooleanExpr{implIsString(args[0])}, nil
}

func strStringLength(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("requires 1, but got %d", len(args))}
	}
	if s, ok := args[0].(ast.StringExpr); ok {
		return ast.IntNum(len(s)), nil
	} else {
		return nil, &Error{Message: fmt.Sprintf("type mismatch: expected string, but got %s", ast.TypeString(args[0]))}
	}
}

type strComparator func(ast.StringExpr, ast.StringExpr) bool

func compForAll(args []ast.Expr, comp strComparator) (ast.Expr, error) {
	if len(args) == 0 {
		return nil, &Error{Message: "wrong number of arguments"}
	}
	a := args[0].(ast.StringExpr)
	for _, b := range args[1:] {
		if bs, ok := b.(ast.StringExpr); !ok || !comp(a, bs) {
			return ast.BooleanExpr{false}, nil
		}
	}
	return ast.BooleanExpr{true}, nil
}

func strStringEqual(args []ast.Expr) (ast.Expr, error) {
	return compForAll(args, func(a, b ast.StringExpr) bool { return a == b })
}

func strStringLT(args []ast.Expr) (ast.Expr, error) {
	return compForAll(args, func(a, b ast.StringExpr) bool { return a < b })
}

func strStringGT(args []ast.Expr) (ast.Expr, error) {
	return compForAll(args, func(a, b ast.StringExpr) bool { return a > b })
}

func strStringLTE(args []ast.Expr) (ast.Expr, error) {
	return compForAll(args, func(a, b ast.StringExpr) bool { return a <= b })
}

func strStringGTE(args []ast.Expr) (ast.Expr, error) {
	return compForAll(args, func(a, b ast.StringExpr) bool { return a >= b })
}

func strSubstring(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 3 {
		return nil, &Error{Message: fmt.Sprintf("required 3, but got %d", len(args))}
	}
	var s ast.StringExpr
	var ok bool
	var start, end ast.IntNum
	if s, ok = args[0].(ast.StringExpr); !ok {
		return nil, &Error{Message: fmt.Sprintf("string required, but got %s", ast.TypeString(args[0]))}
	}
	if start, ok = args[1].(ast.IntNum); !ok {
		return nil, &Error{Message: fmt.Sprintf("integer required, but got %s", ast.TypeString(args[1]))}
	}
	if end, ok = args[2].(ast.IntNum); !ok {
		return nil, &Error{Message: fmt.Sprintf("integer required, but got %s", ast.TypeString(args[2]))}
	}
	if start < 0 || len(s) <= int(end) {
		return nil, &Error{Message: fmt.Sprintf("out of range: %d %d", start, end)}
	}
	return ast.StringExpr(s[start:end]), nil
}

func strStringAppend(args []ast.Expr) (ast.Expr, error) {
	if len(args) == 0 {
		return nil, &Error{Message: fmt.Sprintf("required at least 1, but got %d", len(args))}
	}
	var buffer bytes.Buffer
	for _, e := range args {
		if s, ok := e.(ast.StringExpr); !ok {
			return nil, &Error{Message: fmt.Sprintf("expected string, but got %s", ast.TypeString(e))}
		} else {
			buffer.WriteString(string(s))
		}
	}
	return ast.StringExpr(buffer.String()), nil
}
