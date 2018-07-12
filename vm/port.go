package vm

import (
	"fmt"

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

func implWriteString(p ast.OutputPort, s ast.StringExpr) (ast.Expr, error) {
	n, err := p.Writer.Write([]byte(s))
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("failed to write: %v", err)}
	}
	if n != len(string(s)) {
		return nil, &Error{Message: fmt.Sprintf("actual written size: %d, expected %d", n, len(string(s)))}
	}

	return ast.Undefined{}, nil
}

func portWriteString(args []ast.Expr) (ast.Expr, error) {
	if len(args) == 0 {
		return nil, &Error{Message: fmt.Sprintf("required at least 1, but got %d", len(args))}
	}
	s, ok := args[0].(ast.StringExpr)
	if !ok {
		return nil, &Error{Message: fmt.Sprintf("expected string, but got %s", ast.TypeString(args[0]))}
	}

	port := CurrentVM.Output
	if len(args) > 1 {
		port, ok = args[1].(ast.OutputPort)
		if !ok {
			return nil, &Error{Message: fmt.Sprintf("expected port, but got %s", ast.TypeString(args[1]))}
		}
	}

	if len(args) == 3 {
		start, ok := args[2].(ast.IntNum)
		if !ok {
			return nil, &Error{Message: fmt.Sprintf("expected int, but got %s", ast.TypeString(args[2]))}
		}
		s = s[start:]
	} else if len(args) == 4 {
		start, ok := args[2].(ast.IntNum)
		if !ok {
			return nil, &Error{Message: fmt.Sprintf("expected int, but got %s", ast.TypeString(args[2]))}
		}

		end, ok := args[3].(ast.IntNum)
		if !ok {
			return nil, &Error{Message: fmt.Sprintf("expected int, but got %s", ast.TypeString(args[3]))}
		}
		if start < 0 || len(s) <= int(end) {
			return nil, &Error{Message: fmt.Sprintf("out of range: %d %d", start, end)}
		}
		s = s[start:end]
	}
	return implWriteString(port, s)
}

func portDisplay(args []ast.Expr) (ast.Expr, error) {
	if len(args) == 0 {
		return nil, &Error{Message: fmt.Sprintf("required at least 1, but got %d", len(args))}
	}

	var s ast.StringExpr
	e := args[0]
	if implIsString(e) {
		s = e.(ast.StringExpr)
	} else if implIsNumber(e) {
		s = ast.NumberToString(e)
	} else {
		return nil, &Error{Message: fmt.Sprintf("expected string, but got %s", ast.TypeString(args[0]))}
	}

	var ok bool
	port := CurrentVM.Output
	if len(args) > 1 {
		port, ok = args[1].(ast.OutputPort)
		if !ok {
			return nil, &Error{Message: fmt.Sprintf("expected port, but got %s", ast.TypeString(args[1]))}
		}
	}
	return implWriteString(port, s)
}

func portNewline(args []ast.Expr) (ast.Expr, error) {
	port := CurrentVM.Output
	var ok bool
	if len(args) == 1 {
		port, ok = args[1].(ast.OutputPort)
		if !ok {
			return nil, &Error{Message: fmt.Sprintf("expected port, but got %s", ast.TypeString(args[1]))}
		}
	}
	return implWriteString(port, "\n")
}

func portIsInputPort(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("required 1, but got %d", len(args))}
	}
	if _, ok := args[0].(ast.InputPort); ok {
		return ast.BooleanExpr{true}, nil
	} else {
		return ast.BooleanExpr{false}, nil
	}
}

func portIsOutputPort(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("required 1, but got %d", len(args))}
	}
	if _, ok := args[0].(ast.OutputPort); ok {
		return ast.BooleanExpr{true}, nil
	} else {
		return ast.BooleanExpr{false}, nil
	}
}

func portIsTextualPort(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("required 1, but got %d", len(args))}
	}
	switch args[0].(type) {
	case ast.OutputPort:
		p := args[0].(ast.OutputPort)
		return ast.BooleanExpr{!p.Binary}, nil
	case ast.InputPort:
		p := args[0].(ast.InputPort)
		return ast.BooleanExpr{!p.Binary}, nil
	default:
		return ast.BooleanExpr{false}, nil
	}
}

func portIsBinaryPort(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("required 1, but got %d", len(args))}
	}
	switch args[0].(type) {
	case ast.OutputPort:
		p := args[0].(ast.OutputPort)
		return ast.BooleanExpr{p.Binary}, nil
	case ast.InputPort:
		p := args[0].(ast.InputPort)
		return ast.BooleanExpr{p.Binary}, nil
	default:
		return ast.BooleanExpr{false}, nil
	}
}

func portIsPort(args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("required 1, but got %d", len(args))}
	}
	switch args[0].(type) {
	case ast.OutputPort:
		return ast.BooleanExpr{true}, nil
	case ast.InputPort:
		return ast.BooleanExpr{true}, nil
	default:
		return ast.BooleanExpr{false}, nil
	}
}
