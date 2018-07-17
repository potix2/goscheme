package vm

import (
	"fmt"

	"github.com/potix2/goscheme/ast"

	"github.com/sirupsen/logrus"
)

func evalProgram(exprs []ast.Expr, env *ast.Env) (ast.Expr, error) {
	//TODO: envが必要なもの(program)とそれ以外を分けてapplyの引数からenvを削除する
	//logrus.Debugf("eval: %#v (env:%v)\n", e, env)
	if op, ok := exprs[0].(ast.IdentExpr); ok {
		switch op.Lit {
		case "lambda":
			return ast.LambdaExpr{exprs[1], exprs[2], env}, nil
		case "define":
			//(define <variable> <expression>)
			if id, ok := exprs[1].(ast.IdentExpr); ok {
				if val, err := Eval(exprs[2], env); err == nil {
					env.Bind(id.Lit, val)
					return id, nil
				}
			}

			//(define (<variable> . <formal>) <body>)
			//    => (define <variable> (lambda <formal> <body>))
			if formals, ok := exprs[1].(ast.AppExpr); ok {
				if variable, ok := formals.Exprs[0].(ast.IdentExpr); ok {
					env.Bind(variable.Lit, ast.LambdaExpr{ast.AppExpr{formals.Exprs[1:]}, exprs[2], env})
					return variable, nil
				}
			}
		case "set!":
			if id, ok := exprs[1].(ast.IdentExpr); ok {
				if val, err := Eval(exprs[2], env); err == nil {
					env.Bind(id.Lit, val)
					return id, nil
				}
			}
		case "quote":
			if len(exprs) != 2 {
				return nil, &Error{Message: fmt.Sprintf("expected 1 args, but got %d\n", len(exprs)-1)}
			}
			return exprs[1], nil
			/*
				case "list":
					vals, err := evalValues(exprs[1:], env)
					if err != nil {
						return nil, err
					}
					return ast.AppExpr{vals}, nil
			*/
		case "if":
			if len(exprs) != 4 && len(exprs) != 3 {
				return nil, &Error{Message: fmt.Sprintf("expected 2 or 3 args, but got %d\n", len(exprs)-1)}
			}

			test, err := Eval(exprs[1], env)
			if err != nil {
				return nil, err
			}
			if tv, ok := test.(ast.BooleanExpr); ok && !tv.Lit {
				//alternate
				if len(exprs) == 4 {
					result, err := Eval(exprs[3], env)
					if err != nil {
						return nil, err
					}
					return result, nil
				} else {
					return ast.Undefined{}, nil
				}
			} else {
				//consequent
				result, err := Eval(exprs[2], env)
				if err != nil {
					return nil, err
				}
				return result, nil
			}
		case "begin":
			if len(exprs) <= 1 {
				return nil, &Error{Message: fmt.Sprintf("required at least 1, but got %d", len(exprs)-1)}
			}

			var ret ast.Expr
			var err error
			for _, e := range exprs[1:] {
				ret, err = Eval(e, env)
				if err != nil {
					return nil, err
				}
			}
			return ret, nil
		}
	}

	op, err := Eval(exprs[0], env)
	if err != nil {
		return nil, err
	}

	vals, err := evalValues(exprs[1:], env)
	if err != nil {
		return nil, err
	}

	return apply(op, vals)
}

func apply(op ast.Expr, vals []ast.Expr) (ast.Expr, error) {
	if p, ok := op.(ast.PrimitiveProcExpr); ok {
		return p.Proc(vals)
	}

	if l, ok := op.(ast.LambdaExpr); ok {
		if vars, ok := l.Args.(ast.AppExpr); ok {
			if len(vars.Exprs) != len(vals) {
				return nil, &Error{Message: fmt.Sprintf("expected %d args, but got %d\n", len(vars.Exprs), len(vals))}
			}

			newEnv := Extend(l.Closure, map[string]ast.Expr{})
			for i, a := range vars.Exprs {
				if id, ok := a.(ast.IdentExpr); ok {
					newEnv.Bind(id.Lit, vals[i])
				}
			}
			return Eval(l.Body, newEnv)
		}

		if argList, ok := l.Args.(ast.IdentExpr); ok {
			newEnv := Extend(l.Closure, map[string]ast.Expr{})
			newEnv.Bind(argList.Lit, ast.QuoteExpr{recMakeListFromSlice(vals)})
			return Eval(l.Body, newEnv)
		}
	}

	return nil, &Error{Message: fmt.Sprintf("got unapplicable expression: op=%#v, vals=%#v\n", op, vals)}
}

func Eval(e ast.Expr, env *ast.Env) (ast.Expr, error) {
	logrus.Debugf("eval: %#v (env:%v)\n", e, env)
	if isVariable(e) {
		return e, nil
	}

	if isList(e) {
		return evalProgram(listToSlice(e), env)
	}

	//eval(op operands) => (apply eval(op) operands)
	if a, ok := e.(ast.AppExpr); ok {
		return evalProgram(a.Exprs, env)
	}

	if qe, ok := e.(ast.QuoteExpr); ok {
		return qe.Datum, nil
	}

	if ide, ok := e.(ast.IdentExpr); ok {
		exp, err := Lookup(env, ide.Lit)
		if err != nil {
			return nil, err
		}
		return Eval(exp, env)
	}
	return e, nil
}

func evalValues(args []ast.Expr, env *ast.Env) ([]ast.Expr, error) {
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
	case ast.IntNum, ast.RealNum, ast.RatNum, ast.CompNum, ast.BooleanExpr, ast.PrimitiveProcExpr, ast.InputPort, ast.OutputPort, ast.StringExpr:
		return true
	default:
		return false
	}
}

func isList(e ast.Expr) bool {
	if p, ok := e.(ast.PairExpr); ok {
		return p.IsList()
	} else {
		return false
	}
}
