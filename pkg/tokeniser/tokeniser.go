package tokeniser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Token string

const (
	Id      Token = "id"
	Op      Token = "op"
	Num     Token = "number"
	Str     Token = "string"
	Exp     Token = "exp"
	Bool    Token = "bool"
	NewLine Token = "newline"
)

// Keywords
const (
	method  = "meth"
	main    = "main"
	forLoop = "for"
)

type Tokeniser struct {
	Tokens map[int]map[Token]interface{}
	key    int
	val    string
}

func NewTokeniser() *Tokeniser {
	return &Tokeniser{
		Tokens: make(map[int]map[Token]interface{}),
		key:    0,
		val:    "",
	}
}

// Tokenise works out a list of IDs, Ops, and Nums from
// the given source code
func (t *Tokeniser) Tokenise(code string) (*Tokeniser, error) {
	for _, r := range code {
		var previous_character rune
		if len(t.val) > 0 {
			previous_character = rune(t.val[len(t.val)-1])
		}

		// Handle strings encapsulated inside quotes
		if len(t.val) > 0 && t.val[0] == '"' {
			if r == '"' && t.val[len(t.val)-1] != '\\' {
				t.val += string(r)
				t.addVal(r)
				continue
			}
			t.val += string(r)
			continue
		}

		switch {
		// Handle arg-separating whitespace
		case unicode.IsSpace(r):
			if t.valIsOperator() {
				t.addRune(rune(t.val[0]))
			}
			if len(t.val) > 0 {
				t.addVal(r)
			}
			switch r {
			case '\r', '\n':
				t.addRune(r)
			}

		case unicode.IsSymbol(r):
			switch r {
			case '>':
				switch len(t.val) {
				case 1:
					if previous_character == '-' {
						t.val += string(r)
						t.addVal(r)
						continue
					}
				}
				t.addVal(r)
				t.val += string(r)
			case '<':
				t.addVal(r)
				t.val += string(r)
			case '+':
				if len(t.val) == 1 && previous_character == '+' {
					t.val += string(r)
					t.addVal(r)
					continue
				}
				t.addVal(r)
				t.val += string(r)
			case '=':
				switch len(t.val) {
				case 1:
					switch previous_character {
					case '+', '-', '=', '>', '<':
						t.val += string(r)
						t.addVal(r)
						continue
					}
				}
				t.addVal(r)
				t.val += string(r)
			default:
				t.val += string(r)
			}

		case unicode.IsPunct(r):
			switch r {
			case '-', '&', '|':
				if len(t.val) == 1 && previous_character == r {
					t.val += string(r)
					t.addVal(r)
					continue
				}
				t.addVal(r)
				t.val += string(r)
			case ':', ',', '.', '(', ')', '{', '}', '[', ']', '!':
				t.addVal(r)
				t.addRune(r)
			default:
				t.val += string(r)
			}

		case unicode.IsNumber(r):
			if t.valIsOperator() {
				t.addRune(rune(t.val[0]))
			} else if len(t.val) > 0 && !unicode.IsNumber(previous_character) && !unicode.IsLetter(previous_character) {
				t.addVal(r)
			}
			t.val += string(r)

		case unicode.IsLetter(r):
			if t.valIsOperator() {
				t.addRune(rune(t.val[0]))
			}
			t.val += string(r)
		default:
			return t, fmt.Errorf("invalid token: %s", string(r))
		}

	}

	return t, nil
}

func (t *Tokeniser) addVal(r rune) {
	t.Tokens[t.key] = make(map[Token]interface{})

	switch value, err := strconv.Atoi(t.val); {
	case err == nil:
		t.Tokens[t.key][Num] = value
	case isMultiCharOperator(t.val):
		t.Tokens[t.key][Op] = t.val
	case isExpression(t.val):
		t.Tokens[t.key][Exp] = t.val
	case isReserved(t.val):
		t.Tokens[t.key][Id] = t.val
	case isBool(t.val):
		t.Tokens[t.key][Bool] = t.val
	case strings.HasPrefix(t.val, "\"") && strings.HasSuffix(t.val, "\""):
		t.Tokens[t.key][Str] = t.val
	default:
		t.Tokens[t.key][Id] = t.val
	}

	if len(t.val) > 0 {
		t.key++
		t.val = ""
	}
}

func (t *Tokeniser) addRune(r rune) {
	t.Tokens[t.key] = make(map[Token]interface{})
	switch {
	case isSingleCharOperator(r):
		t.Tokens[t.key][Op] = r
	case isStatementSeparator(r):
		t.Tokens[t.key][NewLine] = ';'
	}

	t.val = ""
	t.key++
}

func isStatementSeparator(r rune) bool {
	switch r {
	case ';', '\n', '\r':
		return true
	}
	return false
}

func (t *Tokeniser) valIsOperator() bool {
	switch t.val {
	case "(", ")", "[", "]", "{", "}", "<", ">", "!", "?", ",", ".", "+", "-", "*", "/", "=", ":":
		return true
	}
	return false
}

func isSingleCharOperator(r rune) bool {
	switch r {
	case '(', ')', '[', ']', '{', '}', '<', '>', '!', '?', ',', '.', '+', '-', '*', '/', '=', ':':
		return true
	}
	return false
}

func isMultiCharOperator(val string) bool {
	switch val {
	case "->", "-/-", "&&", "||":
		return true
	}
	return false
}

func isExpression(val string) bool {
	switch val {
	case "==", "++", "--", "+=", "-=", "<=", ">=":
		return true
	}
	return false
}

func isReserved(val string) bool {
	switch val {
	case method, main, forLoop:
		return true
	}
	return false
}

func isBool(val string) bool {
	switch val {
	case "true", "false":
		return true
	}
	return false
}
