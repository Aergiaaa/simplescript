package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Aergiaaa/idiotic_interpreter/evaluator"
	"github.com/Aergiaaa/idiotic_interpreter/lexer"
	"github.com/Aergiaaa/idiotic_interpreter/object"
	"github.com/Aergiaaa/idiotic_interpreter/parser"
	"github.com/Aergiaaa/idiotic_interpreter/repl"
)

func main() {
	if len(os.Args) >= 3 && os.Args[1] == "run" {
		filename := os.Args[2]

		// check the extension
		if filepath.Ext(filename) != ".il" {
			fmt.Fprintf(os.Stderr, "Error file must have .il extension\n")
			os.Exit(1)
		}

		// exec file mode
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		input := string(content)

		env := object.InitEnv()
		l := lexer.InitLexer(input)
		p := parser.InitParser(l)
		program := p.Parse()

		if len(p.Errors()) != 0 {
			for _, err := range p.Errors() {
				fmt.Fprintf(os.Stderr, "Parser Error: %s\n", err)
			}
			os.Exit(1)
		}

		result := evaluator.Eval(program, env)
		if result != nil && result.Type() == object.ERR_OBJ {
			fmt.Fprintf(os.Stderr, "%s\n", result.Inspect())
			os.Exit(1)
		}
		return
	}

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! welcome to idiotic stupid language!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
