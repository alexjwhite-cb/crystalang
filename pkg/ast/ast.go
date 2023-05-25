package ast

import (
	"bytes"
	"fmt"
	"github.com/alexjwhite-cb/jet/pkg/token"
)

type (
	Node interface {
		TokenLiteral() string
		fmt.Stringer
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

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type Ident struct {
	// This is used for token.IDENT tokens
	Token token.Token
	Value string
}

func (i *Ident) exprNode()            {}
func (i *Ident) TokenLiteral() string { return i.Token.Literal }
func (i *Ident) String() string       { return i.Value }

type IntLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntLiteral) exprNode()            {}
func (i *IntLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntLiteral) String() string       { return i.Token.Literal }

type ValueStmt struct {
	Token token.Token
	Name  *Ident
	Value Expr
}

func (vs *ValueStmt) stmtNode()            {}
func (vs *ValueStmt) TokenLiteral() string { return vs.Token.Literal }
func (vs *ValueStmt) String() string {
	var out bytes.Buffer
	out.WriteString(vs.Name.String())
	out.WriteString(" = ")
	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}
	return out.String()
}

type ExpressionStmt struct {
	Token      token.Token
	Expression Expr
}

func (es *ExpressionStmt) stmtNode()            {}
func (es *ExpressionStmt) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStmt) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
