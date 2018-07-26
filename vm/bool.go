package vm

import (
	"fmt"

	"github.com/potix2/goscheme/scm"
)

func boolNot(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}

	if b, ok := args[0].(scm.BooleanExpr); ok {
		if !b.Lit {
			return scm.BooleanExpr{true}, nil
		}
	}

	return scm.BooleanExpr{false}, nil
}

func boolIsBoolean(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}

	if _, ok := args[0].(scm.BooleanExpr); ok {
		return scm.BooleanExpr{true}, nil
	} else {
		return scm.BooleanExpr{false}, nil
	}
}

func boolIsProcedure(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}
	switch args[0].(type) {
	case scm.LambdaExpr, scm.PrimitiveProcExpr:
		return scm.BooleanExpr{true}, nil
	default:
		return scm.BooleanExpr{false}, nil
	}
}
