package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexjwhite-cb/jet/pkg/abstracter"
	"github.com/alexjwhite-cb/jet/pkg/tokenizer"
)

func main() {
	out, err := tokenizer.NewTokenizer().Tokenize(strings.Join(os.Args[1:], "\n"))
	if err != nil {
		panic(err)
	}
	ast := abstracter.NewAbstractSyntaxTree().Abstract(out.Tokens)
	fmt.Printf("%+v\n", ast)
}
