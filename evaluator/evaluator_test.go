package evaluator

import (
	"testing"

	"github.com/Aergiaaa/simplescript/lexer"
	"github.com/Aergiaaa/simplescript/object"
	"github.com/Aergiaaa/simplescript/parser"
)

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
			"one": 10 - 9,
			two: 1 + 1,
			"thr" + "ee": 6 / 2,
			4: 4,
			true: 5,
			false: 6
	}`

	evaled := testEval(input)
	res, ok := evaled.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaled, evaled)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(res.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(res.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := res.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Val, expectedValue)
	}
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			"index out of bound: -1",
		},
	}

	for _, tt := range tests {
		evaled := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaled, int64(integer))
			continue
		}

		msg, ok := tt.expected.(string)
		if ok {
			testErrorObject(t, evaled, msg)
			continue
		}

		testNullObject(t, evaled)
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)

	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elems) != 3 {
		t.Fatalf("array has wrong num of elems. got=%d",
			len(result.Elems))
	}

	testIntegerObject(t, result.Elems[0], 1)
	testIntegerObject(t, result.Elems[1], 4)
	testIntegerObject(t, result.Elems[2], 6)
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments, got=2, want=1"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}

func TestStringConcatenationError(t *testing.T) {
	input := `"Hello" + " " + "World""!"`

	evaluated := testEval(input)

	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("object is not Error. got=%T (%+v)", evaluated, evaluated)
	}

	if errObj.Message == "" {
		t.Errorf("Expected error message for malformed string syntax")
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestIfInsideFunc(t *testing.T) {
	input := `
	let sub = ft(x,y) {
		if (y == 0) {
			return x + x;
		} else {
			return x - y;
		}
	};

	sub(1,0)
	`

	testIntegerObject(t, testEval(input), 2)
}

func TestClosures(t *testing.T) {
	input := `
		let newAdder = ft(x) {
			ft(y) { x + y };
		};
		let addTwo = newAdder(2);
		addTwo(2);`

	testIntegerObject(t, testEval(input), 4)
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = ft(x) { x; }; identity(5);", 5},
		{"let identity = ft(x) { return x; }; identity(5);", 5},
		{"let double = ft(x) { x * 2; }; double(5);", 10},
		{"let add = ft(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = ft(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"ft(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "ft(x) { x + 2; };"

	evaluated := testEval(input)
	ft, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(ft.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			ft.Parameters)
	}

	if ft.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", ft.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if ft.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, ft.Body.String())
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			`{"name": "Monkey"}[ft(x) { x }];`,
			"unusable as hash key: FUNCTION",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOL",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOL",
		},
		{
			"-true",
			"unknown operator: -BOOL",
		},
		{
			"true + false;",
			"unknown operator: BOOL + BOOL",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOL + BOOL",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOL + BOOL",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}
			`,
			"unknown operator: BOOL + BOOL",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}

				return 1;
			}
			`, 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoolObject(t, evaluated, tt.expected)
	}
}

func TestEvalBoolExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 <= 2", true},
		{"1 > 2", false},
		{"1 >= 2", false},
		{"1 < 1", false},
		{"1 <= 1", true},
		{"1 > 1", false},
		{"1 >= 1", true},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoolObject(t, evaluated, tt.expected)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.InitLexer(input)
	p := parser.InitParser(l)
	program := p.Parse()

	if len(p.Errors()) > 0 {
		return &object.Error{
			Message: p.Errors()[0],
		}
	}

	env := object.InitEnv()

	return Eval(program, env)
}

func testBoolObject(t *testing.T, obj object.Object, expected bool) bool {
	res, ok := obj.(*object.Bool)
	if !ok {
		t.Errorf("object is not bool, got=%T (%+v)", obj, obj)
		return false
	}

	if res.Value != expected {
		t.Errorf("got wrong value, got=%t, expected=%t", expected, res.Value)
		return false
	}

	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func testErrorObject(t *testing.T, obj object.Object, msg string) bool {
	err, ok := obj.(*object.Error)
	if !ok {
		t.Errorf("object is not ERROR, got=%T (%+v)", obj, obj)
		return false
	}

	if err.Message != msg {
		t.Errorf("wrong error message, got=%q, expected=%q", err.Message, msg)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL, got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
