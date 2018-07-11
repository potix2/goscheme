package vm

import (
	"github.com/potix2/goscheme/ast"
)

func makePrimitive(op string, proc ast.PrimitiveProc) ast.PrimitiveProcExpr {
	return ast.PrimitiveProcExpr{Operator: op, Proc: proc}
}

func SetupPrimitives(e *ast.Env) {

	//arithmetic operators
	e.Bind("+", makePrimitive("+", arithAdd))
	e.Bind("*", makePrimitive("*", arithMul))
	e.Bind("-", makePrimitive("-", arithSub))
	e.Bind("/", makePrimitive("/", arithDiv))
	e.Bind("=", makePrimitive("=", arithEqual))
	e.Bind("<", makePrimitive("<", arithLessThan))
	e.Bind(">", makePrimitive(">", arithGreaterThan))
	e.Bind("<=", makePrimitive("<=", arithLessThanEqual))
	e.Bind(">=", makePrimitive(">=", arithGreaterThanEuqal))
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

	//port
	e.Bind("current-input-port", makePrimitive("current-input-port", portCurrentInputPort))
	e.Bind("current-output-port", makePrimitive("current-output-port", portCurrentOutputPort))
	e.Bind("current-error-port", makePrimitive("current-error-port", portCurrentErrorPort))

	//string
	e.Bind("string?", makePrimitive("string?", strIsString))
	e.Bind("string-length", makePrimitive("string-length", strStringLength))
	e.Bind("string=?", makePrimitive("string=?", strStringEqual))
	e.Bind("string<?", makePrimitive("string<?", strStringLT))
	e.Bind("string<=?", makePrimitive("string<=?", strStringLTE))
	e.Bind("string>?", makePrimitive("string>?", strStringGT))
	e.Bind("string>=?", makePrimitive("string>=?", strStringGTE))
	e.Bind("substring", makePrimitive("substring", strSubstring))
	e.Bind("string-append", makePrimitive("string-append", strStringAppend))
}
