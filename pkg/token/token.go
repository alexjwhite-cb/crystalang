package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//	Identifiers
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	LBRACE = "{"
	RBRACE = "}"
	LPAREN = "("
	RPAREN = ")"
	LBRACK = "["
	RBRACK = "]"

	//	Delimiters
	STOP      = "."
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	//	Operators
	DIVIDE      = "/"
	MULTIPLY    = "*"
	PLUS        = "+"
	INCREMENT   = "++"
	MINUS       = "-"
	DECREMENT   = "--"
	ASSIGN      = "="
	MINUSASSIGN = "-="
	PLUSASSIGN  = "+="
	EQUAL       = "=="
	NOTEQUAL    = "!="
	PASSTHROUGH = "->"
	NOT         = "!"
	QUESTION    = "?"
	LESSTHAN    = "<"
	LESSOREQUAL = "<="
	MORETHAN    = ">"
	MOREOREQUAL = ">="
	AND         = "&&"
	OR          = "||"

	//	Keywords
	METHOD   = "METHOD"
	FOR      = "FOR"
	IF       = "IF"
	ELSE     = "ELSE"
	DESCRIBE = "DESCRIBE"
	OBJECT   = "OBJECT"
	OVERLOAD = "OVERLOAD"
	IN       = "in"
	ERROR    = "error"
	TRUE     = "true"
	FALSE    = "false"
)

var keywords = map[string]TokenType{
	"meth":     METHOD,
	"for":      FOR,
	"if":       IF,
	"else":     ELSE,
	"describe": DESCRIBE,
	"object":   OBJECT,
	"overload": OVERLOAD,
	"in":       IN,
	"error":    ERROR,
	"true":     TRUE,
	"false":    FALSE,
}

var operators = map[string]TokenType{
	"*":  MULTIPLY,
	"/":  DIVIDE,
	"+":  PLUS,
	"++": INCREMENT,
	"-":  MINUS,
	"--": DECREMENT,
	"=":  ASSIGN,
	"-=": MINUSASSIGN,
	"+=": PLUSASSIGN,
	"==": EQUAL,
	"!":  NOT,
	"!=": NOTEQUAL,
	"->": PASSTHROUGH,
	"<":  LESSTHAN,
	"<=": LESSOREQUAL,
	">":  MORETHAN,
	">=": MOREOREQUAL,
	"&&": AND,
	"||": OR,
}

// Token represents the token-types available
// during tokenization.
type Token struct {
	Type    TokenType
	Literal string
	Start   int
	Line    int
}

// LookupIdent checks if an identity is a reserved word
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// LookupOperator checks if a string of operators is valid
func LookupOperator(op string) TokenType {
	if tok, ok := operators[op]; ok {
		return tok
	}
	return ILLEGAL
}
