package compiler

import (
	"fmt"
	"testing"

	"github.com/Aergiaaa/simplescript/ast"
	"github.com/Aergiaaa/simplescript/code"
	"github.com/Aergiaaa/simplescript/lexer"
	"github.com/Aergiaaa/simplescript/object"
	"github.com/Aergiaaa/simplescript/parser"
)

type compilerTestCase struct {
	input                string
	expectedConsts       []any
	expectedInstructions []code.Instructions
}

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:          "!true",
			expectedConsts: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
		{
			input:          "1 >= 2",
			expectedConsts: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpGTE),
				code.Make(code.OpPop),
			},
		},
		{
			input:          "1 > 2",
			expectedConsts: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpGT),
				code.Make(code.OpPop),
			},
		},
		{
			input:          "1 < 2",
			expectedConsts: []any{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpGT),
				code.Make(code.OpPop),
			},
		},
		{
			input:          "1 == 2",
			expectedConsts: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:          "1 != 2",
			expectedConsts: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpNEQ),
				code.Make(code.OpPop),
			},
		},
		{
			input:          "true == false",
			expectedConsts: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:          "true != false",
			expectedConsts: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpNEQ),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:          "-1",
			expectedConsts: []any{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpMin),
				code.Make(code.OpPop),
			},
		},
		{
			input:          "true",
			expectedConsts: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{
			input:          "false",
			expectedConsts: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			input:          "1 - 2",
			expectedConsts: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			input:          "1 * 2",
			expectedConsts: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			input:          "2 / 1",
			expectedConsts: []any{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
		{
			input:          "1 + 2",
			expectedConsts: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:          "1; 2",
			expectedConsts: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConst, 1),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := Init()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("instruction comparation failed: %s", err)
		}

		err = testConsts(tt.expectedConsts, bytecode.Constants)
		if err != nil {
			t.Fatalf("constants comparation failed: %s", err)
		}
	}
}

func testInstructions(expected []code.Instructions, instructions code.Instructions) error {
	concatted := concatInstructions(expected)

	if len(instructions) != len(concatted) {
		return fmt.Errorf("wrong instructions length,\nwant=%q\nget=%q",
			len(concatted), len(instructions))
	}

	for i, ins := range concatted {
		if instructions[i] != ins {
			return fmt.Errorf("wrong instruction at %d,\nwant=%q\nget=%q",
				i, ins, instructions[i])
		}
	}

	return nil
}

func testConsts(expected []any, object []object.Object) error {
	if len(expected) != len(object) {
		return fmt.Errorf("wrong number of constants, want=%d,got=%q",
			len(expected), len(object))
	}

	for i, consts := range expected {
		switch c := consts.(type) {
		case int:
			if err := testIntegerObject(int64(c), object[i]); err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s",
					i, err)
			}
		}
	}

	return nil
}

func testIntegerObject(expected int64, obj object.Object) any {
	res, ok := obj.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not integer, get=%T (%+v)",
			obj, obj)
	}

	if res.Value != expected {
		return fmt.Errorf("object wrong value, expected=%d, get=%d", expected, res.Value)
	}

	return nil
}

func concatInstructions(ins []code.Instructions) code.Instructions {
	var res code.Instructions

	for _, in := range ins {
		res = append(res, in...)
	}
	return res
}

func parse(input string) *ast.Program {
	l := lexer.Init(input)
	p := parser.Init(l)

	return p.Parse()
}
