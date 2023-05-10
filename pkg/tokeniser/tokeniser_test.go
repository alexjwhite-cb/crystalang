package tokeniser

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
			name:   "And Operator",
			in:     andOperator,
			expect: []Token{},
			//expect: map[int]map[Token]interface{}{
			//	0: {Id: "meth"}, 1: {Id: "main"}, 2: {Op: '{'}, 3: {NewLine: ';'},
			//	4: {Id: "if"}, 5: {Id: "a"}, 6: {Op: "=="}, 7: {Num: 2}, 8: {Op: '*'}, 9: {Num: 2}, 10: {Op: "&&"}, 11: {Op: '!'}, 12: {Id: "b"}, 13: {Op: '{'}, 14: {NewLine: ';'},
			//	15: {Op: '('}, 16: {Bool: "true"}, 17: {Op: ')'}, 18: {Op: "->"}, 19: {NewLine: ';'},
			//	20: {Op: '}'}, 21: {Id: "else"}, 22: {Op: '{'}, 23: {NewLine: ';'},
			//	24: {Op: '('}, 25: {Bool: "false"}, 26: {Op: ')'}, 27: {Op: "->"}, 28: {NewLine: ';'},
			//	29: {Op: '}'}, 30: {NewLine: ';'},
			//	31: {Op: '}'},
			//},
		},
		{
			name:   "Or Operator",
			in:     orOperator,
			expect: []Token{},
			//expect: map[int]map[Token]interface{}{
			//	0: {Id: "meth"}, 1: {Id: "main"}, 2: {Op: '{'}, 3: {NewLine: ';'},
			//	4: {Id: "if"}, 5: {Id: "a"}, 6: {Op: "=="}, 7: {Num: 2}, 8: {Op: '*'}, 9: {Num: 2}, 10: {Op: "||"}, 11: {Op: '!'}, 12: {Id: "b"}, 13: {Op: '{'}, 14: {NewLine: ';'},
			//	15: {Op: '('}, 16: {Bool: "true"}, 17: {Op: ')'}, 18: {Op: "->"}, 19: {NewLine: ';'},
			//	20: {Op: '}'}, 21: {NewLine: ';'},
			//	22: {Op: '('}, 23: {Bool: "false"}, 24: {Op: ')'}, 25: {Op: "->"}, 26: {NewLine: ';'},
			//	27: {Op: '}'},
			//},
		},
		{
			name:   "Simple Function",
			in:     func1,
			expect: []Token{},
			//expect: map[int]map[Token]interface{}{
			//	0: {Id: "meth"}, 1: {Id: "NewGuitar"}, 2: {Op: ':'}, 3: {Id: "tuning"}, 4: {Op: '{'}, 5: {NewLine: ';'},
			//	6: {Id: "guitar"}, 7: {Op: '='}, 8: {Id: "Guitar"}, 9: {Op: "->"}, 10: {Id: "new"}, 11: {NewLine: ';'},
			//	12: {Id: "tuning"}, 13: {Op: '='}, 14: {Id: "tuning"}, 15: {Op: "->"}, 16: {Id: "toUpper"}, 17: {NewLine: ';'},
			//	18: {Id: "if"}, 19: {Id: "tuning"}, 20: {Op: "->"}, 21: {Op: '!'}, 22: {Id: "inValidTunings"}, 23: {Op: '{'}, 24: {NewLine: ';'},
			//	25: {Op: '('}, 26: {Id: "error"}, 27: {Op: ':'}, 28: {Str: "\"{tuning} is not a valid tuning\""}, 29: {Op: ')'}, 30: {Op: "->"}, 31: {NewLine: ';'},
			//	32: {Op: '}'}, 33: {NewLine: ';'},
			//	34: {Id: "for"}, 35: {Id: "i"}, 36: {Op: ','}, 37: {Id: "t"}, 38: {Id: "in"}, 39: {Id: "array"}, 40: {Op: '{'}, 41: {NewLine: ';'},
			//	42: {Id: "if"}, 43: {Id: "t"}, 44: {Op: "->"}, 45: {Id: "len"}, 46: {Op: "=="}, 47: {Num: 1}, 48: {Op: '{'}, 49: {NewLine: ';'},
			//	50: {Id: "t"}, 51: {Op: '='}, 52: {Str: "\" \""}, 53: {Op: '+'}, 54: {Id: "t"}, 55: {NewLine: ';'},
			//	56: {Op: '}'}, 57: {NewLine: ';'},
			//	58: {Id: "guitar"}, 59: {Op: '.'}, 60: {Id: "Tuning"}, 61: {Op: '['}, 62: {Id: "i"}, 63: {Op: '+'}, 64: {Num: 1}, 65: {Op: ']'}, 66: {Op: '='}, 67: {Id: "t"}, 68: {NewLine: ';'},
			//	69: {Op: '}'}, 70: {NewLine: ';'},
			//	71: {Op: '('}, 72: {Id: "guitar"}, 73: {Op: ')'}, 74: {Op: "->"}, 75: {NewLine: ';'},
			//	76: {Op: '}'},
			//},
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
