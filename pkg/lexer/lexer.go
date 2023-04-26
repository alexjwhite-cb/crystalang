package lexer

import (
	"unicode"
)

//

// continue ->
// break -/-
// return ()->
//
// meth Add: a, b {
//  (a + b)->
// }

//  meth NewGuitar: tuning {
//  	guitar = Guitar->new
//  	tuning = tuning->toUpper
//  	if tuning->!inValidTunings {
//  	    (error: "{tuning} is not a valid tuning")->
//      }
//  	for i, t in array {
//  	    if t->len == 1 {
//  	        t = " " + t
//  	    }
//  	    guitar.Tuning[i+1] = t
//      }
//  	(guitar)->
//  }

const (
	lBracket    = '('
	rBracket    = ')'
	lSquare     = '['
	rSquare     = ']'
	lBrace      = '{'
	rBrace      = '}'
	lArrow      = '<'
	rArrow      = '>'
	exclamation = '!'
	question    = '?'
	comma       = ','
	stop        = '.'
	plus        = '+'
	hyphen      = '-'
	multiply    = '*'
	divide      = '/'
	and         = '&'
	or          = '|'
	equals      = '='
	colon       = ':'
	quote       = '"'
)

type Lexer struct {
	lexicon map[int]string
	key     int
	val     string
}

func NewLexer() *Lexer {
	return &Lexer{
		lexicon: make(map[int]string),
		key:     0,
		val:     "",
	}
}

// Lex works out a list of IDs, Ops, and Nums from
// the given source code
func (l *Lexer) Lex(code string) (map[int]string, error) {
	for _, r := range code {

		// Handle strings encapsulated inside quotes
		if len(l.val) > 0 && l.val[0] == quote {
			if r == quote && l.val[len(l.val)-1] != '\\' {
				l.val += string(r)
				l.addValThenReset()
				continue
			}
			l.val += string(r)
			continue
		}

		// Handle arg-separating whitespace
		if unicode.IsSpace(r) {
			l.addValThenReset()
			continue
		}

		// Calculate functional tokens
		if unicode.IsSymbol(r) {
			switch r {
			case rArrow:
				l.val += string(r)
				l.addValThenReset()
			//case lArrow:
			//case equals:
			//case plus:
			default:
				l.val += string(r)
			}
		}

		// Calculate functional tokens
		if unicode.IsPunct(r) {
			switch r {
			case hyphen:
				l.addValThenReset()
				l.val += string(r)
			case colon, comma, stop, lBracket, rBracket, lBrace, rBrace, lSquare, rSquare:
				l.addValThenReset()
				l.lexicon[l.key] = string(r)
				l.key++
			default:
				l.val += string(r)
			}
		}

		// Handle numeric values
		if unicode.IsNumber(r) {
			l.val += string(r)
		}

		// Handle string calues
		if unicode.IsLetter(r) {
			l.val += string(r)
		}
	}

	return l.lexicon, nil
}

func (l *Lexer) addValThenReset() {
	if len(l.val) > 0 {
		l.lexicon[l.key] = l.val
		l.key++
		l.val = ""
	}
}
