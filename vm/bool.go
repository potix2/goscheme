package vm

import (
	"fmt"

	"github.com/potix2/goscheme/scm"
)

func boolNot(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}

	if b, ok := args[0].(scm.Boolean); ok {
		if !b.Lit {
			return scm.Boolean{true}, nil
		}
	}

	return scm.Boolean{false}, nil
}

func boolIsBoolean(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}

	if _, ok := args[0].(scm.Boolean); ok {
		return scm.Boolean{true}, nil
	} else {
		return scm.Boolean{false}, nil
	}
}

func boolIsProcedure(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}
	switch args[0].(type) {
	case scm.Lambda, scm.PrimitiveProc:
		return scm.Boolean{true}, nil
	default:
		return scm.Boolean{false}, nil
	}
}
