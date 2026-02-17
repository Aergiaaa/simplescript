package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/Aergiaaa/simplescript/ast"
	"github.com/Aergiaaa/simplescript/compiler"
	"github.com/Aergiaaa/simplescript/evaluator"
	"github.com/Aergiaaa/simplescript/lexer"
	"github.com/Aergiaaa/simplescript/object"
	"github.com/Aergiaaa/simplescript/parser"
	"github.com/Aergiaaa/simplescript/vm"
)

const PROMPT = ">>"

type COMPILE_MODE bool

const (
	COMPILE   COMPILE_MODE = true
	INTERPRET COMPILE_MODE = false
)

func Start(in io.Reader, out io.Writer, mode COMPILE_MODE) {
	buffer := bufio.NewScanner(in)
	env := object.InitEnv()
	macroEnv := object.InitEnv()

	for {
		var err error

		fmt.Printf(PROMPT)

		readed := buffer.Scan()
		if !readed {
			return
		}

		line := buffer.Text()

		for strings.HasSuffix(line, "\\") {
			line = strings.TrimSuffix(line, "\\")
			fmt.Printf("..")

			if !buffer.Scan() {
				return
			}

			line += "\n" + buffer.Text()
		}

		// init
		l := lexer.Init(line)
		p := parser.Init(l)

		// parsing program
		program := p.Parse()
		if len(p.Errors()) != 0 {
			printParseError(out, p.Errors())
			continue
		}

		var expanded ast.Node
		if mode == INTERPRET {
			// evaluate macros
			evaluator.DefineMacros(program, macroEnv)
			expanded, err = evaluator.ExpandMacros(program, macroEnv)
			if err != nil {
				io.WriteString(out, err.Error())
				io.WriteString(out, "\n")
			}
			// evaluating
			evaled := evaluator.Eval(expanded, env)
			if evaled != nil {
				// io.WriteString(out, evaled.Inspect())
				// io.WriteString(out, "\n")
			}
		}

		if mode == COMPILE {
			// init compiler
			comp := compiler.Init()
			err = comp.Compile(program)
			if err != nil {
				fmt.Fprintf(out, "compilation failed:\n %s\n", err)
				continue
			}

			// init vm
			machine := vm.Init(comp.Bytecode())
			err = machine.Run()
			if err != nil {
				fmt.Fprintf(out, "failed executing bytecode:\n %s\n", err)
				continue
			}

			lastPop := machine.LastPoppedStackElem()
			io.WriteString(out, lastPop.Inspect())
			io.WriteString(out, "\n")
		}

	}
}

func printParseError(out io.Writer, errors []string) {
	io.WriteString(out, "parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
