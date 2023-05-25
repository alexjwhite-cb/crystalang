package ast

import (
	"github.com/alexjwhite-cb/jet/pkg/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Stmt{
			&ValueStmt{
				Token: token.Token{Type: token.IDENT, Literal: "x"},
				Name: &Ident{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Value: "x",
				},
				Value: &Ident{
					Token: token.Token{Type: token.IDENT, Literal: "y"},
					Value: "y",
				},
			},
		},
	}

	if program.String() != "x = y" {
		t.Errorf("program.String() wrong. got '%s'", program.String())
	}
}