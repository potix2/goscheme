package vm

import (
	"github.com/potix2/goscheme/ast"
)

func makePrimitive(op string, proc ast.PrimitiveProc) ast.PrimitiveProcExpr {
	return ast.PrimitiveProcExpr{Operator: op, Proc: proc}
}

func SetupPrimitives(e *ast.Env) error {

	//arithmetic operators
	e.Bind("+", makePrimitive("+", arithAdd))
	e.Bind("*", makePrimitive("*", arithMul))
	e.Bind("-", makePrimitive("-", arithSub))
	e.Bind("/", makePrimitive("/", arithDiv))
	e.Bind(">", makePrimitive(">", arithGreaterThan))
	e.Bind("<", makePrimitive("<", arithLessThan))
	e.Bind("number?", makePrimitive("number?", arithIsNumber))

	//boolean operators
	e.Bind("not", makePrimitive("not", boolNot))
	e.Bind("boolean?", makePrimitive("boolean?", boolIsBoolean))
	e.Bind("procedure?", makePrimitive("procedure?", boolIsProcedure))
	//e.Bind("char?", makePrimitive("char?", boolIsChar))
	//e.Bind("string?", makePrimitive("string?", boolIsString))
	//e.Bind("vector?", makePrimitive("vector?", boolIsVector))

	//list operatators
	e.Bind("cons", makePrimitive("cons", listCons))
	e.Bind("car", makePrimitive("car", listCar))
	e.Bind("cdr", makePrimitive("cdr", listCdr))
	e.Bind("pair?", makePrimitive("pair?", listIsPair))
	e.Bind("list?", makePrimitive("list?", listIsList))
	e.Bind("null?", makePrimitive("null?", listIsNull))
	return nil
}
