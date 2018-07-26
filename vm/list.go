package vm

import (
	"fmt"

	"github.com/potix2/goscheme/scm"
)

func listCons(args []scm.Object) (scm.Object, error) {
	return scm.PairExpr{args[0], args[1]}, nil
}

func listCar(args []scm.Object) (scm.Object, error) {
	if p, ok := args[0].(scm.PairExpr); ok {
		return p.Car, nil
	}
	return nil, &Error{Message: fmt.Sprintf("pair required, but got %#v", args[0])}
}

func listCdr(args []scm.Object) (scm.Object, error) {
	if p, ok := args[0].(scm.PairExpr); ok {
		return p.Cdr, nil
	}
	return nil, &Error{Message: fmt.Sprintf("pair required, but got %#v", args[0])}
}

func recMakeListFromSlice(elems []scm.Object) scm.Object {
	if len(elems) == 0 {
		return scm.PairExpr{}
	} else {
		return scm.PairExpr{elems[0], recMakeListFromSlice(elems[1:])}
	}
}

func listList(args []scm.Object) (scm.Object, error) {
	return recMakeListFromSlice(args), nil
}

func listIsPair(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}
	if _, ok := args[0].(scm.PairExpr); ok {
		return scm.BooleanExpr{true}, nil
	}
	return scm.BooleanExpr{false}, nil
}

func listIsNull(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}
	if p, ok := args[0].(scm.PairExpr); ok {
		if p.IsEmptyList() {
			return scm.BooleanExpr{true}, nil
		}
	}
	return scm.BooleanExpr{false}, nil
}

func listIsList(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("not requires 1, but got %d", len(args))}
	}
	if p, ok := args[0].(scm.PairExpr); ok {
		return scm.BooleanExpr{p.IsList()}, nil
	}
	return scm.BooleanExpr{false}, nil
}

func recListToSlice(p scm.PairExpr, ret []scm.Object) []scm.Object {
	if p.IsEmptyList() {
		return ret
	} else {
		return recListToSlice(p.Cdr.(scm.PairExpr), append(ret, p.Car))
	}
}

func listToSlice(head scm.Object) []scm.Object {
	var ret []scm.Object
	return recListToSlice(head.(scm.PairExpr), ret)
}
