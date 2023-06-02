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
	POSTFIX  // ->
	EQUALS   // == or !=
	LESSMORE // < or >
	SUM      // + or -
	PRODUCT  // * or /
	PREFIX   // -x or !x
	CALL     // myFunction()
	INDEX    // array[index]
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
	token.LPAREN:      CALL,
	token.PASSTHROUGH: POSTFIX,
	token.LBRACK:      INDEX,
}

type (
	prefixParseFn  func() ast.Expr
	infixParseFn   func(ast.Expr) ast.Expr
	postfixParseFn func(ast.Stmt) ast.Stmt
)

func (p *Parser) registerPrefix(tType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tType] = fn
}
func (p *Parser) registerInfix(tType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tType] = fn
}
func (p *Parser) registerPostfix(tType token.TokenType, fn postfixParseFn) {
	p.postfixParseFns[tType] = fn
}

type Parser struct {
	l               *lexer.Lexer
	curToken        token.Token
	peekToken       token.Token
	prefixParseFns  map[token.TokenType]prefixParseFn
	infixParseFns   map[token.TokenType]infixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn
	errors          []string
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
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.METHOD, p.parseFuncLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACK, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashMap)
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
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACK, p.parseIndexExpression)
	p.postfixParseFns = make(map[token.TokenType]postfixParseFn)
	p.registerPostfix(token.PASSTHROUGH, p.parseReturnStatement)

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
	var exp ast.Stmt
	switch p.curToken.Type {
	case token.IDENT:
		if p.peekTokenIs(token.ASSIGN) {
			return p.parseValueStatement()
		}
		exp = p.parseExpressionStatement()
	case token.METHOD, token.DESCRIBE, token.OBJECT:
		if p.peekTokenIs(token.IDENT) {
			return p.parseDeclarationStatement()
		}
	default:
		exp = p.parseExpressionStatement()
	}
	exp = p.parsePostfixExpressionStatement(exp)
	return exp
}

//
// This section is all about returning statements
//

func (p *Parser) parseValueStatement() *ast.ValueStmt {
	stmt := &ast.ValueStmt{Token: p.curToken}
	stmt.Name = &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.NEWLINE) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement(left ast.Stmt) ast.Stmt {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	expr, ok := left.(*ast.ExpressionStmt)
	if !ok {
		msg := fmt.Sprintf("line %v, col %v: illegal return", p.curToken.Line, p.curToken.Col)
		p.errors = append(p.errors, msg)
		return nil
	}
	stmt.Value = expr.Expression
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStmt {
	stmt := &ast.ExpressionStmt{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.PASSTHROUGH) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseDeclarationStatement() ast.Stmt {
	decl := &ast.DeclarationStmt{Token: p.curToken}
	decl.Declaration = p.parseDeclaration()
	return decl
}

//
// This section is all about returning declarations and functions
//

func (p *Parser) parseDeclaration() ast.Decl {
	ident := &ast.Ident{
		Token: p.peekToken,
		Value: p.peekToken.Literal,
	}
	switch p.curToken.Type {
	case token.METHOD:
		meth := p.parseFunctionDeclaration()
		meth.Name = ident
		return meth
	case token.DESCRIBE:
	case token.OBJECT:
	}
	return nil
}

func (p *Parser) parseFunctionDeclaration() *ast.FuncDeclaration {
	fn := &ast.FuncDeclaration{Token: p.curToken}
	p.nextToken()
	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		fn.Parameters = p.parseFunctionParameters()
	} else {
		fn.Parameters = nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fn.Body = p.parseBlockStatement()
	return fn
}

func (p *Parser) parseFuncLiteral() ast.Expr {
	lit := &ast.FuncLiteral{Token: p.curToken}
	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		lit.Parameters = p.parseFunctionParameters()
	} else {
		lit.Parameters = nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Stmt{}
	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseFunctionParameters() []*ast.Ident {
	var idents []*ast.Ident
	if p.peekTokenIs(token.RBRACE) {
		return nil
	}
	p.nextToken()
	ident := &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}
	idents = append(idents, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}
		idents = append(idents, ident)
	}

	return idents
}

func (p *Parser) parseIntLiteral() ast.Expr {
	lit := &ast.IntLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("line %v, col %v: could not parse %q as integer",
			p.curToken.Line, p.curToken.Col, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expr {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expr {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACK)
	return array
}

func (p *Parser) parseHashMap() ast.Expr {
	hash := &ast.HashMap{Token: p.curToken}
	hash.Pairs = make(map[ast.Expr]ast.Expr)

hashLoop:
	for !p.peekTokenIs(token.RBRACE) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		for p.curTokenIs(token.NEWLINE) {
			if p.peekTokenIs(token.RBRACE) {
				break hashLoop
			}
			p.nextToken()
		}
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value
		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return hash
}

func (p *Parser) noPrefixParseFnError(t token.Token) {
	errToken := fmt.Sprintf("%s", t.Type)
	if t.Type == token.ILLEGAL {
		errToken = fmt.Sprintf("%s (%s)", t.Type, t.Literal)
	}
	msg := fmt.Sprintf("line %v, col %v: no prefix parse function for %s found", t.Line, t.Col, errToken)
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

	// Temporarily handle a scenario like: x 4 "hello"
	switch {
	case p.peekTokenIs(token.IDENT), p.peekTokenIs(token.INT), p.peekTokenIs(token.STRING):
		msg := fmt.Sprintf("line %v, col %v: no operator found between %q and %q",
			p.curToken.Line, p.curToken.Col, p.curToken.Literal, p.peekToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	for !p.peekTokenIs(token.SEMICOLON) &&
		!p.peekTokenIs(token.NEWLINE) &&
		!p.peekTokenIs(token.EOF) &&
		!p.peekTokenIs(token.PASSTHROUGH) &&
		prio < p.peekPriority() {
		//if _, ok := leftExp.(*ast.IfExpression); !ok {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
		//}
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

func (p *Parser) parsePostfixExpressionStatement(left ast.Stmt) ast.Stmt {
	postfix := p.postfixParseFns[p.curToken.Type]
	if postfix == nil {
		return left
	}
	stmt := postfix(left)
	return stmt
}

func (p *Parser) parseIndexExpression(left ast.Expr) ast.Expr {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACK) {
		return nil
	}
	return exp
}

func (p *Parser) parseCallExpression(function ast.Expr) ast.Expr {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Args = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expr {
	var args []ast.Expr
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) parseIdentity() ast.Expr {
	return &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expr {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expr {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expr {
	list := []ast.Expr{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}
	return list
}

func (p *Parser) parseIfExpression() ast.Expr {
	expression := &ast.IfExpression{Token: p.curToken}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()
	// TODO: Add Else if; change below to a for-loop
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		// TODO: !peekTokenIs(token.IF) {} else {
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}
	return expression
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
	msg := fmt.Sprintf("line %v, col %v: expected %s, got %s", t.Line, t.Col, expect, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
