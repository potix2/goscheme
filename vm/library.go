package vm

import (
	//"fmt"
	//"io"
	//"io/ioutil"
	//"path"
	//"strings"

	"github.com/potix2/goscheme/scm"
	//"github.com/potix2/goscheme/parser"
)

type Library struct {
	Name     string
	Symbols  []string
	TopLevel *scm.Env
}

var libraries map[string]Library
var libraryPaths []string

/*
// name is like `scheme.base` `srfi`
func Load(name string) error {
	filepath := strings.Replace(name, ".", "/", -1)
	sourcePath, err := lookup(filepath + ".scm")
	if err != nil {
		return err
	}

	source, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	code := string(source)

	env := NewEnv()
	InitVM(env)
	SetupPrimitives(env)
	exprs, err := parser.Read(code)
	expr, err := Eval(exprs[0], env)
	if err != nil {
		if lib, ok := expr.(Library); ok {
			libraries[name] = lib
			return nil
		} else {
			return
		}
}

func lookup(p string) (string, error) {
	for _, libPath := range libraryPaths {
		path.Join(libPath, p)
	}
}
*/

/*
func libImport(args []scm.Object) (scm.Object, error) {
}

func libExport(args []scm.Object) (scm.Object, error) {
}
*/
