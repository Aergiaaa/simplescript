package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Aergiaaa/simplescript/evaluator"
	"github.com/Aergiaaa/simplescript/lexer"
	"github.com/Aergiaaa/simplescript/object"
	"github.com/Aergiaaa/simplescript/parser"
	"github.com/Aergiaaa/simplescript/repl"
)

func main() {

	// default `idiot run <script.simp>`
	if len(os.Args) >= 3 && os.Args[1] == "run" {
		readFile()
	}

	// repl mode
	var mode repl.COMPILE_MODE = repl.COMPILE
	if os.Args[1] == "i" || os.Args[1] == "c" {
		switch os.Args[1] {
		case "i":
			mode = repl.INTERPRET
		case "c":
			mode = repl.COMPILE
		}
	}

	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! welcome to idiotic stupid language!\n", user.Username)
	fmt.Printf("MODE: %s\n", mode.String())

	repl.Start(os.Stdin, os.Stdout, mode)
}

func readFile() {
	filename := os.Args[2]

	// check the extension
	if filepath.Ext(filename) != ".simp" {
		fmt.Fprintf(os.Stderr, "Error file must have .simp extension\n")
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
	l := lexer.Init(input)
	p := parser.Init(l)
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
}
