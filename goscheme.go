package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/potix2/goscheme/parser"
	"github.com/potix2/goscheme/vm"
)

func main() {
	var following bool
	var source string
	scanner := bufio.NewScanner(os.Stdin)

	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	env := vm.NewEnv()
	vm.SetupPrimitives(env)
	for {
		if !following {
			fmt.Print("> ")
		}

		if !scanner.Scan() {
			break
		}

		source += scanner.Text()
		if source == "" {
			continue
		}
		if source == "(exit)" {
			break
		}

		expr, err := parser.Read(source)
		if e, ok := err.(*parser.Error); ok {
			es := e.Error()
			if strings.HasSuffix(es, "unexpected $end") {
				following = true
				continue
			} else {
				following = true
				fmt.Fprintln(os.Stderr, e)
				continue
			}
		}

		following = false

		v, err := vm.Eval(expr, env)
		if e, ok := err.(*vm.Error); ok {
			fmt.Fprintln(os.Stderr, e)
			continue
		}

		var b bytes.Buffer
		vm.Print(&b, v)
		b.Write([]byte("\n"))
		b.WriteTo(os.Stdout)
		source = ""
	}

	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			fmt.Fprintf(os.Stderr, "ReadString error:", err)
			//return 12
			return
		}
	}
}
