package vm

import (
	"github.com/potix2/goscheme/ast"
)

func makePrimitive(op string, proc ast.PrimitiveProc) ast.PrimitiveProcExpr {
	return ast.PrimitiveProcExpr{Operator: op, Proc: proc}
}

func bindPrimitive(e *ast.Env, op string, proc ast.PrimitiveProc) {
	e.Bind(op, makePrimitive(op, proc))
}

func SetupPrimitives(e *ast.Env) {

	//arithmetic operators
	bindPrimitive(e, "+", arithAdd)
	bindPrimitive(e, "*", arithMul)
	bindPrimitive(e, "-", arithSub)
	bindPrimitive(e, "/", arithDiv)
	bindPrimitive(e, "=", arithEqual)
	bindPrimitive(e, "<", arithLessThan)
	bindPrimitive(e, ">", arithGreaterThan)
	bindPrimitive(e, "<=", arithLessThanEqual)
	bindPrimitive(e, ">=", arithGreaterThanEuqal)
	bindPrimitive(e, "number?", arithIsNumber)
	bindPrimitive(e, "number->string", arithNumberToString)
	bindPrimitive(e, "string->number", arithStringToNumber)

	//boolean operators
	bindPrimitive(e, "not", boolNot)
	bindPrimitive(e, "boolean?", boolIsBoolean)
	bindPrimitive(e, "procedure?", boolIsProcedure)
	//bindPrimitive(e, "char?", boolIsChar))
	//bindPrimitive(e, "vector?", "vector?", boolIsVector))

	//list operatators
	bindPrimitive(e, "cons", listCons)
	bindPrimitive(e, "car", listCar)
	bindPrimitive(e, "cdr", listCdr)
	bindPrimitive(e, "list", listList)
	bindPrimitive(e, "pair?", listIsPair)
	bindPrimitive(e, "list?", listIsList)
	bindPrimitive(e, "null?", listIsNull)

	//port
	bindPrimitive(e, "current-input-port", portCurrentInputPort)
	bindPrimitive(e, "current-output-port", portCurrentOutputPort)
	bindPrimitive(e, "current-error-port", portCurrentErrorPort)
	bindPrimitive(e, "write-string", portWriteString)
	bindPrimitive(e, "display", portDisplay)
	bindPrimitive(e, "newline", portNewline)
	bindPrimitive(e, "input-port?", portIsInputPort)
	bindPrimitive(e, "output-port?", portIsOutputPort)
	bindPrimitive(e, "textual-port?", portIsTextualPort)
	bindPrimitive(e, "binary-port?", portIsBinaryPort)
	bindPrimitive(e, "port?", portIsPort)

	//string
	bindPrimitive(e, "string?", strIsString)
	bindPrimitive(e, "string-length", strStringLength)
	bindPrimitive(e, "string=?", strStringEqual)
	bindPrimitive(e, "string<?", strStringLT)
	bindPrimitive(e, "string<=?", strStringLTE)
	bindPrimitive(e, "string>?", strStringGT)
	bindPrimitive(e, "string>=?", strStringGTE)
	bindPrimitive(e, "substring", strSubstring)
	bindPrimitive(e, "string-append", strStringAppend)

	//system
	bindPrimitive(e, "load", sysLoad)

	//env
	bindPrimitive(e, "eval", envEval)
	bindPrimitive(e, "apply", envApply)
	bindPrimitive(e, "interaction-environment", envInteractionEnvironment)
}
