package vm

import (
	"fmt"
	"io/ioutil"

	"github.com/potix2/goscheme/ast"
	"github.com/potix2/goscheme/parser"
)

func implLoad(filename string, env *ast.Env) error {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	code := string(source)
	exprs, err := parser.Read(code)
	if err != nil {
		return err
	}
	for _, expr := range exprs {
		_, err := Eval(expr, env)
		if err != nil {
			return err
		}
	}
	return nil
}

//(load filename)
//(load filename environment-specifier)
func sysLoad(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("required 1, but got %d", len(args))}
	}

	var s ast.StringExpr
	var ok bool
	if s, ok = args[0].(ast.StringExpr); !ok {
		return nil, &Error{Message: fmt.Sprintf("expected string, but got %s", ast.TypeString(args[0]))}
	}

	//TODO: pass from optional args
	env := getInteractionEnvironment()
	err := implLoad(string(s), env)
	if err != nil {
		return nil, err
	}
	return ast.Undefined{}, nil
}
