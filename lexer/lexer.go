package lexer

// meth Add: a, b {
//  return a + b
// }

// meth NewGuitar: tuning {
// 	guitar = new: Guitar
// 	tuning = toUpper: tuning
// 	if !validTuningRgx.MatchString: tuning {
// 		return error: "{tuning} is not a valid tuning"
// 	}
// 	for i, t in validTuningRgx.FindAllString(tuning, -1) {
// 		if len: t == 1 {
// 			t = " " + t
// 		}
// 		guitar.Tuning[i+1] = t
// 	}
// 	return guitar
// }

import (
	"regexp"
)

var (
	split = regexp.MustCompile(`\r+|\n+| +`)
	// validVar = regexp.MustCompile(`([0-9]?[A-Z])\w+`)
)

// Lex works out a list of IDs, Ops, and Nums from
// the given source code
func Lex(code string) ([]string, error) {
	for _, word := range split.Split(code, -1) {
		println(word)
		// Contains

		//
		// Contains New Line
	}
}
