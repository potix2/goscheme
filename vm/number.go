package vm

import (
	"fmt"

	"github.com/potix2/goscheme/ast"
)

func arithAdd(args []ast.Expr) (ast.Expr, error) {
	var ret ast.Number
	ret = ast.IntNum(0)
	for _, a := range args {
		if a0, ok := a.(ast.Number); ok {
			ret = ret.Add(a0)
		} else {
			return nil, &Error{Message: fmt.Sprintf("invalid number %#v", a)}
		}
	}
	return ret, nil
}

func arithSub(args []ast.Expr) (ast.Expr, error) {
	var ret ast.Number
	ret, ok := args[0].(ast.Number)
	if !ok {
		return nil, &Error{Message: fmt.Sprintf("invalid number %#v", ret)}
	}

	for _, a := range args[1:] {
		if a0, ok := a.(ast.Number); ok {
			ret = ret.Sub(a0)
		} else {
			return nil, &Error{Message: fmt.Sprintf("invalid number %#v", a)}
		}
	}
	return ret, nil
}

func arithMul(args []ast.Expr) (ast.Expr, error) {
	var ret ast.Number
	ret = ast.IntNum(1)
	for _, a := range args {
		if a0, ok := a.(ast.Number); ok {
			ret = ret.Mul(a0)
		} else {
			return nil, &Error{Message: fmt.Sprintf("invalid number %#v", a)}
		}
	}
	return ret, nil
}

func arithDiv(args []ast.Expr) (ast.Expr, error) {
	if len(args) == 0 {
		return nil, &Error{Message: "this procedure requires at least one argument"}
	}

	var ret ast.Number
	ret = args[0].(ast.Number)
	for _, a := range args[1:] {
		if a0, ok := a.(ast.Number); ok {
			ret = ret.Div(a0)
		} else {
			return nil, &Error{Message: fmt.Sprintf("invalid number %#v", a)}
		}
	}
	return ret, nil
}

func arithEqual(args []ast.Expr) (ast.Expr, error) {
	l := args[0].(ast.Number)
	r := args[1].(ast.Number)
	return ast.BooleanExpr{ast.EqNum(l, r)}, nil
}

func arithGreaterThan(args []ast.Expr) (ast.Expr, error) {
	l := args[0].(ast.Number)
	r := args[1].(ast.Number)
	if _, ok := l.(ast.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}
	if _, ok := r.(ast.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}

	return ast.BooleanExpr{ast.GTNum(l, r)}, nil
}

func arithLessThan(args []ast.Expr) (ast.Expr, error) {
	l := args[0].(ast.Number)
	r := args[1].(ast.Number)
	if _, ok := l.(ast.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}
	if _, ok := r.(ast.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}

	return ast.BooleanExpr{ast.LTNum(l, r)}, nil
}

func arithGreaterThanEuqal(args []ast.Expr) (ast.Expr, error) {
	l := args[0].(ast.Number)
	r := args[1].(ast.Number)
	if _, ok := l.(ast.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}
	if _, ok := r.(ast.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}

	return ast.BooleanExpr{ast.GTENum(l, r)}, nil
}

func arithLessThanEqual(args []ast.Expr) (ast.Expr, error) {
	l := args[0].(ast.Number)
	r := args[1].(ast.Number)
	if _, ok := l.(ast.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}
	if _, ok := r.(ast.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}

	return ast.BooleanExpr{ast.LTENum(l, r)}, nil
}

func implIsNumber(expr ast.Expr) bool {
	switch expr.(type) {
	case ast.IntNum, ast.RealNum, ast.RatNum, ast.CompNum:
		return true
	default:
		return false
	}
}

func arithIsNumber(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("requires 1, but got %d", len(args))}
	}
	return ast.BooleanExpr{implIsNumber(args[0])}, nil
}

func arithNumberToString(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("requires 1, but got %d", len(args))}
	}
	if implIsNumber(args[0]) {
		return ast.NumberToString(args[0]), nil
	} else {
		return nil, &Error{Message: fmt.Sprintf("expected number, but got %s", ast.TypeString(args[0]))}
	}
}

func arithStringToNumber(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("requires 1, but got %d", len(args))}
	}
	if s, ok := args[0].(ast.StringExpr); ok {
		return ast.StringToNumber(string(s)), nil
	} else {
		return nil, &Error{Message: fmt.Sprintf("expected string, but got %s", ast.TypeString(args[0]))}
	}
}
