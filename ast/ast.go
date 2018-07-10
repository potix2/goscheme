package ast

import (
	"io"
)

type Env struct {
	Values map[string]Expr
	Parent *Env
}

type Expr interface {
	Print(io.Writer)
}

type PrimitiveProc func([]Expr) (Expr, error)

type (
	IdentExpr struct {
		Lit string
	}
	BooleanExpr struct {
		Lit bool
	}
	QuoteExpr struct {
		Datum Expr
	}
	PairExpr struct {
		Car Expr
		Cdr Expr
	}
	AppExpr struct {
		Exprs []Expr
	}
	//(lambda (x) (+ 1 x))
	LambdaExpr struct {
		Args    Expr
		Body    Expr
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

func MakeListFromSlice(exprs []Expr) PairExpr {
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
	if len(x.Exprs) > 0 {
		x.Exprs[0].Print(output)
		for _, op := range x.Exprs[1:] {
			output.Write([]byte(" "))
			op.Print(output)
		}
	}
	output.Write([]byte(")"))
}

func (x LambdaExpr) Print(output io.Writer) {
	output.Write([]byte("(lambda "))
	x.Args.Print(output)
	output.Write([]byte(" "))
	x.Body.Print(output)
	output.Write([]byte(")"))
}

func (x PrimitiveProcExpr) Print(output io.Writer) {
	output.Write([]byte(x.Operator))
}

func (x Undefined) Print(output io.Writer) {
	output.Write([]byte("#undef"))
}

func (x InputPort) Print(output io.Writer) {
	output.Write([]byte("#<iport>"))
}

func (x OutputPort) Print(output io.Writer) {
	output.Write([]byte("#<oport>"))
}

func (e *Env) Bind(name string, value Expr) {
	e.Values[name] = value
}
