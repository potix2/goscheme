package vm

import (
	"bytes"
	"fmt"

	"github.com/potix2/goscheme/scm"
)

func implIsString(expr scm.Object) bool {
	_, ok := expr.(scm.String)
	return ok
}

func strIsString(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("requires 1, but got %d", len(args))}
	}
	return scm.Boolean{implIsString(args[0])}, nil
}

func strStringLength(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("requires 1, but got %d", len(args))}
	}
	if s, ok := args[0].(scm.String); ok {
		return scm.IntNum(len(s)), nil
	} else {
		return nil, &Error{Message: fmt.Sprintf("type mismatch: expected string, but got %s", scm.TypeString(args[0]))}
	}
}

type strComparator func(scm.String, scm.String) bool

func compForAll(args []scm.Object, comp strComparator) (scm.Object, error) {
	if len(args) == 0 {
		return nil, &Error{Message: "wrong number of arguments"}
	}
	a := args[0].(scm.String)
	for _, b := range args[1:] {
		if bs, ok := b.(scm.String); !ok || !comp(a, bs) {
			return scm.Boolean{false}, nil
		}
	}
	return scm.Boolean{true}, nil
}

func strStringEqual(args []scm.Object) (scm.Object, error) {
	return compForAll(args, func(a, b scm.String) bool { return a == b })
}

func strStringLT(args []scm.Object) (scm.Object, error) {
	return compForAll(args, func(a, b scm.String) bool { return a < b })
}

func strStringGT(args []scm.Object) (scm.Object, error) {
	return compForAll(args, func(a, b scm.String) bool { return a > b })
}

func strStringLTE(args []scm.Object) (scm.Object, error) {
	return compForAll(args, func(a, b scm.String) bool { return a <= b })
}

func strStringGTE(args []scm.Object) (scm.Object, error) {
	return compForAll(args, func(a, b scm.String) bool { return a >= b })
}

func strSubstring(args []scm.Object) (scm.Object, error) {
	if len(args) != 3 {
		return nil, &Error{Message: fmt.Sprintf("required 3, but got %d", len(args))}
	}
	var s scm.String
	var ok bool
	var start, end scm.IntNum
	if s, ok = args[0].(scm.String); !ok {
		return nil, &Error{Message: fmt.Sprintf("string required, but got %s", scm.TypeString(args[0]))}
	}
	if start, ok = args[1].(scm.IntNum); !ok {
		return nil, &Error{Message: fmt.Sprintf("integer required, but got %s", scm.TypeString(args[1]))}
	}
	if end, ok = args[2].(scm.IntNum); !ok {
		return nil, &Error{Message: fmt.Sprintf("integer required, but got %s", scm.TypeString(args[2]))}
	}
	if start < 0 || len(s) <= int(end) {
		return nil, &Error{Message: fmt.Sprintf("out of range: %d %d", start, end)}
	}
	return scm.String(s[start:end]), nil
}

func strStringAppend(args []scm.Object) (scm.Object, error) {
	if len(args) == 0 {
		return nil, &Error{Message: fmt.Sprintf("required at lescm 1, but got %d", len(args))}
	}
	var buffer bytes.Buffer
	for _, e := range args {
		if s, ok := e.(scm.String); !ok {
			return nil, &Error{Message: fmt.Sprintf("expected string, but got %s", scm.TypeString(e))}
		} else {
			buffer.WriteString(string(s))
		}
	}
	return scm.String(buffer.String()), nil
}
