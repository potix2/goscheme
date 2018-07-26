package scm

import (
	"io"
)

type Env struct {
	Values map[string]Object
	Parent *Env
}

type Object interface {
	Print(io.Writer)
}

type PrimitiveFunc func([]Object) (Object, error)

type String string
type (
	Symbol struct {
		Lit string
	}
	Boolean struct {
		Lit bool
	}
	Quote struct {
		Datum Object
	}
	Pair struct {
		Car Object
		Cdr Object
	}
	Subp struct {
		Objs []Object
	}
	//(lambda (x) (+ 1 x))
	Lambda struct {
		Args    Object
		Body    []Object
		Closure *Env
	}
	//(let
	//  (x (lambda (a) (+ a 2)))
	//  (x 2))

	PrimitiveProc struct {
		Operator string
		Proc     PrimitiveFunc
	}
	InputPort struct {
		Reader io.Reader
		Binary bool
	}
	OutputPort struct {
		Writer io.Writer
		Binary bool
	}

	Undefined struct {
	}
)

func MakeListFromSlice(exprs []Object) Pair {
	if len(exprs) == 0 {
		return Pair{}
	} else {
		return Pair{Car: exprs[0], Cdr: MakeListFromSlice(exprs[1:])}
	}
}

func (x Symbol) Print(output io.Writer) {
	output.Write([]byte(x.Lit))
}

func (x Boolean) Print(output io.Writer) {
	if x.Lit {
		output.Write([]byte("#t"))
	} else {
		output.Write([]byte("#f"))
	}
}

func (x Pair) Print(output io.Writer) {
	output.Write([]byte("("))
	if x.Car != nil {
		x.Car.Print(output)
		if x.Cdr != nil {
			output.Write([]byte(" . "))
			x.Cdr.Print(output)
		}
	}
	output.Write([]byte(")"))
}

func (x Pair) IsEmptyList() bool {
	return x.Car == nil && x.Cdr == nil
}

func (x Pair) IsList() bool {
	if x.IsEmptyList() {
		return true
	}

	if child, ok := x.Cdr.(Pair); ok {
		return child.IsList()
	}
	return false
}

func (x Quote) Print(output io.Writer) {
	output.Write([]byte("'"))
	x.Datum.Print(output)
}

func (x Subp) Print(output io.Writer) {
	output.Write([]byte("("))
	if len(x.Objs) > 0 {
		x.Objs[0].Print(output)
		for _, op := range x.Objs[1:] {
			output.Write([]byte(" "))
			op.Print(output)
		}
	}
	output.Write([]byte(")"))
}

func (x Lambda) Print(output io.Writer) {
	output.Write([]byte("#<closure (#f"))
	x.Args.Print(output)
	output.Write([]byte(")>"))
}

func (x PrimitiveProc) Print(output io.Writer) {
	output.Write([]byte(x.Operator))
}

func (x Undefined) Print(output io.Writer) {
	output.Write([]byte("#<undef>"))
}

func (x InputPort) Print(output io.Writer) {
	output.Write([]byte("#<iport>"))
}

func (x OutputPort) Print(output io.Writer) {
	output.Write([]byte("#<oport>"))
}

func (x String) Print(output io.Writer) {
	output.Write([]byte(x))
}

func (e Env) Print(output io.Writer) {
	output.Write([]byte("#<env>"))
}

func (e *Env) Bind(name string, value Object) {
	e.Values[name] = value
}

func TypeString(expr Object) string {
	switch expr.(type) {
	case Boolean:
		return "boolean"
	case IntNum:
		return "number"
	case RatNum:
		return "number"
	case RealNum:
		return "number"
	case Symbol:
		return "symbol"
	case String:
		return "string"
	case InputPort:
		return "port"
	case OutputPort:
		return "port"
	case Pair:
		return "pair"
	case Lambda:
		return "procedure"
	case PrimitiveProc:
		return "procedure"
	case Quote:
		return "quote"
	case Env:
		return "env"
	default:
		return "unknown"
	}
}
