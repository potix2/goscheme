package vm

import (
	"fmt"

	"github.com/potix2/goscheme/ast"

	"github.com/sirupsen/logrus"
)

func Eval(e ast.Expr, env *ast.Env) (ast.Expr, error) {
	logrus.Debugf("eval: %#v (env:%v)\n", e, env)
	if isVariable(e) {
		return e, nil
	}

	//eval(op operands) => (apply eval(op) operands)
	if a, ok := e.(ast.AppExpr); ok {
		logrus.Debugf("eval: %#v (env:%v)\n", e, env)

		if op, ok := a.Exprs[0].(ast.IdentExpr); ok {
			switch op.Lit {
			case "lambda":
				return ast.LambdaExpr{a.Exprs[1], a.Exprs[2], env}, nil
				//TODO: implement define
			case "define":
				//(define <variable> <expression>)
				if id, ok := a.Exprs[1].(ast.IdentExpr); ok {
					if val, err := Eval(a.Exprs[2], env); err == nil {
						env.Bind(id.Lit, val)
						return id, nil
					}
				}

				//(define (<variable> . <formal>) <body>)
				//    => (define <variable> (lambda <formal> <body>))
				if formals, ok := a.Exprs[1].(ast.AppExpr); ok {
					if variable, ok := formals.Exprs[0].(ast.IdentExpr); ok {
						env.Bind(variable.Lit, ast.LambdaExpr{ast.AppExpr{formals.Exprs[1:]}, a.Exprs[2], env})
						return variable, nil
					}
				}
			case "set!":
				if id, ok := a.Exprs[1].(ast.IdentExpr); ok {
					if val, err := Eval(a.Exprs[2], env); err == nil {
						env.Bind(id.Lit, val)
						return id, nil
					}
				}
			case "quote":
				if len(a.Exprs) != 2 {
					return nil, &Error{Message: fmt.Sprintf("expected 1 args, but got %d\n", len(a.Exprs)-1)}
				}
				return a.Exprs[1], nil
			case "list":
				vals, err := evalValues(a.Exprs[1:], env)
				if err != nil {
					return nil, err
				}
				return ast.AppExpr{vals}, nil
			case "if":
				if len(a.Exprs) != 4 && len(a.Exprs) != 3 {
					return nil, &Error{Message: fmt.Sprintf("expected 2 or 3 args, but got %d\n", len(a.Exprs)-1)}
				}

				test, err := Eval(a.Exprs[1], env)
				if err != nil {
					return nil, err
				}
				if tv, ok := test.(ast.BooleanExpr); ok && !tv.Lit {
					//alternate
					if len(a.Exprs) == 4 {
						result, err := Eval(a.Exprs[3], env)
						if err != nil {
							return nil, err
						}
						return result, nil
					} else {
						return ast.Undefined{}, nil
					}
				} else {
					//consequent
					result, err := Eval(a.Exprs[2], env)
					if err != nil {
						return nil, err
					}
					return result, nil
				}
			}
		}

		op, err := Eval(a.Exprs[0], env)
		if err != nil {
			return nil, err
		}
		vals, err := evalValues(a.Exprs[1:], env)
		if err != nil {
			return nil, err
		}

		return Apply(op, vals)
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

func Apply(op ast.Expr, args []ast.Expr) (ast.Expr, error) {
	logrus.Debugf("call Apply: %#v, %#v\n", op, args)
	if p, ok := op.(ast.PrimitiveProcExpr); ok {
		return p.Proc(args)
	}

	if l, ok := op.(ast.LambdaExpr); ok {
		if vars, ok := l.Args.(ast.AppExpr); ok {
			if len(vars.Exprs) != len(args) {
				return nil, &Error{Message: fmt.Sprintf("expected %d args, but got %d\n", len(vars.Exprs), len(args))}
			}

			newEnv := Extend(l.Closure, map[string]ast.Expr{})
			for i, a := range vars.Exprs {
				if id, ok := a.(ast.IdentExpr); ok {
					newEnv.Bind(id.Lit, args[i])
				}
			}
			return Eval(l.Body, newEnv)
		}
	}

	return nil, &Error{Message: fmt.Sprintf("got unapplicable expression: %#v\n", op)}
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
	case ast.Uint10Expr, ast.BooleanExpr, ast.PrimitiveProcExpr:
		return true
	default:
		return false
	}
}
