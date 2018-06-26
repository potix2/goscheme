package vm

import (
	"io"

	"github.com/potix2/goscheme/ast"
)

type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func Print(output io.Writer, expr ast.Expr) {
	expr.Print(output)
}
