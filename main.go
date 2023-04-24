package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexjwhite-cb/track-lang/lexer"
)

func main() {
	out, _ := lexer.Lex(strings.Join(os.Args[1:], "\n"))
	fmt.Printf("%+v\n", out)
}
