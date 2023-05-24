package ast

import "github.com/alexjwhite-cb/jet/pkg/token"

type (
	Node interface {
		TokenLiteral() string
	}

	// Expr - expression nodes implement the Expr interface.
	// Expressions produce values
	Expr interface {
		Node
		exprNode()
	}

	// Stmt - statement nodes implement the Stmt interface.
	// Statements do not usually produce values
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

type Program struct {
	Statements []Stmt
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type Ident struct {
	// This is used for token.IDENT tokens
	Token token.Token
	Value string
}

func (i *Ident) exprNode()            {}
func (i *Ident) TokenLiteral() string { return i.Token.Literal }

type ValueStmt struct {
	Token token.Token
	Name  *Ident
	Value Expr
}

func (vs *ValueStmt) stmtNode()            {}
func (vs *ValueStmt) TokenLiteral() string { return vs.Token.Literal }
