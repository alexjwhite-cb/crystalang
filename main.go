package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexjwhite-cb/track-lang/pkg/lexer"
)

func main() {
	out, _ := lexer.NewLexer().Lex(strings.Join(os.Args[1:], "\n"))
	fmt.Printf("%+v\n", out)
}
