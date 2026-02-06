package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Aergiaaa/simplescript/evaluator"
	"github.com/Aergiaaa/simplescript/lexer"
	"github.com/Aergiaaa/simplescript/object"
	"github.com/Aergiaaa/simplescript/parser"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	buffer := bufio.NewScanner(in)
	env := object.InitEnv()

	for {
		fmt.Printf(PROMPT)

		readed := buffer.Scan()
		if !readed {
			return
		}

		line := buffer.Text()
		l := lexer.InitLexer(line)
		p := parser.InitParser(l)

		program := p.Parse()
		if len(p.Errors()) != 0 {
			printParseError(out, p.Errors())
			continue
		}

		evaled := evaluator.Eval(program, env)
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
