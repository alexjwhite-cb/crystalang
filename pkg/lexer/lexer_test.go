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
	if tuning->!inValidTunings {
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
				{token.METHOD, "meth", 0, 1},
				{token.IDENT, "main", 5, 1},
				{token.LBRACE, "{", 10, 1},
				{token.RBRACE, "}", 12, 2},
			},
		},
		{
			name: "Return 1",
			in:   fail,
			expect: []token.Token{
				{token.METHOD, "meth", 0, 1},
				{token.IDENT, "main", 5, 1},
				{token.LBRACE, "{", 10, 1},
				{token.LPAREN, "(", 13, 2},
				{token.INT, "1", 14, 2},
				{token.RPAREN, ")", 15, 2},
				{token.PASSTHROUGH, "->", 16, 2},
				{token.RBRACE, "}", 19, 3},
			},
		},
		{
			name: "Return Var",
			in:   returnVar,
			expect: []token.Token{
				{token.METHOD, "meth", 0, 1},
				{token.IDENT, "main", 5, 1},
				{token.LBRACE, "{", 10, 1},
				{token.IDENT, "num", 13, 2},
				{token.ASSIGN, "=", 17, 2},
				{token.INT, "0", 19, 2},
				{token.LPAREN, "(", 22, 3},
				{token.IDENT, "num", 23, 3},
				{token.RPAREN, ")", 26, 3},
				{token.PASSTHROUGH, "->", 27, 3},
				{token.RBRACE, "}", 30, 4},
			},
		},
		{
			name: "Increment Var 1",
			in:   returnIncrement1,
			expect: []token.Token{
				{token.METHOD, "meth", 0, 1},
				{token.IDENT, "main", 5, 1},
				{token.LBRACE, "{", 10, 1},
				{token.IDENT, "num", 13, 2},
				{token.ASSIGN, "=", 17, 2},
				{token.INT, "0", 19, 2},
				{token.IDENT, "num", 22, 3},
				{token.INCREMENT, "++", 25, 3},
				{token.LPAREN, "(", 29, 4},
				{token.IDENT, "num", 30, 4},
				{token.RPAREN, ")", 33, 4},
				{token.PASSTHROUGH, "->", 34, 4},
				{token.RBRACE, "}", 37, 5},
			},
		},
		{
			name: "Increment Var 2",
			in:   returnIncrement2,
			expect: []token.Token{
				{token.METHOD, "meth", 0, 1},
				{token.IDENT, "main", 5, 1},
				{token.LBRACE, "{", 10, 1},
				{token.IDENT, "num", 13, 2},
				{token.ASSIGN, "=", 17, 2},
				{token.INT, "0", 19, 2},
				{token.IDENT, "num", 22, 3},
				{token.PLUSASSIGN, "+=", 26, 3},
				{token.INT, "2", 29, 3},
				{token.LPAREN, "(", 32, 4},
				{token.IDENT, "num", 33, 4},
				{token.RPAREN, ")", 36, 4},
				{token.PASSTHROUGH, "->", 37, 4},
				{token.RBRACE, "}", 40, 5},
			},
		},
		{
			name: "Decrement Var 1",
			in:   returnDecrement1,
			expect: []token.Token{
				{token.METHOD, "meth", 0, 1},
				{token.IDENT, "main", 5, 1},
				{token.LBRACE, "{", 10, 1},
				{token.IDENT, "num", 13, 2},
				{token.ASSIGN, "=", 17, 2},
				{token.INT, "0", 19, 2},
				{token.IDENT, "num", 22, 3},
				{token.DECREMENT, "--", 25, 3},
				{token.LPAREN, "(", 29, 4},
				{token.IDENT, "num", 30, 4},
				{token.RPAREN, ")", 33, 4},
				{token.PASSTHROUGH, "->", 34, 4},
				{token.RBRACE, "}", 37, 5},
			},
		},
		{
			name: "Decrement Var 2",
			in:   returnDecrement2,
			expect: []token.Token{
				{token.METHOD, "meth", 0, 1},
				{token.IDENT, "main", 5, 1},
				{token.LBRACE, "{", 10, 1},
				{token.IDENT, "num", 13, 2},
				{token.ASSIGN, "=", 17, 2},
				{token.INT, "0", 19, 2},
				{token.IDENT, "num", 22, 3},
				{token.MINUSASSIGN, "-=", 26, 3},
				{token.INT, "2", 29, 3},
				{token.LPAREN, "(", 32, 4},
				{token.IDENT, "num", 33, 4},
				{token.RPAREN, ")", 36, 4},
				{token.PASSTHROUGH, "->", 37, 4},
				{token.RBRACE, "}", 40, 5},
			},
		},
		{
			name: "Return String",
			in:   returnString,
			expect: []token.Token{
				{token.METHOD, "meth", 0, 1},
				{token.IDENT, "main", 5, 1},
				{token.LBRACE, "{", 10, 1},
				{token.IDENT, "str", 13, 2},
				{token.ASSIGN, "=", 17, 2},
				{token.STRING, "\"Hello, World!\"", 19, 2},
				{token.LPAREN, "(", 36, 3},
				{token.IDENT, "str", 37, 3},
				{token.RPAREN, ")", 40, 3},
				{token.PASSTHROUGH, "->", 41, 3},
				{token.RBRACE, "}", 44, 4},
			},
		},
		{
			name: "Return EscapeString",
			in:   returnEscapeString,
			expect: []token.Token{
				{token.METHOD, "meth", 0, 1},
				{token.IDENT, "main", 5, 1},
				{token.LBRACE, "{", 10, 1},
				{token.IDENT, "str", 13, 2},
				{token.ASSIGN, "=", 17, 2},
				{token.STRING, "\"\\\"Hello, World!\\\"\"", 19, 2},
				{token.LPAREN, "(", 40, 3},
				{token.IDENT, "str", 41, 3},
				{token.RPAREN, ")", 44, 3},
				{token.PASSTHROUGH, "->", 45, 3},
				{token.RBRACE, "}", 48, 4},
			},
		},
		{
			name: "And Operator",
			in:   andOperator,
			expect: []token.Token{
				{token.METHOD, "meth", 0, 1},
				{token.IDENT, "main", 5, 1},
				{token.LBRACE, "{", 10, 1},
				{token.IF, "if", 13, 2},
				{token.IDENT, "a", 16, 2},
				{token.EQUAL, "==", 18, 2},
				{token.INT, "2", 21, 2},
				{token.MULTIPLY, "*", 23, 2},
				{token.INT, "2", 25, 2},
				{token.AND, "&&", 27, 2},
				{token.NOT, "!", 30, 2},
				{token.IDENT, "b", 31, 2},
				{token.LBRACE, "{", 33, 2},
				{token.LPAREN, "(", 37, 3},
				{token.TRUE, "true", 38, 3},
				{token.RPAREN, ")", 42, 3},
				{token.PASSTHROUGH, "->", 43, 3},
				{token.RBRACE, "}", 47, 4},
				{token.ELSE, "else", 49, 4},
				{token.LBRACE, "{", 54, 4},
				{token.LPAREN, "(", 58, 5},
				{token.FALSE, "false", 59, 5},
				{token.RPAREN, ")", 64, 5},
				{token.PASSTHROUGH, "->", 65, 5},
				{token.RBRACE, "}", 69, 6},
				{token.RBRACE, "}", 71, 7},
			},
		},
		{
			name: "Or Operator",
			in:   orOperator,
			expect: []token.Token{
				{token.METHOD, "meth", 0, 1},
				{token.IDENT, "main", 5, 1},
				{token.LBRACE, "{", 10, 1},
				{token.IF, "if", 13, 2},
				{token.IDENT, "a", 16, 2},
				{token.EQUAL, "==", 18, 2},
				{token.INT, "2", 21, 2},
				{token.MULTIPLY, "*", 23, 2},
				{token.INT, "2", 25, 2},
				{token.OR, "||", 27, 2},
				{token.NOT, "!", 30, 2},
				{token.IDENT, "b", 31, 2},
				{token.LBRACE, "{", 33, 2},
				{token.LPAREN, "(", 37, 3},
				{token.TRUE, "true", 38, 3},
				{token.RPAREN, ")", 42, 3},
				{token.PASSTHROUGH, "->", 43, 3},
				{token.RBRACE, "}", 47, 4},
				{token.LPAREN, "(", 50, 5},
				{token.FALSE, "false", 51, 5},
				{token.RPAREN, ")", 56, 5},
				{token.PASSTHROUGH, "->", 57, 5},
				{token.RBRACE, "}", 60, 6},
			},
		},
		{
			name: "Simple Function",
			in:   func1,
			expect: []token.Token{
				{token.METHOD, "meth", 0, 1},
				{token.IDENT, "NewGuitar", 5, 1},
				{token.COLON, ":", 14, 1},
				{token.IDENT, "tuning", 16, 1},
				{token.LBRACE, "{", 23, 1},
				{token.IDENT, "guitar", 26, 2},
				{token.ASSIGN, "=", 33, 2},
				{token.IDENT, "Guitar", 35, 2},
				{token.PASSTHROUGH, "->", 41, 2},
				{token.IDENT, "new", 43, 2},
				{token.IDENT, "tuning", 48, 3},
				{token.ASSIGN, "=", 55, 3},
				{token.IDENT, "tuning", 57, 3},
				{token.PASSTHROUGH, "->", 63, 3},
				{token.IDENT, "toUpper", 65, 3},
				{token.IF, "if", 74, 4},
				{token.IDENT, "tuning", 77, 4},
				{token.PASSTHROUGH, "->", 83, 4},
				{token.NOT, "!", 85, 4},
				{token.IDENT, "inValidTunings", 86, 4},
				{token.LBRACE, "{", 101, 4},
				{token.LPAREN, "(", 105, 5},
				{token.ERROR, "error", 106, 5},
				{token.COLON, ":", 111, 5},
				{token.STRING, "\"{tuning} is not a valid tuning\"", 113, 5},
				{token.RPAREN, ")", 145, 5},
				{token.PASSTHROUGH, "->", 146, 5},
				{token.RBRACE, "}", 150, 6},
				{token.FOR, "for", 153, 7},
				{token.IDENT, "i", 157, 7},
				{token.COMMA, ",", 158, 7},
				{token.IDENT, "t", 160, 7},
				{token.IN, "in", 162, 7},
				{token.IDENT, "array", 165, 7},
				{token.LBRACE, "{", 171, 7},
				{token.IF, "if", 175, 8},
				{token.IDENT, "t", 178, 8},
				{token.PASSTHROUGH, "->", 179, 8},
				{token.IDENT, "len", 181, 8},
				{token.EQUAL, "==", 185, 8},
				{token.INT, "1", 188, 8},
				{token.LBRACE, "{", 190, 8},
				{token.IDENT, "t", 195, 9},
				{token.ASSIGN, "=", 197, 9},
				{token.STRING, "\" \"", 199, 9},
				{token.PLUS, "+", 203, 9},
				{token.IDENT, "t", 205, 9},
				{token.RBRACE, "}", 209, 10},
				{token.IDENT, "guitar", 213, 11},
				{token.STOP, ".", 219, 11},
				{token.IDENT, "Tuning", 220, 11},
				{token.LBRACK, "[", 226, 11},
				{token.IDENT, "i", 227, 11},
				{token.PLUS, "+", 228, 11},
				{token.INT, "1", 229, 11},
				{token.RBRACK, "]", 230, 11},
				{token.ASSIGN, "=", 232, 11},
				{token.IDENT, "t", 234, 11},
				{token.RBRACE, "}", 237, 12},
				{token.LPAREN, "(", 240, 13},
				{token.IDENT, "guitar", 241, 13},
				{token.RPAREN, ")", 247, 13},
				{token.PASSTHROUGH, "->", 248, 13},
				{token.RBRACE, "}", 251, 14},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result []token.Token
			l := NewLexer(tt.in)
			for l.char != 0 {
				result = append(result, l.NextToken())
			}

			if !reflect.DeepEqual(tt.expect, result) {
				t.Errorf(strings.ReplaceAll(fmt.Sprintf("\nExpected: %+v\nGot: %+v", tt.expect, result), "map", ""))
			}
		})
	}
}
