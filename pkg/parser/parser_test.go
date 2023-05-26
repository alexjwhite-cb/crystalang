package parser

import (
	"fmt"
	"github.com/alexjwhite-cb/jet/pkg/ast"
	"github.com/alexjwhite-cb/jet/pkg/lexer"
	"testing"
)

func TestValueStmts(t *testing.T) {
	input := `
x = 5
y = 10
foo = 9845`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statement. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testValueStmt(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testValueStmt(t *testing.T, s ast.Stmt, name string) bool {
	valueStmt, ok := s.(*ast.ValueStmt)
	if !ok {
		t.Errorf("s not *ast.ValueStmt, got %T", s)
		return false
	}

	if valueStmt.Name.Value != name {
		t.Errorf("valueStmt.Name.Value not '%s', got %s", name, valueStmt.Name.Value)
		return false
	}

	if valueStmt.Name.TokenLiteral() != name {
		t.Errorf("valueStmt.Name.TokenLiteral() not %s, got %s", name, valueStmt.Name.TokenLiteral())
		return false
	}
	return true
}

func TestIdentifierExpr(t *testing.T) {
	input := "foobar"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStmt, got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Ident)
	if !ok {
		t.Fatalf("exp not *ast.Ident, got %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s, got %s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() not %s, got %s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpr(t *testing.T) {
	input := "5"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStmt, got %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral, got %T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d, got %d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("ident.TokenLiteral() not %s, got %s", "5", literal.TokenLiteral())
	}
}

func TestBooleanLiteralExpr(t *testing.T) {
	input := "true"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStmt, got %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral, got %T", stmt.Expression)
	}

	if literal.Value != true {
		t.Errorf("literal.Value not %v, got %v", true, literal.Value)
	}

	if literal.TokenLiteral() != "true" {
		t.Errorf("ident.TokenLiteral() not %s, got %s", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		intValue interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements, got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStmt, got %T", stmt.Expression)
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpressions, got %T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s', got %s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.intValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expr, value int64) bool {
	integ, ok := il.(*ast.IntLiteral)
	if !ok {
		t.Errorf("intlit not *ast.IntLiteral, got %T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Val not %d, got %d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d, got %s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input    string
		leftVal  interface{}
		operator string
		rightVal interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"5 <= 5", 5, "<=", 5},
		{"5 >= 5", 5, ">=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements, got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStmt, got %T", stmt.Expression)
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftVal, tt.operator, tt.rightVal) {
			return
		}
	}
}

func TestOperatorPriorityParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("Expected: %q, got %q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expr, value string) bool {
	ident, ok := exp.(*ast.Ident)
	if !ok {
		t.Errorf("exp not *ast.Identifier, got %T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s, got %s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s, got %s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expr, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean, got %T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t, got %t", value, bo.Value)
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %v, got %v", value, bo.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expr, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled, got %T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expr, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp not *ast.Identifier, got %T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator not %s, got %q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}
