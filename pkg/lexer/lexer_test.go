package lexer

import (
	"reflect"
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

func TestLex(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		expect      map[int]string
		expectError bool
	}{
		{
			name: "Entrypoint",
			in:   entry,
			expect: map[int]string{
				0: "meth", 1: "main", 2: "{",
				3: "}",
			},
		},
		{
			name: "Return 1",
			in:   fail,
			expect: map[int]string{
				0: "meth", 1: "main", 2: "{",
				3: "(", 4: "1", 5: ")", 6: "->",
				7: "}",
			},
		},
		{
			name: "Return Var",
			in:   returnVar,
			expect: map[int]string{
				0: "meth", 1: "main", 2: "{",
				3: "num", 4: "=", 5: "0",
				6: "(", 7: "num", 8: ")", 9: "->",
				10: "}",
			},
		},
		{
			name: "Increment Var 1",
			in:   returnIncrement1,
			expect: map[int]string{
				0: "meth", 1: "main", 2: "{",
				3: "num", 4: "=", 5: "0",
				6: "num", 7: "++",
				8: "(", 9: "num", 10: ")", 11: "->",
				12: "}",
			},
		},
		{
			name: "Increment Var 2",
			in:   returnIncrement2,
			expect: map[int]string{
				0: "meth", 1: "main", 2: "{",
				3: "num", 4: "=", 5: "0",
				6: "num", 7: "+=", 8: "2",
				9: "(", 10: "num", 11: ")", 12: "->",
				13: "}",
			},
		},
		{
			name: "Decrement Var 1",
			in:   returnDecrement1,
			expect: map[int]string{
				0: "meth", 1: "main", 2: "{",
				3: "num", 4: "=", 5: "0",
				6: "num", 7: "--",
				8: "(", 9: "num", 10: ")", 11: "->",
				12: "}",
			},
		},
		{
			name: "Decrement Var 2",
			in:   returnDecrement2,
			expect: map[int]string{
				0: "meth", 1: "main", 2: "{",
				3: "num", 4: "=", 5: "0",
				6: "num", 7: "-=", 8: "2",
				9: "(", 10: "num", 11: ")", 12: "->",
				13: "}",
			},
		},
		{
			name: "Return String",
			in:   returnString,
			expect: map[int]string{
				0: "meth", 1: "main", 2: "{",
				3: "str", 4: "=", 5: "\"Hello, World!\"",
				6: "(", 7: "str", 8: ")", 9: "->",
				10: "}",
			},
		},
		{
			name: "Return EscapeString",
			in:   returnEscapeString,
			expect: map[int]string{
				0: "meth", 1: "main", 2: "{",
				3: "str", 4: "=", 5: "\"\\\"Hello, World!\\\"\"",
				6: "(", 7: "str", 8: ")", 9: "->",
				10: "}",
			},
		},
		{
			name: "Simple Function",
			in:   func1,
			expect: map[int]string{
				0: "meth", 1: "NewGuitar", 2: ":", 3: "tuning", 4: "{",
				5: "guitar", 6: "=", 7: "Guitar", 8: "->", 9: "new",
				10: "tuning", 11: "=", 12: "tuning", 13: "->", 14: "toUpper",
				15: "if", 16: "tuning", 17: "->", 18: "!", 19: "inValidTunings", 20: "{",
				21: "(", 22: "error", 23: ":", 24: "\"{tuning} is not a valid tuning\"", 25: ")", 26: "->",
				27: "}",
				28: "for", 29: "i", 30: ",", 31: "t", 32: "in", 33: "array", 34: "{",
				35: "if", 36: "t", 37: "->", 38: "len", 39: "==", 40: "1", 41: "{",
				42: "t", 43: "=", 44: "\" \"", 45: "+", 46: "t",
				47: "}",
				48: "guitar", 49: ".", 50: "Tuning", 51: "[", 52: "i", 53: "+", 54: "1", 55: "]", 56: "=", 57: "t",
				58: "}",
				59: "(", 60: "guitar", 61: ")", 62: "->",
				63: "}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewLexer().Lex(tt.in)
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

			if !reflect.DeepEqual(tt.expect, result) {
				t.Errorf("\nExpected: %+v\nGot: %+v", tt.expect, result)
			}
		})
	}
}
