package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexjwhite-cb/jet/pkg/abstracter"
	"github.com/alexjwhite-cb/jet/pkg/token"
)

func main() {
	out, err := token.NewTokenizer().Tokenize(strings.Join(os.Args[1:], "\n"))
	if err != nil {
		panic(err)
	}
	ast := abstracter.Abstract(out.Tokens)
	fmt.Printf("%+v\n", ast)
}
