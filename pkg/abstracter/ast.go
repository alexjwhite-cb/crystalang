package abstracter

import (
	"fmt"
	"github.com/alexjwhite-cb/jet/pkg/token"
)

type Keyword int

const (
	Meth Keyword = iota
	For
	If
	Else
	Describe
	Object
	Const
	Int
	Float
	Bool
	String
	Array
	Map
)

var keywords = map[Keyword]string{
	Meth:     "meth",
	For:      "for",
	If:       "if",
	Else:     "else",
	Describe: "describe",
	Object:   "object",
	Const:    "const",
	Int:      "int",
	Float:    "float",
	Bool:     "bool",
	String:   "string",
	Array:    "array",
	Map:      "map",
}

type (
	Node interface {
		Pos() int // key of opening node
		End() int // key of last character + 1
	}

	// Expr - expression nodes implement the Expr interface.
	Expr interface {
		Node
		exprNode()
	}

	// Stmt - statement nodes implement the Stmt interface.
	Stmt interface {
		Node
		stmtNode()
	}

	// Decl - declaration nodes implement the Decl interface.
	Decl interface {
		Node
		declNode()
	}
)

// Position is the structure containing the Token's start, end, and line values.
// Col should be calculated as Start - Start of the previous NewLine token.
type Position struct {
	Start, End, Col, Line int
}

// ArgList describes arguments that are consumed when declaring Functions,
// Descriptors and
type ArgList struct {
	Open *Position // Position of the colon : that initiates the arg list
	List []Arg
}

type Arg struct {
	Loc  *Position
	Name string
}

// A BasicLit node represents a literal of basic type.
type BasicLit struct {
	Loc   *Position  // literal position
	Kind  token.Type // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
	Value string     // literal string; e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
}

// Ident is used for track the position of a preceding Keyword
type Ident struct {
	IdPos  *Position // position of identifier (meth/describe) or nil
	IdName string
}

type AbstractSyntaxTree struct {
	Declaration []Decl
}

func Abstract(tokens []token.Token) *AbstractSyntaxTree {
	// This will need to call the function that consumes main
	fmt.Printf("%+v", tokens[0])
	return &AbstractSyntaxTree{}
}
