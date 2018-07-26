package vm

import (
	"fmt"

	"github.com/potix2/goscheme/scm"
)

func portCurrentInputPort(args []scm.Object) (scm.Object, error) {
	return CurrentVM.Input, nil
}

func portCurrentOutputPort(args []scm.Object) (scm.Object, error) {
	return CurrentVM.Output, nil
}

func portCurrentErrorPort(args []scm.Object) (scm.Object, error) {
	return CurrentVM.Error, nil
}

func implWriteString(p scm.OutputPort, s scm.String) (scm.Object, error) {
	n, err := p.Writer.Write([]byte(s))
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("failed to write: %v", err)}
	}
	if n != len(string(s)) {
		return nil, &Error{Message: fmt.Sprintf("actual written size: %d, expected %d", n, len(string(s)))}
	}

	return scm.Undefined{}, nil
}

func portWriteString(args []scm.Object) (scm.Object, error) {
	if len(args) == 0 {
		return nil, &Error{Message: fmt.Sprintf("required at lescm 1, but got %d", len(args))}
	}
	s, ok := args[0].(scm.String)
	if !ok {
		return nil, &Error{Message: fmt.Sprintf("expected string, but got %s", scm.TypeString(args[0]))}
	}

	port := CurrentVM.Output
	if len(args) > 1 {
		port, ok = args[1].(scm.OutputPort)
		if !ok {
			return nil, &Error{Message: fmt.Sprintf("expected port, but got %s", scm.TypeString(args[1]))}
		}
	}

	if len(args) == 3 {
		start, ok := args[2].(scm.IntNum)
		if !ok {
			return nil, &Error{Message: fmt.Sprintf("expected int, but got %s", scm.TypeString(args[2]))}
		}
		s = s[start:]
	} else if len(args) == 4 {
		start, ok := args[2].(scm.IntNum)
		if !ok {
			return nil, &Error{Message: fmt.Sprintf("expected int, but got %s", scm.TypeString(args[2]))}
		}

		end, ok := args[3].(scm.IntNum)
		if !ok {
			return nil, &Error{Message: fmt.Sprintf("expected int, but got %s", scm.TypeString(args[3]))}
		}
		if start < 0 || len(s) <= int(end) {
			return nil, &Error{Message: fmt.Sprintf("out of range: %d %d", start, end)}
		}
		s = s[start:end]
	}
	return implWriteString(port, s)
}

func portDisplay(args []scm.Object) (scm.Object, error) {
	if len(args) == 0 {
		return nil, &Error{Message: fmt.Sprintf("required at lescm 1, but got %d", len(args))}
	}

	e := args[0]
	s := scm.String(dumpExpr(e))

	var ok bool
	port := CurrentVM.Output
	if len(args) > 1 {
		port, ok = args[1].(scm.OutputPort)
		if !ok {
			return nil, &Error{Message: fmt.Sprintf("expected port, but got %s", scm.TypeString(args[1]))}
		}
	}
	return implWriteString(port, s)
}

func portNewline(args []scm.Object) (scm.Object, error) {
	port := CurrentVM.Output
	var ok bool
	if len(args) == 1 {
		port, ok = args[1].(scm.OutputPort)
		if !ok {
			return nil, &Error{Message: fmt.Sprintf("expected port, but got %s", scm.TypeString(args[1]))}
		}
	}
	return implWriteString(port, "\n")
}

func portIsInputPort(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("required 1, but got %d", len(args))}
	}
	if _, ok := args[0].(scm.InputPort); ok {
		return scm.Boolean{true}, nil
	} else {
		return scm.Boolean{false}, nil
	}
}

func portIsOutputPort(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("required 1, but got %d", len(args))}
	}
	if _, ok := args[0].(scm.OutputPort); ok {
		return scm.Boolean{true}, nil
	} else {
		return scm.Boolean{false}, nil
	}
}

func portIsTextualPort(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("required 1, but got %d", len(args))}
	}
	switch args[0].(type) {
	case scm.OutputPort:
		p := args[0].(scm.OutputPort)
		return scm.Boolean{!p.Binary}, nil
	case scm.InputPort:
		p := args[0].(scm.InputPort)
		return scm.Boolean{!p.Binary}, nil
	default:
		return scm.Boolean{false}, nil
	}
}

func portIsBinaryPort(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("required 1, but got %d", len(args))}
	}
	switch args[0].(type) {
	case scm.OutputPort:
		p := args[0].(scm.OutputPort)
		return scm.Boolean{p.Binary}, nil
	case scm.InputPort:
		p := args[0].(scm.InputPort)
		return scm.Boolean{p.Binary}, nil
	default:
		return scm.Boolean{false}, nil
	}
}

func portIsPort(args []scm.Object) (scm.Object, error) {
	if len(args) != 1 {
		return nil, &Error{Message: fmt.Sprintf("required 1, but got %d", len(args))}
	}
	switch args[0].(type) {
	case scm.OutputPort:
		return scm.Boolean{true}, nil
	case scm.InputPort:
		return scm.Boolean{true}, nil
	default:
		return scm.Boolean{false}, nil
	}
}
