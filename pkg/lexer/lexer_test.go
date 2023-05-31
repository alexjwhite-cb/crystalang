package lexer

import (
	"fmt"
	"github.com/alexjwhite-cb/jet/pkg/token"
	"reflect"
	"strings"
	"testing"
)

const (
	entry = `meth main {
}`
	fail = `meth main {
	(1)->
}`
	returnVar = `meth main {
	num = 0
	(num)->
}`
	returnIncrement1 = `meth main {
	num = 0
	num++
	(num)->
}`
	returnIncrement2 = `meth main {
	num = 0
	num += 2
	(num)->
}`
	returnDecrement1 = `meth main {
	num = 0
	num--
	(num)->
}`
	returnDecrement2 = `meth main {
	num = 0
	num -= 2
	(num)->
}`
	returnString = `meth main {
	str = "Hello, World!"
	(str)->
}`
	returnEscapeString = `meth main {
	str = "\"Hello, World!\""
	(str)->
}`
	andOperator = `meth main {
	if a == 2 * 2 && !b {
		(true)->
	} else {
		(false)->
	}
}`
	orOperator = `meth main {
	if a == 2 * 2 || !b {
		(true)->
	}
	(false)->
}`
	func1 = `meth NewGuitar: tuning {
	guitar = Guitar->new
	tuning = tuning->toUpper
	if !inValidTunings(tuning) {
		(error: "{tuning} is not a valid tuning")->
	}
	for i, t in array {
		if t->len == 1 {
			t = " " + t
		}
		guitar.Tuning[i+1] = t
	}
	(guitar)->
}`
)

func TestLexer_NextToken(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		expect []token.Token
	}{
		{
			name: "Entrypoint",
			in:   entry,
			expect: []token.Token{
				{token.METHOD, "meth", 1, 1},
				{token.IDENT, "main", 6, 1},
				{token.LBRACE, "{", 11, 1},
				{token.NEWLINE, "\n", 12, 1},
				{token.RBRACE, "}", 1, 2},
			},
		},
		{
			name: "Return 1",
			in:   fail,
			expect: []token.Token{
				{token.METHOD, "meth", 1, 1},
				{token.IDENT, "main", 6, 1},
				{token.LBRACE, "{", 11, 1},
				{token.NEWLINE, "\n", 12, 1},
				{token.LPAREN, "(", 2, 2},
				{token.INT, "1", 3, 2},
				{token.RPAREN, ")", 4, 2},
				{token.PASSTHROUGH, "->", 5, 2},
				{token.NEWLINE, "\n", 7, 2},
				{token.RBRACE, "}", 1, 3},
			},
		},
		{
			name: "Return Var",
			in:   returnVar,
			expect: []token.Token{
				{token.METHOD, "meth", 1, 1},
				{token.IDENT, "main", 6, 1},
				{token.LBRACE, "{", 11, 1},
				{token.NEWLINE, "\n", 12, 1},
				{token.IDENT, "num", 2, 2},
				{token.ASSIGN, "=", 6, 2},
				{token.INT, "0", 8, 2},
				{token.NEWLINE, "\n", 9, 2},
				{token.LPAREN, "(", 2, 3},
				{token.IDENT, "num", 3, 3},
				{token.RPAREN, ")", 6, 3},
				{token.PASSTHROUGH, "->", 7, 3},
				{token.NEWLINE, "\n", 9, 3},
				{token.RBRACE, "}", 1, 4},
			},
		},
		{
			name: "Increment Var 1",
			in:   returnIncrement1,
			expect: []token.Token{
				{token.METHOD, "meth", 1, 1},
				{token.IDENT, "main", 6, 1},
				{token.LBRACE, "{", 11, 1},
				{token.NEWLINE, "\n", 12, 1},
				{token.IDENT, "num", 2, 2},
				{token.ASSIGN, "=", 6, 2},
				{token.INT, "0", 8, 2},
				{token.NEWLINE, "\n", 9, 2},
				{token.IDENT, "num", 2, 3},
				{token.INCREMENT, "++", 5, 3},
				{token.NEWLINE, "\n", 7, 3},
				{token.LPAREN, "(", 2, 4},
				{token.IDENT, "num", 3, 4},
				{token.RPAREN, ")", 6, 4},
				{token.PASSTHROUGH, "->", 7, 4},
				{token.NEWLINE, "\n", 9, 4},
				{token.RBRACE, "}", 1, 5},
			},
		},
		{
			name: "Increment Var 2",
			in:   returnIncrement2,
			expect: []token.Token{
				{token.METHOD, "meth", 1, 1},
				{token.IDENT, "main", 6, 1},
				{token.LBRACE, "{", 11, 1},
				{token.NEWLINE, "\n", 12, 1},
				{token.IDENT, "num", 2, 2},
				{token.ASSIGN, "=", 6, 2},
				{token.INT, "0", 8, 2},
				{token.NEWLINE, "\n", 9, 2},
				{token.IDENT, "num", 2, 3},
				{token.PLUSASSIGN, "+=", 6, 3},
				{token.INT, "2", 9, 3},
				{token.NEWLINE, "\n", 10, 3},
				{token.LPAREN, "(", 2, 4},
				{token.IDENT, "num", 3, 4},
				{token.RPAREN, ")", 6, 4},
				{token.PASSTHROUGH, "->", 7, 4},
				{token.NEWLINE, "\n", 9, 4},
				{token.RBRACE, "}", 1, 5},
			},
		},
		{
			name: "Decrement Var 1",
			in:   returnDecrement1,
			expect: []token.Token{
				{token.METHOD, "meth", 1, 1},
				{token.IDENT, "main", 6, 1},
				{token.LBRACE, "{", 11, 1},
				{token.NEWLINE, "\n", 12, 1},
				{token.IDENT, "num", 2, 2},
				{token.ASSIGN, "=", 6, 2},
				{token.INT, "0", 8, 2},
				{token.NEWLINE, "\n", 9, 2},
				{token.IDENT, "num", 2, 3},
				{token.DECREMENT, "--", 5, 3},
				{token.NEWLINE, "\n", 7, 3},
				{token.LPAREN, "(", 2, 4},
				{token.IDENT, "num", 3, 4},
				{token.RPAREN, ")", 6, 4},
				{token.PASSTHROUGH, "->", 7, 4},
				{token.NEWLINE, "\n", 9, 4},
				{token.RBRACE, "}", 1, 5},
			},
		},
		{
			name: "Decrement Var 2",
			in:   returnDecrement2,
			expect: []token.Token{
				{token.METHOD, "meth", 1, 1},
				{token.IDENT, "main", 6, 1},
				{token.LBRACE, "{", 11, 1},
				{token.NEWLINE, "\n", 12, 1},
				{token.IDENT, "num", 2, 2},
				{token.ASSIGN, "=", 6, 2},
				{token.INT, "0", 8, 2},
				{token.NEWLINE, "\n", 9, 2},
				{token.IDENT, "num", 2, 3},
				{token.MINUSASSIGN, "-=", 6, 3},
				{token.INT, "2", 9, 3},
				{token.NEWLINE, "\n", 10, 3},
				{token.LPAREN, "(", 2, 4},
				{token.IDENT, "num", 3, 4},
				{token.RPAREN, ")", 6, 4},
				{token.PASSTHROUGH, "->", 7, 4},
				{token.NEWLINE, "\n", 9, 4},
				{token.RBRACE, "}", 1, 5},
			},
		},
		{
			name: "Return String",
			in:   returnString,
			expect: []token.Token{
				{token.METHOD, "meth", 1, 1},
				{token.IDENT, "main", 6, 1},
				{token.LBRACE, "{", 11, 1},
				{token.NEWLINE, "\n", 12, 1},
				{token.IDENT, "str", 2, 2},
				{token.ASSIGN, "=", 6, 2},
				{token.STRING, "Hello, World!", 8, 2},
				{token.NEWLINE, "\n", 23, 2},
				{token.LPAREN, "(", 2, 3},
				{token.IDENT, "str", 3, 3},
				{token.RPAREN, ")", 6, 3},
				{token.PASSTHROUGH, "->", 7, 3},
				{token.NEWLINE, "\n", 9, 3},
				{token.RBRACE, "}", 1, 4},
			},
		},
		{
			name: "Return EscapeString",
			in:   returnEscapeString,
			expect: []token.Token{
				{token.METHOD, "meth", 1, 1},
				{token.IDENT, "main", 6, 1},
				{token.LBRACE, "{", 11, 1},
				{token.NEWLINE, "\n", 12, 1},
				{token.IDENT, "str", 2, 2},
				{token.ASSIGN, "=", 6, 2},
				{token.STRING, "\\\"Hello, World!\\\"", 8, 2},
				{token.NEWLINE, "\n", 27, 2},
				{token.LPAREN, "(", 2, 3},
				{token.IDENT, "str", 3, 3},
				{token.RPAREN, ")", 6, 3},
				{token.PASSTHROUGH, "->", 7, 3},
				{token.NEWLINE, "\n", 9, 3},
				{token.RBRACE, "}", 1, 4},
			},
		},
		{
			name: "And Operator",
			in:   andOperator,
			expect: []token.Token{
				{token.METHOD, "meth", 1, 1},
				{token.IDENT, "main", 6, 1},
				{token.LBRACE, "{", 11, 1},
				{token.NEWLINE, "\n", 12, 1},
				{token.IF, "if", 2, 2},
				{token.IDENT, "a", 5, 2},
				{token.EQUAL, "==", 7, 2},
				{token.INT, "2", 10, 2},
				{token.MULTIPLY, "*", 12, 2},
				{token.INT, "2", 14, 2},
				{token.AND, "&&", 16, 2},
				{token.NOT, "!", 19, 2},
				{token.IDENT, "b", 20, 2},
				{token.LBRACE, "{", 22, 2},
				{token.NEWLINE, "\n", 23, 2},
				{token.LPAREN, "(", 3, 3},
				{token.TRUE, "true", 4, 3},
				{token.RPAREN, ")", 8, 3},
				{token.PASSTHROUGH, "->", 9, 3},
				{token.NEWLINE, "\n", 11, 3},
				{token.RBRACE, "}", 2, 4},
				{token.ELSE, "else", 4, 4},
				{token.LBRACE, "{", 9, 4},
				{token.NEWLINE, "\n", 10, 4},
				{token.LPAREN, "(", 3, 5},
				{token.FALSE, "false", 4, 5},
				{token.RPAREN, ")", 9, 5},
				{token.PASSTHROUGH, "->", 10, 5},
				{token.NEWLINE, "\n", 12, 5},
				{token.RBRACE, "}", 2, 6},
				{token.NEWLINE, "\n", 3, 6},
				{token.RBRACE, "}", 1, 7},
			},
		},
		{
			name: "Or Operator",
			in:   orOperator,
			expect: []token.Token{
				{token.METHOD, "meth", 1, 1},
				{token.IDENT, "main", 6, 1},
				{token.LBRACE, "{", 11, 1},
				{token.NEWLINE, "\n", 12, 1},
				{token.IF, "if", 2, 2},
				{token.IDENT, "a", 5, 2},
				{token.EQUAL, "==", 7, 2},
				{token.INT, "2", 10, 2},
				{token.MULTIPLY, "*", 12, 2},
				{token.INT, "2", 14, 2},
				{token.OR, "||", 16, 2},
				{token.NOT, "!", 19, 2},
				{token.IDENT, "b", 20, 2},
				{token.LBRACE, "{", 22, 2},
				{token.NEWLINE, "\n", 23, 2},
				{token.LPAREN, "(", 3, 3},
				{token.TRUE, "true", 4, 3},
				{token.RPAREN, ")", 8, 3},
				{token.PASSTHROUGH, "->", 9, 3},
				{token.NEWLINE, "\n", 11, 3},
				{token.RBRACE, "}", 2, 4},
				{token.NEWLINE, "\n", 3, 4},
				{token.LPAREN, "(", 2, 5},
				{token.FALSE, "false", 3, 5},
				{token.RPAREN, ")", 8, 5},
				{token.PASSTHROUGH, "->", 9, 5},
				{token.NEWLINE, "\n", 11, 5},
				{token.RBRACE, "}", 1, 6},
			},
		},
		{
			name: "Simple Function",
			in:   func1,
			expect: []token.Token{
				{token.METHOD, "meth", 1, 1},
				{token.IDENT, "NewGuitar", 6, 1},
				{token.COLON, ":", 15, 1},
				{token.IDENT, "tuning", 17, 1},
				{token.LBRACE, "{", 24, 1},
				{token.NEWLINE, "\n", 25, 1},
				{token.IDENT, "guitar", 2, 2},
				{token.ASSIGN, "=", 9, 2},
				{token.IDENT, "Guitar", 11, 2},
				{token.PASSTHROUGH, "->", 17, 2},
				{token.IDENT, "new", 19, 2},
				{token.NEWLINE, "\n", 22, 2},
				{token.IDENT, "tuning", 2, 3},
				{token.ASSIGN, "=", 9, 3},
				{token.IDENT, "tuning", 11, 3},
				{token.PASSTHROUGH, "->", 17, 3},
				{token.IDENT, "toUpper", 19, 3},
				{token.NEWLINE, "\n", 26, 3},
				{token.IF, "if", 2, 4},
				{token.NOT, "!", 5, 4},
				{token.IDENT, "inValidTunings", 6, 4},
				{token.LPAREN, "(", 20, 4},
				{token.IDENT, "tuning", 21, 4},
				{token.RPAREN, ")", 27, 4},
				{token.LBRACE, "{", 29, 4},
				{token.NEWLINE, "\n", 30, 4},
				{token.LPAREN, "(", 3, 5},
				{token.ERROR, "error", 4, 5},
				{token.COLON, ":", 9, 5},
				{token.STRING, "{tuning} is not a valid tuning", 11, 5},
				{token.RPAREN, ")", 43, 5},
				{token.PASSTHROUGH, "->", 44, 5},
				{token.NEWLINE, "\n", 46, 5},
				{token.RBRACE, "}", 2, 6},
				{token.NEWLINE, "\n", 3, 6},
				{token.FOR, "for", 2, 7},
				{token.IDENT, "i", 6, 7},
				{token.COMMA, ",", 7, 7},
				{token.IDENT, "t", 9, 7},
				{token.IN, "in", 11, 7},
				{token.IDENT, "array", 14, 7},
				{token.LBRACE, "{", 20, 7},
				{token.NEWLINE, "\n", 21, 7},
				{token.IF, "if", 3, 8},
				{token.IDENT, "t", 6, 8},
				{token.PASSTHROUGH, "->", 7, 8},
				{token.IDENT, "len", 9, 8},
				{token.EQUAL, "==", 13, 8},
				{token.INT, "1", 16, 8},
				{token.LBRACE, "{", 18, 8},
				{token.NEWLINE, "\n", 19, 8},
				{token.IDENT, "t", 4, 9},
				{token.ASSIGN, "=", 6, 9},
				{token.STRING, " ", 8, 9},
				{token.PLUS, "+", 12, 9},
				{token.IDENT, "t", 14, 9},
				{token.NEWLINE, "\n", 15, 9},
				{token.RBRACE, "}", 3, 10},
				{token.NEWLINE, "\n", 4, 10},
				{token.IDENT, "guitar", 3, 11},
				{token.STOP, ".", 9, 11},
				{token.IDENT, "Tuning", 10, 11},
				{token.LBRACK, "[", 16, 11},
				{token.IDENT, "i", 17, 11},
				{token.PLUS, "+", 18, 11},
				{token.INT, "1", 19, 11},
				{token.RBRACK, "]", 20, 11},
				{token.ASSIGN, "=", 22, 11},
				{token.IDENT, "t", 24, 11},
				{token.NEWLINE, "\n", 25, 11},
				{token.RBRACE, "}", 2, 12},
				{token.NEWLINE, "\n", 3, 12},
				{token.LPAREN, "(", 2, 13},
				{token.IDENT, "guitar", 3, 13},
				{token.RPAREN, ")", 9, 13},
				{token.PASSTHROUGH, "->", 10, 13},
				{token.NEWLINE, "\n", 12, 13},
				{token.RBRACE, "}", 1, 14},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result []token.Token
			l := New(tt.in)
			for l.char != 0 {
				result = append(result, l.NextToken())
			}

			if !reflect.DeepEqual(tt.expect, result) {
				t.Errorf(strings.ReplaceAll(fmt.Sprintf("\nExpected: %+v\nGot: %+v", tt.expect, result), "map", ""))
			}
		})
	}
}
