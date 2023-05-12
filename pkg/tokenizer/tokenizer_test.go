package tokenizer

import (
	"fmt"
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

func TestTokenise(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		expect      []Token
		expectError bool
	}{
		{
			name: "Entrypoint",
			in:   entry,
			expect: []Token{
				{Id, "meth", 0, 4, 1},
				{Id, "main", 5, 9, 1},
				{Op, '{', 10, 11, 1},
				{NewLine, ';', 11, 12, 1},
				{Op, '}', 12, 13, 2},
			},
		},
		{
			name: "Return 1",
			in:   fail,
			expect: []Token{
				{Id, "meth", 0, 4, 1},
				{Id, "main", 5, 9, 1},
				{Op, '{', 10, 11, 1},
				{NewLine, ';', 11, 12, 1},
				{Op, '(', 13, 14, 2},
				{Num, 1, 14, 15, 2},
				{Op, ')', 15, 16, 2},
				{Op, "->", 16, 18, 2},
				{NewLine, ';', 18, 19, 2},
				{Op, '}', 19, 20, 3},
			},
		},
		{
			name: "Return Var",
			in:   returnVar,
			expect: []Token{
				{Id, "meth", 0, 4, 1},
				{Id, "main", 5, 9, 1},
				{Op, '{', 10, 11, 1},
				{NewLine, ';', 11, 12, 1},
				{Id, "num", 13, 16, 2},
				{Op, '=', 17, 18, 2},
				{Num, 0, 19, 20, 2},
				{NewLine, ';', 20, 21, 2},
				{Op, '(', 22, 23, 3},
				{Id, "num", 23, 26, 3},
				{Op, ')', 26, 27, 3},
				{Op, "->", 27, 29, 3},
				{NewLine, ';', 29, 30, 3},
				{Op, '}', 30, 31, 4},
			},
		},
		{
			name: "Increment Var 1",
			in:   returnIncrement1,
			expect: []Token{
				{Id, "meth", 0, 4, 1},
				{Id, "main", 5, 9, 1},
				{Op, '{', 10, 11, 1},
				{NewLine, ';', 11, 12, 1},
				{Id, "num", 13, 16, 2},
				{Op, '=', 17, 18, 2},
				{Num, 0, 19, 20, 2},
				{NewLine, ';', 20, 21, 2},
				{Id, "num", 22, 25, 3},
				{Op, "++", 25, 27, 3},
				{NewLine, ';', 27, 28, 3},
				{Op, '(', 29, 30, 4},
				{Id, "num", 30, 33, 4},
				{Op, ')', 33, 34, 4},
				{Op, "->", 34, 36, 4},
				{NewLine, ';', 36, 37, 4},
				{Op, '}', 37, 38, 5},
			},
		},
		{
			name: "Increment Var 2",
			in:   returnIncrement2,
			expect: []Token{
				{Id, "meth", 0, 4, 1},
				{Id, "main", 5, 9, 1},
				{Op, '{', 10, 11, 1},
				{NewLine, ';', 11, 12, 1},
				{Id, "num", 13, 16, 2},
				{Op, '=', 17, 18, 2},
				{Num, 0, 19, 20, 2},
				{NewLine, ';', 20, 21, 2},
				{Id, "num", 22, 25, 3},
				{Op, "+=", 26, 28, 3},
				{Num, 2, 29, 30, 3},
				{NewLine, ';', 30, 31, 3},
				{Op, '(', 32, 33, 4},
				{Id, "num", 33, 36, 4},
				{Op, ')', 36, 37, 4},
				{Op, "->", 37, 39, 4},
				{NewLine, ';', 39, 40, 4},
				{Op, '}', 40, 41, 5},
			},
		},
		{
			name: "Decrement Var 1",
			in:   returnDecrement1,
			expect: []Token{
				{Id, "meth", 0, 4, 1},
				{Id, "main", 5, 9, 1},
				{Op, '{', 10, 11, 1},
				{NewLine, ';', 11, 12, 1},
				{Id, "num", 13, 16, 2},
				{Op, '=', 17, 18, 2},
				{Num, 0, 19, 20, 2},
				{NewLine, ';', 20, 21, 2},
				{Id, "num", 22, 25, 3},
				{Op, "--", 25, 27, 3},
				{NewLine, ';', 27, 28, 3},
				{Op, '(', 29, 30, 4},
				{Id, "num", 30, 33, 4},
				{Op, ')', 33, 34, 4},
				{Op, "->", 34, 36, 4},
				{NewLine, ';', 36, 37, 4},
				{Op, '}', 37, 38, 5},
			},
		},
		{
			name: "Decrement Var 2",
			in:   returnDecrement2,
			expect: []Token{
				{Id, "meth", 0, 4, 1},
				{Id, "main", 5, 9, 1},
				{Op, '{', 10, 11, 1},
				{NewLine, ';', 11, 12, 1},
				{Id, "num", 13, 16, 2},
				{Op, '=', 17, 18, 2},
				{Num, 0, 19, 20, 2},
				{NewLine, ';', 20, 21, 2},
				{Id, "num", 22, 25, 3},
				{Op, "-=", 26, 28, 3},
				{Num, 2, 29, 30, 3},
				{NewLine, ';', 30, 31, 3},
				{Op, '(', 32, 33, 4},
				{Id, "num", 33, 36, 4},
				{Op, ')', 36, 37, 4},
				{Op, "->", 37, 39, 4},
				{NewLine, ';', 39, 40, 4},
				{Op, '}', 40, 41, 5},
			},
		},
		{
			name: "Return String",
			in:   returnString,
			expect: []Token{
				{Id, "meth", 0, 4, 1},
				{Id, "main", 5, 9, 1},
				{Op, '{', 10, 11, 1},
				{NewLine, ';', 11, 12, 1},
				{Id, "str", 13, 16, 2},
				{Op, '=', 17, 18, 2},
				{Str, "\"Hello, World!\"", 19, 34, 2},
				{NewLine, ';', 34, 35, 2},
				{Op, '(', 36, 37, 3},
				{Id, "str", 37, 40, 3},
				{Op, ')', 40, 41, 3},
				{Op, "->", 41, 43, 3},
				{NewLine, ';', 43, 44, 3},
				{Op, '}', 44, 45, 4},
			},
		},
		{
			name: "Return EscapeString",
			in:   returnEscapeString,
			expect: []Token{
				{Id, "meth", 0, 4, 1},
				{Id, "main", 5, 9, 1},
				{Op, '{', 10, 11, 1},
				{NewLine, ';', 11, 12, 1},
				{Id, "str", 13, 16, 2},
				{Op, '=', 17, 18, 2},
				{Str, "\"\\\"Hello, World!\\\"\"", 19, 38, 2},
				{NewLine, ';', 38, 39, 2},
				{Op, '(', 40, 41, 3},
				{Id, "str", 41, 44, 3},
				{Op, ')', 44, 45, 3},
				{Op, "->", 45, 47, 3},
				{NewLine, ';', 47, 48, 3},
				{Op, '}', 48, 49, 4},
			},
		},
		{
			name: "And Operator",
			in:   andOperator,
			expect: []Token{
				{Id, "meth", 0, 4, 1},
				{Id, "main", 5, 9, 1},
				{Op, '{', 10, 11, 1},
				{NewLine, ';', 11, 12, 1},
				{Id, "if", 13, 15, 2},
				{Id, "a", 16, 17, 2},
				{Op, "==", 18, 20, 2},
				{Num, 2, 21, 22, 2},
				{Op, '*', 23, 24, 2},
				{Num, 2, 25, 26, 2},
				{Op, "&&", 27, 29, 2},
				{Op, '!', 30, 31, 2},
				{Id, "b", 31, 32, 2},
				{Op, '{', 33, 34, 2},
				{NewLine, ';', 34, 35, 2},
				{Op, '(', 37, 38, 3},
				{Bool, true, 38, 42, 3},
				{Op, ')', 42, 43, 3},
				{Op, "->", 43, 45, 3},
				{NewLine, ';', 45, 46, 3},
				{Op, '}', 47, 48, 4},
				{Id, "else", 49, 53, 4},
				{Op, '{', 54, 55, 4},
				{NewLine, ';', 55, 56, 4},
				{Op, '(', 58, 59, 5},
				{Bool, false, 59, 64, 5},
				{Op, ')', 64, 65, 5},
				{Op, "->", 65, 67, 5},
				{NewLine, ';', 67, 68, 5},
				{Op, '}', 69, 70, 6},
				{NewLine, ';', 70, 71, 6},
				{Op, '}', 71, 72, 7},
			},
		},
		{
			name: "Or Operator",
			in:   orOperator,
			expect: []Token{
				{Id, "meth", 0, 4, 1},
				{Id, "main", 5, 9, 1},
				{Op, '{', 10, 11, 1},
				{NewLine, ';', 11, 12, 1},
				{Id, "if", 13, 15, 2},
				{Id, "a", 16, 17, 2},
				{Op, "==", 18, 20, 2},
				{Num, 2, 21, 22, 2},
				{Op, '*', 23, 24, 2},
				{Num, 2, 25, 26, 2},
				{Op, "||", 27, 29, 2},
				{Op, '!', 30, 31, 2},
				{Id, "b", 31, 32, 2},
				{Op, '{', 33, 34, 2},
				{NewLine, ';', 34, 35, 2},
				{Op, '(', 37, 38, 3},
				{Bool, true, 38, 42, 3},
				{Op, ')', 42, 43, 3},
				{Op, "->", 43, 45, 3},
				{NewLine, ';', 45, 46, 3},
				{Op, '}', 47, 48, 4},
				{NewLine, ';', 48, 49, 4},
				{Op, '(', 50, 51, 5},
				{Bool, false, 51, 56, 5},
				{Op, ')', 56, 57, 5},
				{Op, "->", 57, 59, 5},
				{NewLine, ';', 59, 60, 5},
				{Op, '}', 60, 61, 6},
			},
		},
		{
			name: "Simple Function",
			in:   func1,
			expect: []Token{
				{Id, "meth", 0, 4, 1},
				{Id, "NewGuitar", 5, 14, 1},
				{Op, ':', 14, 15, 1},
				{Id, "tuning", 16, 22, 1},
				{Op, '{', 23, 24, 1},
				{NewLine, ';', 24, 25, 1},
				{Id, "guitar", 26, 32, 2},
				{Op, '=', 33, 34, 2},
				{Id, "Guitar", 35, 41, 2},
				{Op, "->", 41, 43, 2},
				{Id, "new", 43, 46, 2},
				{NewLine, ';', 46, 47, 2},
				{Id, "tuning", 48, 54, 3},
				{Op, '=', 55, 56, 3},
				{Id, "tuning", 57, 63, 3},
				{Op, "->", 63, 65, 3},
				{Id, "toUpper", 65, 72, 3},
				{NewLine, ';', 72, 73, 3},
				{Id, "if", 74, 76, 4},
				{Id, "tuning", 77, 83, 4},
				{Op, "->", 83, 85, 4},
				{Op, '!', 85, 86, 4},
				{Id, "inValidTunings", 86, 100, 4},
				{Op, '{', 101, 102, 4},
				{NewLine, ';', 102, 103, 4},
				{Op, '(', 105, 106, 5},
				{Id, "error", 106, 111, 5},
				{Op, ':', 111, 112, 5},
				{Str, "\"{tuning} is not a valid tuning\"", 113, 145, 5},
				{Op, ')', 145, 146, 5},
				{Op, "->", 146, 148, 5},
				{NewLine, ';', 148, 149, 5},
				{Op, '}', 150, 151, 6},
				{NewLine, ';', 151, 152, 6},
				{Id, "for", 153, 156, 7},
				{Id, "i", 157, 158, 7},
				{Op, ',', 158, 159, 7},
				{Id, "t", 160, 161, 7},
				{Id, "in", 162, 164, 7},
				{Id, "array", 165, 170, 7},
				{Op, '{', 171, 172, 7},
				{NewLine, ';', 172, 173, 7},
				{Id, "if", 175, 177, 8},
				{Id, "t", 178, 179, 8},
				{Op, "->", 179, 181, 8},
				{Id, "len", 181, 184, 8},
				{Op, "==", 185, 187, 8},
				{Num, 1, 188, 189, 8},
				{Op, '{', 190, 191, 8},
				{NewLine, ';', 191, 192, 8},
				{Id, "t", 195, 196, 9},
				{Op, '=', 197, 198, 9},
				{Str, "\" \"", 199, 202, 9},
				{Op, '+', 203, 204, 9},
				{Id, "t", 205, 206, 9},
				{NewLine, ';', 206, 207, 9},
				{Op, '}', 209, 210, 10},
				{NewLine, ';', 210, 211, 10},
				{Id, "guitar", 213, 219, 11},
				{Op, '.', 219, 220, 11},
				{Id, "Tuning", 220, 226, 11},
				{Op, '[', 226, 227, 11},
				{Id, "i", 227, 228, 11},
				{Op, '+', 228, 229, 11},
				{Num, 1, 229, 230, 11},
				{Op, ']', 230, 231, 11},
				{Op, '=', 232, 233, 11},
				{Id, "t", 234, 235, 11},
				{NewLine, ';', 235, 236, 11},
				{Op, '}', 237, 238, 12},
				{NewLine, ';', 238, 239, 12},
				{Op, '(', 240, 241, 13},
				{Id, "guitar", 241, 247, 13},
				{Op, ')', 247, 248, 13},
				{Op, "->", 248, 250, 13},
				{NewLine, ';', 250, 251, 13},
				{Op, '}', 251, 252, 14},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewTokenizer().Tokenize(tt.in)
			switch {
			case tt.expectError && err != nil:
				return
			case tt.expectError && err == nil:
				t.Errorf("Expected error, got none.")
				return
			case err != nil:
				t.Errorf("Unexpected error: %s", err)
				return
			}

			if !reflect.DeepEqual(tt.expect, result.Tokens) {
				t.Errorf(strings.ReplaceAll(fmt.Sprintf("\nExpected: %+v\nGot: %+v", tt.expect, result.Tokens), "map", ""))
			}
		})
	}
}
