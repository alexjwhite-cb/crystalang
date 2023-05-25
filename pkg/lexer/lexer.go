package lexer

import (
	"github.com/alexjwhite-cb/jet/pkg/token"
	"unicode"
)

// Lexer is a stateful struct that evolves as
// it iterates over the characters in the script/program being imported.
type Lexer struct {
	input        string
	position     int
	readPosition int
	start        int
	end          int
	line         int
	char         rune
}

// New instantiates a new Lexer
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.readChar()
	return l
}

func newToken(tokenType token.TokenType, literal rune, start, line int) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(literal),
		Start:   start,
		Line:    line,
	}
}

// NextToken works out a token map from
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
func (l *Lexer) NextToken() token.Token {
	// Helpful debug
	//fmt.Printf("%v: %s\n", i, string(r))
	var tok token.Token

	l.skipWhitespace()

	switch l.char {
	case '{':
		tok = newToken(token.LBRACE, l.char, l.position, l.line)
	case '}':
		tok = newToken(token.RBRACE, l.char, l.position, l.line)
	case '(':
		tok = newToken(token.LPAREN, l.char, l.position, l.line)
	case ')':
		tok = newToken(token.RPAREN, l.char, l.position, l.line)
	case '[':
		tok = newToken(token.LBRACK, l.char, l.position, l.line)
	case ']':
		tok = newToken(token.RBRACK, l.char, l.position, l.line)
	case ',':
		tok = newToken(token.COMMA, l.char, l.position, l.line)
	case '.':
		tok = newToken(token.STOP, l.char, l.position, l.line)
	case ';':
		tok = newToken(token.SEMICOLON, l.char, l.position, l.line)
	case ':':
		tok = newToken(token.COLON, l.char, l.position, l.line)
	case '?':
		tok = newToken(token.QUESTION, l.char, l.position, l.line)
	case '*':
		tok = newToken(token.MULTIPLY, l.char, l.position, l.line)
	case '/':
		tok = newToken(token.DIVIDE, l.char, l.position, l.line)
	case '\n', '\r':
		tok = newToken(token.NEWLINE, l.char, l.position, l.line)
		l.line++
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:

		switch {
		case unicode.IsLetter(l.char):
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Start = l.start
			tok.Line = l.line
			return tok

		case unicode.IsNumber(l.char):
			tok.Literal = l.readInt()
			tok.Type = token.INT
			tok.Start = l.start
			tok.Line = l.line
			return tok

		case l.isOperator():
			tok.Literal = l.readOperator()
			tok.Type = token.LookupOperator(tok.Literal)
			tok.Start = l.start
			tok.Line = l.line
			return tok

		case l.char == '"':
			tok.Literal = l.readString()
			tok.Type = token.STRING
			tok.Start = l.start
			tok.Line = l.line
			return tok

		default:
			tok = newToken(token.ILLEGAL, l.char, l.position, l.line)
		}

	}
	l.readChar()
	return tok

}

// readChar updates the current character being inspected
// and allows look-ahead.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = rune(l.input[l.readPosition])
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// readIdentifier concisely reads variables, function names, and keywords
func (l *Lexer) readIdentifier() string {
	l.start = l.position
	for unicode.IsLetter(l.char) || unicode.IsDigit(l.char) || l.char == '_' {
		l.readChar()
	}
	return l.input[l.start:l.position]
}

// readInt concisely reads integer values
func (l *Lexer) readInt() string {
	l.start = l.position
	for unicode.IsNumber(l.char) {
		l.readChar()
	}
	return l.input[l.start:l.position]
}

// readInt concisely reads integer values
func (l *Lexer) readString() string {
	l.start = l.position
	l.readChar()
	for l.char != '"' && l.char != 0 {
		if l.char == '\\' && l.input[l.readPosition] == '"' {
			l.readChar()
		}
		l.readChar()
	}
	l.readChar()
	return l.input[l.start:l.position]
}

// IsOperator checks to see if the current string value
// is equivalent to anything in the single-character operator list.
func (l *Lexer) isOperator() bool {
	switch l.char {
	case '+', '-', '=', '<', '>', '|', '&', '!':
		return true
	}
	return false
}

// readOperator concisely reads complex operators
func (l *Lexer) readOperator() string {
	l.start = l.position
	for l.isOperator() {
		l.readChar()
	}
	return l.input[l.start:l.position]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.char) {
		switch l.char {
		case '\n', '\r':
			return
		}
		l.readChar()
	}
}
