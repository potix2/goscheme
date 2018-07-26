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

type PrimitiveProc func([]Object) (Object, error)

type StringExpr string
type (
	IdentExpr struct {
		Lit string
	}
	BooleanExpr struct {
		Lit bool
	}
	QuoteExpr struct {
		Datum Object
	}
	PairExpr struct {
		Car Object
		Cdr Object
	}
	AppExpr struct {
		Objs []Object
	}
	//(lambda (x) (+ 1 x))
	LambdaExpr struct {
		Args    Object
		Body    []Object
		Closure *Env
	}
	//(let
	//  (x (lambda (a) (+ a 2)))
	//  (x 2))

	PrimitiveProcExpr struct {
		Operator string
		Proc     PrimitiveProc
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

func MakeListFromSlice(exprs []Object) PairExpr {
	if len(exprs) == 0 {
		return PairExpr{}
	} else {
		return PairExpr{Car: exprs[0], Cdr: MakeListFromSlice(exprs[1:])}
	}
}

func (x IdentExpr) Print(output io.Writer) {
	output.Write([]byte(x.Lit))
}

func (x BooleanExpr) Print(output io.Writer) {
	if x.Lit {
		output.Write([]byte("#t"))
	} else {
		output.Write([]byte("#f"))
	}
}

func (x PairExpr) Print(output io.Writer) {
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

func (x PairExpr) IsEmptyList() bool {
	return x.Car == nil && x.Cdr == nil
}

func (x PairExpr) IsList() bool {
	if x.IsEmptyList() {
		return true
	}

	if child, ok := x.Cdr.(PairExpr); ok {
		return child.IsList()
	}
	return false
}

func (x QuoteExpr) Print(output io.Writer) {
	output.Write([]byte("'"))
	x.Datum.Print(output)
}

func (x AppExpr) Print(output io.Writer) {
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

func (x LambdaExpr) Print(output io.Writer) {
	output.Write([]byte("#<closure (#f"))
	x.Args.Print(output)
	output.Write([]byte(")>"))
}

func (x PrimitiveProcExpr) Print(output io.Writer) {
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

func (x StringExpr) Print(output io.Writer) {
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
	case BooleanExpr:
		return "boolean"
	case IntNum:
		return "number"
	case RatNum:
		return "number"
	case RealNum:
		return "number"
	case IdentExpr:
		return "symbol"
	case StringExpr:
		return "string"
	case InputPort:
		return "port"
	case OutputPort:
		return "port"
	case PairExpr:
		return "pair"
	case LambdaExpr:
		return "procedure"
	case PrimitiveProcExpr:
		return "procedure"
	case QuoteExpr:
		return "quote"
	case Env:
		return "env"
	default:
		return "unknown"
	}
}
