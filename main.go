package main

import (
	"fmt"
	"github.com/alexjwhite-cb/jet/pkg/repl"
	"os"
)

func main() {
	fmt.Printf("Welcome to the Jet programming language!\n")
	repl.Start(os.Stdin, os.Stdout)
}
