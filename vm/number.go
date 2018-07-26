package vm

import (
	"fmt"

	"github.com/potix2/goscheme/scm"
)

func arithAdd(args []scm.Object) (scm.Object, error) {
	var ret scm.Number
	ret = scm.IntNum(0)
	for _, a := range args {
		if a0, ok := a.(scm.Number); ok {
			ret = ret.Add(a0)
		} else {
			return nil, &Error{Message: fmt.Sprintf("invalid number %#v", a)}
		}
	}
	return ret, nil
}

func arithSub(args []scm.Object) (scm.Object, error) {
	var ret scm.Number
	ret, ok := args[0].(scm.Number)
	if !ok {
		return nil, &Error{Message: fmt.Sprintf("invalid number %#v", ret)}
	}

	for _, a := range args[1:] {
		if a0, ok := a.(scm.Number); ok {
			ret = ret.Sub(a0)
		} else {
			return nil, &Error{Message: fmt.Sprintf("invalid number %#v", a)}
		}
	}
	return ret, nil
}

func arithMul(args []scm.Object) (scm.Object, error) {
	var ret scm.Number
	ret = scm.IntNum(1)
	for _, a := range args {
		if a0, ok := a.(scm.Number); ok {
			ret = ret.Mul(a0)
		} else {
			return nil, &Error{Message: fmt.Sprintf("invalid number %#v", a)}
		}
	}
	return ret, nil
}

func arithDiv(args []scm.Object) (scm.Object, error) {
	if len(args) == 0 {
		return nil, &Error{Message: "this procedure requires at lescm one argument"}
	}

	var ret scm.Number
	ret = args[0].(scm.Number)
	for _, a := range args[1:] {
		if a0, ok := a.(scm.Number); ok {
			ret = ret.Div(a0)
		} else {
			return nil, &Error{Message: fmt.Sprintf("invalid number %#v", a)}
		}
	}
	return ret, nil
}

func arithEqual(args []scm.Object) (scm.Object, error) {
	l := args[0].(scm.Number)
	r := args[1].(scm.Number)
	return scm.Boolean{scm.EqNum(l, r)}, nil
}

func arithGreaterThan(args []scm.Object) (scm.Object, error) {
	l := args[0].(scm.Number)
	r := args[1].(scm.Number)
	if _, ok := l.(scm.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}
	if _, ok := r.(scm.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}

	return scm.Boolean{scm.GTNum(l, r)}, nil
}

func arithLessThan(args []scm.Object) (scm.Object, error) {
	l := args[0].(scm.Number)
	r := args[1].(scm.Number)
	if _, ok := l.(scm.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}
	if _, ok := r.(scm.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}

	return scm.Boolean{scm.LTNum(l, r)}, nil
}

func arithGreaterThanEuqal(args []scm.Object) (scm.Object, error) {
	l := args[0].(scm.Number)
	r := args[1].(scm.Number)
	if _, ok := l.(scm.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}
	if _, ok := r.(scm.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}

	return scm.Boolean{scm.GTENum(l, r)}, nil
}

func arithLessThanEqual(args []scm.Object) (scm.Object, error) {
	l := args[0].(scm.Number)
	r := args[1].(scm.Number)
	if _, ok := l.(scm.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}
	if _, ok := r.(scm.CompNum); ok {
		return nil, &Error{Message: "real number is required"}
	}

	return scm.Boolean{scm.LTENum(l, r)}, nil
}

func implIsNumber(expr scm.Object) bool {
	switch expr.(type) {
	case scm.IntNum, scm.RealNum, scm.RatNum, scm.CompNum:
		return true
	default:
		return false
	}
}

func arithIsNumber(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("requires 1, but got %d", len(args))}
	}
	return scm.Boolean{implIsNumber(args[0])}, nil
}

func arithNumberToString(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("requires 1, but got %d", len(args))}
	}
	if implIsNumber(args[0]) {
		return scm.NumberToString(args[0]), nil
	} else {
		return nil, &Error{Message: fmt.Sprintf("expected number, but got %s", scm.TypeString(args[0]))}
	}
}

func arithStringToNumber(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("requires 1, but got %d", len(args))}
	}
	if s, ok := args[0].(scm.String); ok {
		return scm.StringToNumber(string(s)), nil
	} else {
		return nil, &Error{Message: fmt.Sprintf("expected string, but got %s", scm.TypeString(args[0]))}
	}
}
