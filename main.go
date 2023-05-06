package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexjwhite-cb/crystalang/pkg/tokeniser"
)

func main() {
	out, err := tokeniser.NewTokeniser().Tokenise(strings.Join(os.Args[1:], "\n"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", out)
}
