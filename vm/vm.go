package vm

import (
	"bufio"
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
	input := ast.InputPort{bufio.NewReader(os.Stdin), false}
	output := ast.OutputPort{bufio.NewWriter(os.Stdout), false}
	err := ast.OutputPort{bufio.NewWriter(os.Stderr), false}

	CurrentVM = VM{input, output, err, env}
}
