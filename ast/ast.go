package ast

import (
	"io"
	"strconv"
)

type Expr interface {
	expression()
	Print(io.Writer)
}

type PrimitiveProc func([]Expr) (Expr, error)

type (
	IdentExpr struct {
		Lit string
	}
	Uint10Expr struct {
		Lit int
	}
	AppExpr struct {
		Operator Expr
		Operands []Expr
	}
	PrimitiveProcExpr struct {
		Operator string
		Proc     PrimitiveProc
	}
)

func (x IdentExpr) expression() {}
func (x IdentExpr) Print(output io.Writer) {
	output.Write([]byte(x.Lit))
}

func (x Uint10Expr) expression() {}
func (x Uint10Expr) Print(output io.Writer) {
	output.Write([]byte(strconv.Itoa(x.Lit)))
}

func (x AppExpr) expression() {}
func (x AppExpr) Print(output io.Writer) {
	output.Write([]byte("("))
	x.Operator.Print(output)
	for _, op := range x.Operands {
		output.Write([]byte(" "))
		op.Print(output)
	}
	output.Write([]byte(")"))
}

func (x PrimitiveProcExpr) expression() {}
func (x PrimitiveProcExpr) Print(output io.Writer) {
	output.Write([]byte(x.Operator))
}
