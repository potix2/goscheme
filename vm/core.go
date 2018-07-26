package vm

import (
	"bytes"
	"fmt"
	"io"

	"github.com/potix2/goscheme/scm"

	"github.com/sirupsen/logrus"
)

//((lambda (x)
//  (lambda (y) (+ x y))) 2)
// =>
//(lambda (y) (+ x y)) {env: x = 2}
func evalProgram(exprs []scm.Object, env *scm.Env) (scm.Object, error) {
	logrus.WithFields(logrus.Fields{
		"expr": dumpExpr(scm.AppExpr{exprs}),
		"env":  dumpEnv(env),
	}).Debug("evalProgram")

	if op, ok := exprs[0].(scm.IdentExpr); ok {
		switch op.Lit {
		case "lambda":
			return scm.LambdaExpr{exprs[1], exprs[2:], env}, nil
		case "define":
			//(define <variable> <expression>)
			if id, ok := exprs[1].(scm.IdentExpr); ok {
				if val, err := Eval(exprs[2], env); err == nil {
					env.Bind(id.Lit, val)
					return id, nil
				}
			}

			//(define (<variable> . <formal>) <body>)
			//    => (define <variable> (lambda <formal> <body>))
			if formals, ok := exprs[1].(scm.AppExpr); ok {
				if variable, ok := formals.Objs[0].(scm.IdentExpr); ok {
					env.Bind(variable.Lit, scm.LambdaExpr{scm.AppExpr{formals.Objs[1:]}, exprs[2:], env})
					return variable, nil
				}
			}
		case "set!":
			if id, ok := exprs[1].(scm.IdentExpr); ok {
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
		case "if":
			if len(exprs) != 4 && len(exprs) != 3 {
				return nil, &Error{Message: fmt.Sprintf("expected 2 or 3 args, but got %d\n", len(exprs)-1)}
			}

			test, err := Eval(exprs[1], env)
			if err != nil {
				return nil, err
			}
			if tv, ok := test.(scm.BooleanExpr); ok && !tv.Lit {
				//alternate
				if len(exprs) == 4 {
					result, err := Eval(exprs[3], env)
					if err != nil {
						return nil, err
					}
					return result, nil
				} else {
					return scm.Undefined{}, nil
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
				return nil, &Error{Message: fmt.Sprintf("required at lescm 1, but got %d", len(exprs)-1)}
			}

			var ret scm.Object
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

func apply(op scm.Object, vals []scm.Object) (scm.Object, error) {
	if p, ok := op.(scm.PrimitiveProcExpr); ok {
		return p.Proc(vals)
	}

	if l, ok := op.(scm.LambdaExpr); ok {
		if vars, ok := l.Args.(scm.AppExpr); ok {
			if len(vars.Objs) != len(vals) {
				return nil, &Error{Message: fmt.Sprintf("expected %d args, but got %d\n", len(vars.Objs), len(vals))}
			}

			newEnv := Extend(l.Closure, map[string]scm.Object{})
			for i, a := range vars.Objs {
				if id, ok := a.(scm.IdentExpr); ok {
					newEnv.Bind(id.Lit, vals[i])
				}
			}

			var ret scm.Object
			var err error
			for _, e := range l.Body {
				ret, err = Eval(e, newEnv)
				if err != nil {
					return nil, err
				}
			}
			return ret, nil
		}

		if argList, ok := l.Args.(scm.IdentExpr); ok {
			newEnv := Extend(l.Closure, map[string]scm.Object{})
			newEnv.Bind(argList.Lit, recMakeListFromSlice(vals))

			var ret scm.Object
			var err error
			for _, e := range l.Body {
				ret, err = Eval(e, newEnv)
				if err != nil {
					return nil, err
				}
			}
			return ret, nil
		}
	}

	return nil, &Error{Message: fmt.Sprintf("got unapplicable expression: op=%#v, vals=%#v\n", op, vals)}
}

func Eval(e scm.Object, env *scm.Env) (scm.Object, error) {
	logrus.WithFields(logrus.Fields{
		"expr": dumpExpr(e),
		"env":  dumpEnv(env),
	}).Debug("eval")

	if isVariable(e) {
		return e, nil
	}

	//eval(op operands) => (apply eval(op) operands)
	if a, ok := e.(scm.AppExpr); ok {
		return evalProgram(a.Objs, env)
	}

	if qe, ok := e.(scm.QuoteExpr); ok {
		return qe.Datum, nil
	}

	if ide, ok := e.(scm.IdentExpr); ok {
		exp, err := Lookup(env, ide.Lit)
		if err != nil {
			return nil, err
		}
		return Eval(exp, env)
	}
	return e, nil
}

func evalValues(args []scm.Object, env *scm.Env) ([]scm.Object, error) {
	ret := make([]scm.Object, 0, len(args))
	for _, arg := range args {
		v, err := Eval(arg, env)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}
	return ret, nil
}

func isVariable(e scm.Object) bool {
	switch e.(type) {
	case scm.IntNum, scm.RealNum, scm.RatNum, scm.CompNum, scm.BooleanExpr, scm.PrimitiveProcExpr, scm.InputPort, scm.OutputPort, scm.StringExpr, scm.PairExpr:
		return true
	default:
		return false
	}
}

func isList(e scm.Object) bool {
	if p, ok := e.(scm.PairExpr); ok {
		return p.IsList()
	} else {
		return false
	}
}

func dumpExpr(e scm.Object) string {
	var buf bytes.Buffer
	e.Print(&buf)
	return buf.String()
}

func dumpEnvImpl(w io.Writer, env *scm.Env) {
	w.Write([]byte("{"))
	for k, _ := range env.Values {
		fmt.Fprintf(w, "%s, ", k)
	}
	if env.Parent != nil {
		fmt.Fprintf(w, "parent(%p):", env.Parent)
		dumpEnvImpl(w, env.Parent)
	}
	w.Write([]byte("}"))
}

func dumpEnv(env *scm.Env) string {
	var buf bytes.Buffer
	dumpEnvImpl(&buf, env)
	return buf.String()
}
