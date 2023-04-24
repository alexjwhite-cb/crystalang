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
	l_bracket   = '('
	r_bracket   = ')'
	l_sqBracket = '['
	r_sqBracket = ']'
	l_brace     = '{'
	r_brace     = '}'
	exclamation = '!'
	question    = '?'
	plus        = '+'
	hyphen      = '-'
	multiply    = '*'
	divide      = '/'
	and         = '&'
	or          = '|'
	l_arrow     = '<'
	r_arrow     = '>'
	equals      = '='
	colon       = ':'
	quote       = '"'
)

// Lex works out a list of IDs, Ops, and Nums from
// the given source code
func Lex(code string) (map[int]string, error) {
	lexicon := make(map[int]string)
	key := 0
	var val string

	for _, r := range []rune(code) {

		if unicode.IsSpace(r) {
			if len(val) > 0 {
				lexicon[key] = val
				key++
				val = ""
			}
			continue
		}

		if unicode.IsPunct(r) {
			switch r {
			case l_bracket:
				fallthrough
			case r_bracket:
				fallthrough
			case l_brace:
				fallthrough
			case r_brace:
				fallthrough
			case l_sqBracket:
				fallthrough
			case r_sqBracket:
				if len(val) > 0 {
					lexicon[key] = val
					key++
					val = ""
				}
				lexicon[key] = string(r)
				key++
			default:
				val += string(r)
			}
		}

		if unicode.IsNumber(r) {
			val += string(r)
		}

		if unicode.IsLetter(r) {
			val += string(r)
		}
	}

	return lexicon, nil
}
