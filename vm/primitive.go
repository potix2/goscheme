package vm

import (
	"github.com/potix2/goscheme/ast"
)

func makePrimitive(op string, proc ast.PrimitiveProc) ast.PrimitiveProcExpr {
	return ast.PrimitiveProcExpr{Operator: op, Proc: proc}
}

func SetupPrimitives(e *Env) error {

	//arithmetic operators
	e.SetVariable("+", makePrimitive("+", arithAdd))
	e.SetVariable("*", makePrimitive("*", arithMul))
	e.SetVariable("-", makePrimitive("-", arithSub))
	e.SetVariable("/", makePrimitive("/", arithDiv))

	//boolean operators
	return nil
}
