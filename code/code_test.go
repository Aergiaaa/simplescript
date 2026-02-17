package code

import "testing"

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConst, []int{65535}, 2},
	}

	for _, tt := range tests {
		ins := Make(tt.op, tt.operands...)

		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		operandsRead, n := ReadOperands(def, ins[1:])
		if n != tt.bytesRead {
			t.Fatalf("n wrong. expected=%d, got=%d", tt.bytesRead, n)
		}

		for i, exp := range tt.operands {
			if operandsRead[i] != exp {
				t.Errorf("operand wrong. expected=%d, got=%d", exp, operandsRead[i])
			}
		}
	}
}

func TestInstructionsStr(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpConst, 2),
		Make(OpConst, 65535),
	}

	expected := `0000 OpAdd
0001 OpConst 2
0004 OpConst 65535
`

	var concatted Instructions
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted.\nwant=%q\ngot=%q",
			expected, concatted.String())
	}
}

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{
			OpConst, []int{65534}, []byte{byte(OpConst), 255, 254},
		},
		{
			OpAdd, []int{}, []byte{byte(OpAdd)},
		},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction wrong length, get=%d,expected=%d",
				len(instruction), len(tt.expected))
		}

		for i, b := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("wrong byte at position %d, get=%d, expected=%d", i, b, instruction[i])
			}
		}
	}
}
