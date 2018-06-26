package vm

import (
	"log"

	"github.com/potix2/goscheme/ast"
)

func Eval(e ast.Expr, env *Env) (ast.Expr, error) {
	log.Printf("eval: %#v (env:%v)\n", e, env)
	if isVariable(e) {
		return e, nil
	}

	//eval(op operands) => (apply eval(op) operands)
	if a, ok := e.(ast.AppExpr); ok {
		log.Printf("eval: %#v (env:%v)\n", e, env)
		op, err := Eval(a.Operator, env)
		if err != nil {
			return nil, err
		}
		vals, err := evalValues(a.Operands, env)
		if err != nil {
			return nil, err
		}

		return Apply(op, vals)
	}

	if ide, ok := e.(ast.IdentExpr); ok {
		exp, err := env.Lookup(ide.Lit)
		if err != nil {
			return nil, err
		}
		return Eval(exp, env)
	}
	return e, nil
}

func Apply(op ast.Expr, args []ast.Expr) (ast.Expr, error) {
	log.Printf("call Apply: %#v, %#v\n", op, args)
	if p, ok := op.(ast.PrimitiveProcExpr); ok {
		return p.Proc(args)
	}
	return nil, &Error{Message: "Not implemented"}
}

func evalValues(args []ast.Expr, env *Env) ([]ast.Expr, error) {
	ret := make([]ast.Expr, 0, len(args))
	for _, arg := range args {
		v, err := Eval(arg, env)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}
	return ret, nil
}

func isVariable(e ast.Expr) bool {
	switch e.(type) {
	case ast.Uint10Expr, ast.PrimitiveProcExpr:
		return true
	default:
		return false
	}
}

//((lambda (x) (+ 1 x)) 2) => (+ 1 2) => 3
