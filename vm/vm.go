package vm

import (
	"os"

	"github.com/potix2/goscheme/ast"
)

var CurrentVM VM

type VM struct {
	Input  ast.InputPort
	Output ast.OutputPort
	Error  ast.OutputPort
	Env    *ast.Env
}

func InitVM(env *ast.Env) {
	input := ast.InputPort{os.Stdin, false}
	output := ast.OutputPort{os.Stdout, false}
	err := ast.OutputPort{os.Stderr, false}

	CurrentVM = VM{input, output, err, env}
}
