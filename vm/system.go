package vm

import (
	"fmt"
	"io/ioutil"

	"github.com/potix2/goscheme/parser"
	"github.com/potix2/goscheme/scm"
)

func implLoad(filename string, env *scm.Env) error {
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
func sysLoad(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("required 1, but got %d", len(args))}
	}

	var s scm.StringExpr
	var ok bool
	if s, ok = args[0].(scm.StringExpr); !ok {
		return nil, &Error{Message: fmt.Sprintf("expected string, but got %s", scm.TypeString(args[0]))}
	}

	//TODO: pass from optional args
	env := getInteractionEnvironment()
	err := implLoad(string(s), env)
	if err != nil {
		return nil, err
	}
	return scm.Undefined{}, nil
}
