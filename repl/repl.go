package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/Aergiaaa/simplescript/evaluator"
	"github.com/Aergiaaa/simplescript/lexer"
	"github.com/Aergiaaa/simplescript/object"
	"github.com/Aergiaaa/simplescript/parser"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	buffer := bufio.NewScanner(in)
	env := object.InitEnv()
	macroEnv := object.InitEnv()

	for {
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

		l := lexer.InitLexer(line)
		p := parser.InitParser(l)

		program := p.Parse()
		if len(p.Errors()) != 0 {
			printParseError(out, p.Errors())
			continue
		}

		evaluator.DefineMacros(program, macroEnv)
		expanded, err := evaluator.ExpandMacros(program, macroEnv)
		if err != nil {
			io.WriteString(out, err.Error())
			io.WriteString(out, "\n")
		}

		evaled := evaluator.Eval(expanded, env)
		if evaled != nil {
			io.WriteString(out, evaled.Inspect())
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
