package parser

import (
	"fmt"
	"github.com/alexjwhite-cb/jet/pkg/ast"
	"github.com/alexjwhite-cb/jet/pkg/lexer"
	"testing"
)

func TestValueStmts(t *testing.T) {
	tests := []struct {
		input       string
		expectIdent string
		expectValue interface{}
	}{
		{"x = 5", "x", 5},
		{"y = true", "y", true},
		{"foobar = y", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got %d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testValueStmt(t, stmt, tt.expectIdent) {
			return
		}

		val := stmt.(*ast.ValueStmt).Value
		if !testLiteralExpression(t, val, tt.expectValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input       string
		expectValue interface{}
	}{
		{"(5)->", 5},
		{"(true)->", true},
		{"(y)->", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Errorf("%+v", program.Statements[0])
			t.Fatalf("stmt not *ast.ReturnStatement, got %T", program.Statements[0])
		}

		if !testLiteralExpression(t, stmt.Value, tt.expectValue) {
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
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
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

func TestIfExpression(t *testing.T) {
	input := `if x < y { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got %d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStmt, got %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expressions is not ast.IfExpression, got %T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements, got %d\n", len(exp.Consequence.Statements))
		return
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStmt, got %T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements is not nil, got %+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if x < y { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got %d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStmt, got %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expressions is not ast.IfExpression, got %T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements, got %d\n", len(exp.Consequence.Statements))
		return
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStmt, got %T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("consequence is not 1 statements, got %d\n", len(exp.Alternative.Statements))
		return
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStmt, got %T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFuncDeclarationParsing(t *testing.T) {
	input := `meth add: x, y { x + y }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got %d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.DeclarationStmt)
	if !ok {
		t.Fatalf("Statements[0] is not ast.DeclarationStmt, got %T", program.Statements[0])
	}

	function, ok := stmt.Declaration.(*ast.FuncDeclaration)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FuncLiteral, got %T", stmt.Declaration)
	}

	if function.Name.Value != "add" {
		t.Fatalf("function.Value is not %q, got %q", "add", function.Name.Value)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function.Parameters does not contain %d parameters, got %d", 2, len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements does not contain %d statements, got %d", 1, len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not ast.ExpressionStmt, got %T", function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFuncLiteralParsing(t *testing.T) {
	input := `meth: x, y { x + y }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got %d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStmt, got %T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FuncLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FuncLiteral, got %T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function.Parameters does not contain %d parameters, got %d", 2, len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements does not contain %d statements, got %d", 1, len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not ast.ExpressionStmt, got %T", function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input        string
		expectParams []string
	}{
		{"meth {}", []string{}},
		{"meth: x {}", []string{"x"}},
		{"meth: x, y, z {}", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStmt)
		function := stmt.Expression.(*ast.FuncLiteral)

		if len(function.Parameters) != len(tt.expectParams) {
			t.Errorf("unexpected param length, want %d, got %d", len(function.Parameters), len(tt.expectParams))
		}

		for i, ident := range tt.expectParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got %d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStmt, got %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression, got %T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Args) != 3 {
		t.Fatalf("unexpected arg length, got %d", len(exp.Args))
	}

	testLiteralExpression(t, exp.Args[0], 1)
	testInfixExpression(t, exp.Args[1], 2, "*", 3)
	testInfixExpression(t, exp.Args[2], 4, "+", 5)
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

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStmt)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral, got %T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q, got %q", "hello world", literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, _ := program.Statements[0].(*ast.ExpressionStmt)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral, got %T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3, got %d", len(array.Elements))
	}

	testLiteralExpression(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1+1]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, _ := program.Statements[0].(*ast.ExpressionStmt)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not ast.IndexExpression, got %T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStmt)
	hash, ok := stmt.Expression.(*ast.HashMap)
	if !ok {
		t.Fatalf("exp is not ast.HashMap, got %T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length, got %d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for k, v := range hash.Pairs {
		literal, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral, got %T", k)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, v, expectedValue)
	}
}

func TestParsingHashLiteralsIntKeys(t *testing.T) {
	input := `{1: "one", 2: "two", 3: "three"}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStmt)
	hash, ok := stmt.Expression.(*ast.HashMap)
	if !ok {
		t.Fatalf("exp is not ast.HashMap, got %T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length, got %d", len(hash.Pairs))
	}

	expected := map[int64]string{
		1: "one",
		2: "two",
		3: "three",
	}

	for k, v := range hash.Pairs {
		literal, ok := k.(*ast.IntLiteral)
		if !ok {
			t.Errorf("key is not ast.IntLiteral, got %T", k)
		}
		if v.String() != expected[literal.Value] {
			t.Errorf("String at key does not match expectation, got %s", expected[literal.Value])
		}
	}
}

func TestParsingEmptyHashLiterals(t *testing.T) {
	input := `{}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStmt)
	hash, ok := stmt.Expression.(*ast.HashMap)
	if !ok {
		t.Fatalf("exp is not ast.HashMap, got %T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length, got %d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0+1, "two": 10-8, "three": 15/5}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStmt)
	hash, ok := stmt.Expression.(*ast.HashMap)
	if !ok {
		t.Fatalf("exp is not ast.HashMap, got %T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length, got %d", len(hash.Pairs))
	}

	tests := map[string]func(expr ast.Expr){
		"one": func(e ast.Expr) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expr) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expr) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral, got %T", key)
			continue
		}
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q", literal.String())
			continue
		}
		testFunc(value)
	}
}
