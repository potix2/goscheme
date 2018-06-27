package vm

import (
	"fmt"

	"github.com/potix2/goscheme/ast"
)

func boolNot(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}

	if b, ok := args[0].(ast.BooleanExpr); ok {
		if !b.Lit {
			return ast.BooleanExpr{true}, nil
		}
	}

	return ast.BooleanExpr{false}, nil
}

func boolIsBoolean(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}

	if _, ok := args[0].(ast.BooleanExpr); ok {
		return ast.BooleanExpr{true}, nil
	} else {
		return ast.BooleanExpr{false}, nil
	}
}

func boolIsProcedure(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}
	switch args[0].(type) {
	case ast.LambdaExpr, ast.PrimitiveProcExpr:
		return ast.BooleanExpr{true}, nil
	default:
		return ast.BooleanExpr{false}, nil
	}
}
