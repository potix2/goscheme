package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/potix2/goscheme/parser"
	"github.com/potix2/goscheme/vm"

	"github.com/sirupsen/logrus"
)

type Config struct {
	debugLog string
}

func runInteractive() {
	var following bool
	var source string

	scanner := bufio.NewScanner(os.Stdin)

	env := vm.NewEnv()
	vm.SetInteractionEnvironment(env)
	vm.InitVM(env)
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

		exprs, err := parser.Read(source)
		if e, ok := err.(*parser.Error); ok {
			es := e.Error()
			if strings.HasSuffix(es, "unexpected $end") {
				following = true
				source += "\n"
				continue
			} else {
				following = true
				fmt.Fprintln(os.Stderr, e)
				continue
			}
		}

		following = false

		for _, expr := range exprs {
			v, err := vm.Eval(expr, env)
			if e, ok := err.(*vm.Error); ok {
				fmt.Fprintln(os.Stderr, e)
				continue
			}

			var b bytes.Buffer
			vm.Print(&b, v)
			b.Write([]byte("\n"))
			b.WriteTo(os.Stdout)
		}
		source = ""
	}

	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			logrus.Fatalf("ReadString error: %v", err)
		}
	}
}

func runFromFile(sourcePath string) {
	source, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		logrus.Fatal(err)
	}

	code := string(source)
	logrus.Debug(code)

	env := vm.NewEnv()
	vm.InitVM(env)
	vm.SetupPrimitives(env)
	exprs, err := parser.Read(code)
	if err != nil {
		logrus.Fatalf("parse error: %v", err)
	}

	for _, expr := range exprs {
		v, err := vm.Eval(expr, env)
		if err != nil {
			logrus.Fatalf("eval error: %v", err)
		}

		var b bytes.Buffer
		vm.Print(&b, v)
		logrus.Debug(b.String())
	}
}

func main() {
	c := &Config{}

	flag.StringVar(&c.debugLog, "debuglog", "", "Set a path of debug log")
	flag.Parse()

	if c.debugLog != "" {
		logrus.SetLevel(logrus.DebugLevel)

		f, err := os.OpenFile(c.debugLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			logrus.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		logrus.SetOutput(f)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	if len(flag.Args()) == 0 {
		runInteractive()
	} else {
		for _, f := range flag.Args() {
			runFromFile(f)
		}
	}
}
