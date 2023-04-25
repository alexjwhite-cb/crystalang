package lexer

import "unicode"

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

		if unicode.IsSpace(r) {
			l.addValThenReset()
			continue
		}

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

		if unicode.IsPunct(r) {
			switch r {
			case hyphen:
				l.addValThenReset()
				l.val += string(r)
			case colon:
				fallthrough
			case comma:
				fallthrough
			case stop:
				fallthrough
			case lBracket:
				fallthrough
			case rBracket:
				fallthrough
			case lBrace:
				fallthrough
			case rBrace:
				fallthrough
			case lSquare:
				fallthrough
			case rSquare:
				l.addValThenReset()
				l.lexicon[l.key] = string(r)
				l.key++
			default:
				l.val += string(r)
			}
		}

		if unicode.IsNumber(r) {
			l.val += string(r)
		}

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
