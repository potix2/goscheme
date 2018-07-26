package vm

import (
	"io"

	"github.com/potix2/goscheme/scm"
)

type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func Print(output io.Writer, expr scm.Expr) {
	expr.Print(output)
}
