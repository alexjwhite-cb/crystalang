package repl

import (
	"bufio"
	"fmt"
	"github.com/alexjwhite-cb/jet/pkg/lexer"
	"github.com/alexjwhite-cb/jet/pkg/parser"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		fmt.Printf("%s\n", program.String())
		//for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		//	fmt.Fprintf(out, "%+v\n", tok)
		//}
	}
}
