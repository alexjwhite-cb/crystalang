package parser

import (
	"fmt"
	"github.com/alexjwhite-cb/jet/pkg/ast"
	"github.com/alexjwhite-cb/jet/pkg/lexer"
	"github.com/alexjwhite-cb/jet/pkg/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS   // == or !=
	LESSMORE // < or >
	SUM      // + or -
	PRODUCT  // * or /
	PREFIX   // -x or !x
	CALL     // myFunction()
)

var priority = map[token.TokenType]int{
	token.EQUAL:       EQUALS,
	token.NOTEQUAL:    EQUALS,
	token.MOREOREQUAL: EQUALS,
	token.LESSOREQUAL: EQUALS,
	token.LESSTHAN:    LESSMORE,
	token.MORETHAN:    LESSMORE,
	token.PLUS:        SUM,
	token.MINUS:       SUM,
	token.MULTIPLY:    PRODUCT,
	token.DIVIDE:      PRODUCT,
}

type (
	prefixParseFn func() ast.Expr
	infixParseFn  func(ast.Expr) ast.Expr
)

func (p *Parser) registerPrefix(tType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tType] = fn
}
func (p *Parser) registerInfix(tType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tType] = fn
}

type Parser struct {
	l              *lexer.Lexer
	curToken       token.Token
	peekToken      token.Token
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
	errors         []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentity)
	p.registerPrefix(token.INT, p.parseIntLiteral)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.EQUAL, p.parseInfixExpression)
	p.registerInfix(token.NOTEQUAL, p.parseInfixExpression)
	p.registerInfix(token.MOREOREQUAL, p.parseInfixExpression)
	p.registerInfix(token.LESSOREQUAL, p.parseInfixExpression)
	p.registerInfix(token.LESSTHAN, p.parseInfixExpression)
	p.registerInfix(token.MORETHAN, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.MULTIPLY, p.parseInfixExpression)
	p.registerInfix(token.DIVIDE, p.parseInfixExpression)

	//	 Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Stmt{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Stmt {
	switch p.curToken.Type {
	case token.IDENT:
		if p.peekTokenIs(token.ASSIGN) {
			return p.parseValueStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseValueStatement() *ast.ValueStmt {
	stmt := &ast.ValueStmt{}
	stmt.Name = &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: Handle Expressions
	for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.NEWLINE) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStmt {
	stmt := &ast.ExpressionStmt{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseIntLiteral() ast.Expr {
	lit := &ast.IntLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("line %v, col %v: could not parse %q as integer",
			p.curToken.Line, p.curToken.Start, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) noPrefixParseFnError(t token.Token) {
	errToken := fmt.Sprintf("%s", t.Type)
	if t.Type == token.ILLEGAL {
		errToken = fmt.Sprintf("%s (%s)", t.Type, t.Literal)
	}
	msg := fmt.Sprintf("line %v, col %v: no prefix parse function for %s found", t.Line, t.Start, errToken)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPriority() int {
	if prio, ok := priority[p.peekToken.Type]; ok {
		return prio
	}
	return LOWEST
}

func (p *Parser) curPriority() int {
	if prio, ok := priority[p.curToken.Type]; ok {
		return prio
	}
	return LOWEST
}

func (p *Parser) parseExpression(prio int) ast.Expr {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) &&
		!p.peekTokenIs(token.NEWLINE) &&
		!p.peekTokenIs(token.EOF) &&
		prio < p.peekPriority() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expr {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expr) ast.Expr {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	prio := p.curPriority()
	p.nextToken()
	expression.Right = p.parseExpression(prio)
	return expression
}

func (p *Parser) parseIdentity() ast.Expr {
	return &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t, p.peekToken)
	return false
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(expect token.TokenType, t token.Token) {
	peek := p.peekToken.Type
	if peek == token.NEWLINE {
		peek = "newline"
	}
	msg := fmt.Sprintf("line %v, col %v: expected %s, got %s", t.Line, t.Start, expect, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
