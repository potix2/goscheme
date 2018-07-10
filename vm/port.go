package vm

import (
	"github.com/potix2/goscheme/ast"
)

func portCurrentInputPort(args []ast.Expr) (ast.Expr, error) {
	return CurrentVM.Input, nil
}

func portCurrentOutputPort(args []ast.Expr) (ast.Expr, error) {
	return CurrentVM.Output, nil
}

func portCurrentErrorPort(args []ast.Expr) (ast.Expr, error) {
	return CurrentVM.Error, nil
}
