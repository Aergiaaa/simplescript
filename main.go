package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/Aergiaaa/idiotic_interpreter/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! welcome to idiotic stupid language!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
