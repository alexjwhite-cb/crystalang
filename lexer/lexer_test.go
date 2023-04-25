package lexer

import (
	"reflect"
	"testing"
)

const (
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
		expect      string
		expectError bool
	}{
		{
			name:   "Simple Function",
			in:     func1,
			expect: "",
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
				t.Errorf("Expected: %s\nGot: %+v", tt.expect, result)
			}
		})
	}
}
