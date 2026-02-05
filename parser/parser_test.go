package parser

import (
	"fmt"
	"testing"

	"github.com/Aergiaaa/idiotic_interpreter/ast"
	"github.com/Aergiaaa/idiotic_interpreter/lexer"
)

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	l := lexer.InitLexer(input)
	p := InitParser(l)
	program := p.Parse()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Func, "add") {
		return
	}
	if len(exp.Args) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Args))
	}
	testLiteralExpression(t, exp.Args[0], 1)
	testInfixExpression(t, exp.Args[1], "*", 2, 3)
	testInfixExpression(t, exp.Args[2], "+", 4, 5)
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "ft() {};", expectedParams: []string{}},
		{input: "ft(x) {};", expectedParams: []string{"x"}},
		{input: "ft(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.InitLexer(tt.input)
		p := InitParser(l)
		program := p.Parse()
		checkParserErrors(t, p)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `ft(x, y) { x + y; }`

	l := lexer.InitLexer(input)
	p := InitParser(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "+", "x", "y")
}

func TestIFElseExpression(t *testing.T) {
	inp := `if (x < y) { x } else { y	}`

	l := lexer.InitLexer(inp)
	p := InitParser(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program body does not contain %d statement, got %d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not Expression statement, got %T\n", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.IFExpression)
	if !ok {
		t.Fatalf("expression is not 'if' expression, got=%T\n", stmt.Expression)
	}

	if !testInfixExpression(t, expr.Condition, "<", "x", "y") {
		return
	}

	if len(expr.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 element, got=%d\n", len(expr.Consequence.Statements))
	}

	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence is not expression statement, got=%T\n", expr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := expr.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alternative is not expression statement, got=%T\n", expr.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestIFExpression(t *testing.T) {
	inp := `if (x < y) { x }`

	l := lexer.InitLexer(inp)
	p := InitParser(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program body does not contain %d statement, got %d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not Expression statement, got %T\n", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.IFExpression)
	if !ok {
		t.Fatalf("expression is not 'if' expression, got=%T\n", stmt.Expression)
	}

	if !testInfixExpression(t, expr.Condition, "<", "x", "y") {
		return
	}

	if len(expr.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 element, got=%d\n", len(expr.Consequence.Statements))
	}

	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not expression statement, got=%T\n", expr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expr.Alternative != nil {
		t.Errorf("expr alternative is not nil, got=%+v\n", expr.Alternative)
	}
}

func TestOpHierarchyParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"!(true == false)",
			"(!(true == false))",
		},
		{
			"true", "true",
		},
		{
			"false", "false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a * b / c == d - e",
			"(((a * b) / c) == (d - e))",
		},
	}

	for _, tt := range tests {
		l := lexer.InitLexer(tt.input)
		p := InitParser(l)
		program := p.Parse()
		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q,got=%q", tt.expected, actual)
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input    string
		leftVal  any
		op       string
		rightVal any
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"false == false", false, "==", false},
		{"false != true", false, "!=", true},
	}

	for _, tt := range infixTests {
		l := lexer.InitLexer(tt.input)
		p := InitParser(l)
		program := p.Parse()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program does not contain %d statement, got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("statement is not expr, got=%T", program.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expr is not infix, got=%T", stmt.Expression)
		}

		if !testLiteralExpression(t, expr.Left, tt.leftVal) {
			return
		}

		if expr.Operator != tt.op {
			t.Fatalf("expr op is not %s, got=%s", tt.op, expr.Operator)
		}

		if !testLiteralExpression(t, expr.Right, tt.rightVal) {
			return
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		intVal   any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.InitLexer(tt.input)
		p := InitParser(l)
		program := p.Parse()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("exprStmt is not Expression Statement, got=%T", program.Statements[0])
		}

		expr, ok := exprStmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expr is not ast.PrefixExpression, get=%T", exprStmt.Expression)
		}
		if expr.Operator != tt.operator {
			t.Fatalf("operator is not %s, got=%s", tt.operator, expr.Operator)
		}

		if !testLiteralExpression(t, expr.Right, tt.intVal) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, val int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expr is not *ast.IntegerLiteral, got=%T", il)
		return false
	}

	if integ.Value != val {
		t.Errorf("value is not %d, got=%d", val, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", val) {
		t.Errorf("integer token literal is not %d, got=%s", val, integ.TokenLiteral())
		return false
	}

	return true
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
		// {"let fool = true;", true},
		// {"let smart = false;", false},
	}

	for _, tt := range tests {
		l := lexer.InitLexer(tt.input)
		p := InitParser(l)
		program := p.Parse()
		checkParserErrors(t, p)

		if len(program.Statements) == 0 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		exprStmt, exprOk := program.Statements[0].(*ast.ExpressionStatement)
		var letExpr *ast.LetStatement
		if !exprOk {
			var letOk bool
			letExpr, letOk = program.Statements[0].(*ast.LetStatement)
			if !letOk {
				t.Fatalf("exprStmt is not Expression Statement, got=%T", program.Statements[0])
			}

			b, ok := letExpr.Value.(*ast.Boolean)
			if !ok {
				t.Fatalf("expr is not Boolean, got=%T", letExpr.Value)
			}
			if b.Value != tt.expected {
				t.Errorf("value is not '%t', got=%t", tt.expected, b.Value)
			}

			return
		}

		b, ok := exprStmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("expr is not Boolean, got=%T", exprStmt.Expression)
		}

		if b.Value != tt.expected {
			t.Errorf("value is not '%t', got=%t", tt.expected, b.Value)
		}
	}

}

func TestIntegralLiteralExpression(t *testing.T) {
	inp := `5;`

	l := lexer.InitLexer(inp)
	p := InitParser(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exprStmt is not Expression Statement, got=%T", program.Statements[0])
	}

	literal, ok := exprStmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expr is not IntegralLiteral, got=%T", exprStmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}

func TestIdentifierExpression(t *testing.T) {
	inp := `foobar;`

	l := lexer.InitLexer(inp)
	p := InitParser(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exprStmt is not Expression Statement, got=%T", program.Statements[0])
	}

	ident, ok := exprStmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expr is not Identifier, got=%T", exprStmt.Expression)
	}
	expectedValue := "foobar"
	if ident.Value != expectedValue {
		t.Errorf("value of ident is not '%s', got=%q", expectedValue, ident.Value)
	}
	if ident.TokenLiteral() != expectedValue {
		t.Errorf("token literal of ident is not '%s', got=%q", expectedValue, ident.TokenLiteral())
	}
}

func TestReturnStatement(t *testing.T) {
	inp := `
		return 5;
		return 6767;
		return 10000000;
	`

	l := lexer.InitLexer(inp)
	p := InitParser(l)

	program := p.Parse()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("program return nil on parse")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("statement doesnt contain 3 statement. got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt is not a return stmt, got=%T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("stmt token literal is not 'return', got=%s", returnStmt.TokenLiteral())
		}
	}
}

func TestLetStatement(t *testing.T) {
	inp := `
		let x = 5;
		let y = 123;
		let foo = 666777;
	`

	l := lexer.InitLexer(inp)
	p := InitParser(l)

	program := p.Parse()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("program return nil on parse")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("statement doesnt contain 3 statement. got %d", len(program.Statements))
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
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("token literal is not 'let', got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement is not LetStatement, got=%q", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("statement value is not '%s', got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("statement name is not '%s', got=%s", name, letStmt.Name)
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, expr ast.Expression, op string, left, right any) bool {
	opExpr, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expr is not infix, got=%T(%s)", expr, expr)
		return false
	}

	if !testLiteralExpression(t, opExpr.Left, left) {
		return false
	}

	if opExpr.Operator != op {
		t.Errorf("operator is not %s, got=%s", op, opExpr.Operator)
		return false
	}

	if !testLiteralExpression(t, opExpr.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	case string:
		return testIdentifier(t, expr, v)
	case bool:
		return testBooleanLiteral(t, expr, v)
	}

	t.Errorf("type of expr is no handled, got=%T", expr)
	return false
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, val bool) bool {
	b, ok := expr.(*ast.Boolean)
	if !ok {
		t.Errorf("expr is not boolean, get=%T", expr)
		return false
	}

	if b.Value != val {
		t.Errorf("value is not '%t', got=%t", val, b.Value)
		return false
	}

	if b.TokenLiteral() != fmt.Sprint(val) {
		t.Errorf("token literal is not '%t', got=%s", val, b.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, expr ast.Expression, val string) bool {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Errorf("expr not Identifier, got=%T", expr)
		return false
	}

	if ident.Value != val {
		t.Errorf("value is not %s, got=%s", val, ident.Value)
		return false
	}

	if ident.TokenLiteral() != val {
		t.Errorf("token literal is not %s, got=%s", val, ident.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d error", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
