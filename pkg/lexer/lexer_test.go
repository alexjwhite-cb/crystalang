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
	returnString = `meth main {
	str = "Hello, World!"
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
			name:   "Simple Function",
			in:     func1,
			expect: map[int]string{},
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
