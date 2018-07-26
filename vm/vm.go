package vm

import (
	"os"

	"github.com/potix2/goscheme/scm"
)

var CurrentVM VM

type VM struct {
	Input  scm.InputPort
	Output scm.OutputPort
	Error  scm.OutputPort
	Env    *scm.Env
}

func InitVM(env *scm.Env) {
	input := scm.InputPort{os.Stdin, false}
	output := scm.OutputPort{os.Stdout, false}
	err := scm.OutputPort{os.Stderr, false}

	CurrentVM = VM{input, output, err, env}
}
