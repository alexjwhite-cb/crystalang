package token

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//	Identifiers
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"
	BOOL   = "BOOL"

	LBRACE = "{"
	RBRACE = "}"
	LPAREN = "("
	RPAREN = ")"
	LBRACK = "["
	RBRACK = "]"

	//	Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	//	Operators
	PLUS        = "+"
	INCREMENT   = "++"
	MINUS       = "-"
	DECREMENT   = "--"
	ASSIGN      = "="
	EQUAL       = "=="
	NOTEQUAL    = "!="
	PASSTHROUGH = "->"
	EXCLAMATION = "!"
	QUESTION    = "?"

	//	Keywords
	METHOD   = "METHOD"
	FOR      = "FOR"
	IF       = "IF"
	ELSE     = "ELSE"
	DESCRIBE = "DESCRIBE"
	OBJECT   = "OBJECT"
	OVERLOAD = "OVERLOAD"
)

// Token represents the token-types available
// during tokenization.
type Token struct {
	Type  TokenType
	Value any
	Start int
	End   int
	Line  int
}

// TODO: Add ; type for statements that happen on the same line
const (
	Id      = "identifier"
	Op      = "operator"
	Num     = "number"
	Str     = "string"
	Bool    = "bool"
	NewLine = "newline"
)

// Tokenizer is a stateful struct that evolves as
// it iterates over the characters in the script/program being imported.
type Tokenizer struct {
	Tokens []Token
	start  int
	end    int
	val    string
	line   int
}

// NewTokenizer instantiates a new Tokenizer
func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		Tokens: []Token{},
		val:    "",
		line:   1,
	}
}

// newValueToken creates a new Token from Tokenizer values
func (t *Tokenizer) newValueToken(tokenType string, value any, pos int) Token {
	return Token{
		Type:  TokenType(tokenType),
		Value: value,
		Start: pos - len(t.val),
		End:   pos,
		Line:  t.line,
	}
}

// newRuneToken creates a new Token from the current rune
func (t *Tokenizer) newRuneToken(tokenType string, value rune, pos int) Token {
	return Token{
		Type:  TokenType(tokenType),
		Value: value,
		Start: pos,
		End:   pos + 1,
		Line:  t.line,
	}
}

// Tokenize works out a token map from
// the given source code by performing the following steps:
//
//  1. Identify the last character that was added to our value tracker
//  2. If we're currently handling a string, add all characters to the value
//     until we find an unescaped closing quote (")
//  3. Make a decision with whitespace. Check to see if our value is a rune
//     that could be a single-character operator. If not, determine the value's
//     token as normal, then handle the space. If it's a new line character, store it.
//  4. Handle symbols and punctuation
//  5. Handle numbers and letters
//  6. Error if we have a character we're not expecting that won't fit into any
//     of the above patterns
func (t *Tokenizer) Tokenize(code string) (*Tokenizer, error) {
	for i, r := range code {
		// Helpful debug
		//fmt.Printf("%v: %s\n", i, string(r))
		var previousCharacter rune
		if len(t.val) > 0 {
			previousCharacter = rune(t.val[len(t.val)-1])
		}

		// Handle strings encapsulated inside quotes
		if len(t.val) > 0 && t.val[0] == '"' {
			if r == '"' && t.val[len(t.val)-1] != '\\' {
				t.val += string(r)
				t.addVal(i + 1)
				continue
			}
			t.val += string(r)
			continue
		}

		switch {
		// Handle arg-separating whitespace
		case unicode.IsSpace(r):
			if t.valIsOperator() {
				t.addRune(rune(t.val[0]), i-1)
			}
			if len(t.val) > 0 {
				t.addVal(i)
			}
			switch r {
			case '\r', '\n':
				t.addRune(r, i)
			}

		case unicode.IsSymbol(r), unicode.IsPunct(r):
			switch r {
			case '>':
				switch len(t.val) {
				case 1:
					if previousCharacter == '-' {
						t.val += string(r)
						t.addVal(i + 1)
						continue
					}
				}
				t.addVal(i)
				t.val += string(r)
			case '<':
				t.addVal(i)
				t.val += string(r)
			case '=':
				switch len(t.val) {
				case 1:
					switch previousCharacter {
					case '+', '-', '=', '>', '<':
						t.val += string(r)
						t.addVal(i + 1)
						continue
					}
				}
				t.addVal(i)
				t.val += string(r)
			case '-', '&', '|', '+':
				if len(t.val) == 1 && previousCharacter == r {
					t.val += string(r)
					t.addVal(i + 1)
					continue
				}
				t.addVal(i)
				t.val += string(r)
			case ':', ',', '.', '(', ')', '{', '}', '[', ']', '!':
				t.addVal(i)
				t.addRune(r, i)
			default:
				t.val += string(r)
			}

		case unicode.IsNumber(r):
			if t.valIsOperator() {
				t.addRune(rune(t.val[0]), i-1)
			} else if len(t.val) > 0 && !unicode.IsNumber(previousCharacter) && !unicode.IsLetter(previousCharacter) {
				t.addVal(i)
			}
			t.val += string(r)

		case unicode.IsLetter(r):
			if t.valIsOperator() {
				t.addRune(rune(t.val[0]), i-1)
			}
			t.val += string(r)

		default:
			return t, fmt.Errorf("invalid token: %s", string(r))
		}
	}

	return t, nil
}

// addVal adds a value to the token map. This function should
// not be used for adding rune operators. Numbers and
// single character variable or function names are
// compatible.
func (t *Tokenizer) addVal(pos int) {
	switch value, err := strconv.Atoi(t.val); {
	case err == nil:
		t.Tokens = append(t.Tokens, t.newValueToken(Num, value, pos))
	case isMultiCharOperator(t.val):
		t.Tokens = append(t.Tokens, t.newValueToken(Op, t.val, pos))
	case isBool(t.val):
		boolVal, err := strconv.ParseBool(t.val)
		if err != nil {
			panic(err)
		}
		t.Tokens = append(t.Tokens, t.newValueToken(Bool, boolVal, pos))
	case strings.HasPrefix(t.val, "\"") && strings.HasSuffix(t.val, "\""):
		t.Tokens = append(t.Tokens, t.newValueToken(Str, t.val, pos))
	default:
		if len(t.val) > 0 {
			t.Tokens = append(t.Tokens, t.newValueToken(Id, t.val, pos))
		}
	}
	t.val = ""
}

// addRune adds a single rune to the token map
func (t *Tokenizer) addRune(r rune, pos int) {
	switch {
	case isSingleCharOperator(r):
		t.Tokens = append(t.Tokens, t.newRuneToken(Op, r, pos))
	case isStatementSeparator(r):
		t.Tokens = append(t.Tokens, t.newRuneToken(NewLine, ';', pos))
		t.line++
	}

	t.val = ""
}

// isStatementSeparator uses ;, \n, and \r as a way of delimiting
// statements. All are interpreted as semicolons
func isStatementSeparator(r rune) bool {
	switch r {
	case ';', '\n', '\r':
		return true
	}
	return false
}

// valIsOperator checks to see if the current string value
// is equivalent to anything in the single-character operator list.
func (t *Tokenizer) valIsOperator() bool {
	switch t.val {
	case "(", ")", "[", "]", "{", "}", "<", ">", "!", "?", ",", ".", "+", "-", "*", "/", "=", ":":
		return true
	}
	return false
}

// isSingleCharOperator checks to see if the given rune is in the
// single character operator list.
func isSingleCharOperator(r rune) bool {
	switch r {
	case '(', ')', '[', ']', '{', '}', '<', '>', '!', '?', ',', '.', '+', '-', '*', '/', '=', ':':
		return true
	}
	return false
}

// isMultiCharOperator checks to see if the given string value
// is equivalent to anything in the multi-character operator list.
func isMultiCharOperator(val string) bool {
	switch val {
	case "->", "-/-", "&&", "||", "==", "++", "--", "+=", "-=", "<=", ">=":
		return true
	}
	return false
}

// isBool checks to see if the given string value
// is equivalent to anything in the boolean list
func isBool(val string) bool {
	switch val {
	case "true", "false":
		return true
	}
	return false
}
