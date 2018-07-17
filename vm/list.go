package vm

import (
	"fmt"

	"github.com/potix2/goscheme/ast"
)

func listCons(args []ast.Expr) (ast.Expr, error) {
	return ast.PairExpr{args[0], args[1]}, nil
}

func listCar(args []ast.Expr) (ast.Expr, error) {
	if p, ok := args[0].(ast.PairExpr); ok {
		return p.Car, nil
	}
	return nil, &Error{Message: fmt.Sprintf("pair required, but got %#v", args[0])}
}

func listCdr(args []ast.Expr) (ast.Expr, error) {
	if p, ok := args[0].(ast.PairExpr); ok {
		return p.Cdr, nil
	}
	return nil, &Error{Message: fmt.Sprintf("pair required, but got %#v", args[0])}
}

func recMakeListFromSlice(elems []ast.Expr) ast.Expr {
	if len(elems) == 0 {
		return ast.PairExpr{}
	} else {
		return ast.PairExpr{elems[0], recMakeListFromSlice(elems[1:])}
	}
}

func listList(args []ast.Expr) (ast.Expr, error) {
	return recMakeListFromSlice(args), nil
}

func listIsPair(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}
	if _, ok := args[0].(ast.PairExpr); ok {
		return ast.BooleanExpr{true}, nil
	}
	return ast.BooleanExpr{false}, nil
}

func listIsNull(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}
	if p, ok := args[0].(ast.PairExpr); ok {
		if p.IsEmptyList() {
			return ast.BooleanExpr{true}, nil
		}
	}
	return ast.BooleanExpr{false}, nil
}

func listIsList(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}
	if p, ok := args[0].(ast.PairExpr); ok {
		return ast.BooleanExpr{p.IsList()}, nil
	}
	return ast.BooleanExpr{false}, nil
}

func recListToSlice(p ast.PairExpr, ret []ast.Expr) []ast.Expr {
	if p.IsEmptyList() {
		return ret
	} else {
		return recListToSlice(p.Cdr.(ast.PairExpr), append(ret, p.Car))
	}
}

func listToSlice(head ast.Expr) []ast.Expr {
	var ret []ast.Expr
	return recListToSlice(head.(ast.PairExpr), ret)
}
