package vm

import (
	"fmt"

	"github.com/Aergiaaa/simplescript/code"
	"github.com/Aergiaaa/simplescript/compiler"
	"github.com/Aergiaaa/simplescript/object"
)

const STACK_SIZE = 2 * 1024

var (
	True  = &object.Bool{Value: true}
	False = &object.Bool{Value: false}
)

type VM struct {
	consts       []object.Object
	instructions code.Instructions

	// stack pointer
	sp    int
	stack []object.Object
}

func Init(bc *compiler.Bytecode) *VM {
	return &VM{
		instructions: bc.Instructions,
		consts:       bc.Constants,

		stack: make([]object.Object, STACK_SIZE),
		sp:    0,
	}
}

func (v *VM) Run() error {
	for ip := 0; ip < len(v.instructions); ip++ {
		op := code.Opcode(v.instructions[ip])
		switch op {
		case code.OpConst:
			constIdx := code.ReadUint16(v.instructions[ip+1:])
			ip += 2

			err := v.push(v.consts[constIdx])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := v.execBinaryOp(op)
			if err != nil {
				return err
			}
		case code.OpEqual, code.OpNEQ, code.OpGT, code.OpGTE:
			err := v.execComparison(op)
			if err != nil {
				return err
			}
		case code.OpMin:
			err := v.execMinOp()
			if err != nil {
				return err
			}
		case code.OpBang:
			err := v.execBangOp()
			if err != nil {
				return err
			}
		case code.OpTrue:
			err := v.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := v.push(False)
			if err != nil {
				return err
			}
		case code.OpPop:
			v.pop()
		}
	}

	return nil
}

func (v *VM) execMinOp() error {
	op := v.pop()

	if op.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", op.Type())
	}

	val := op.(*object.Integer).Value
	return v.push(&object.Integer{
		Value: -val,
	})
}

func (v *VM) execBangOp() error {
	op := v.pop()

	switch op {
	case True:
		return v.push(False)
	case False:
		return v.push(True)
	default:
		return v.push(False)
	}
}

func (v *VM) execBinaryOp(op code.Opcode) error {
	right := v.pop()
	left := v.pop()

	rightT := right.Type()
	leftT := left.Type()

	if rightT == object.INTEGER_OBJ && leftT == object.INTEGER_OBJ {
		return v.execBinaryIntOp(op, left, right)
	}

	return fmt.Errorf("unsupported type of binary op: %s %s",
		leftT, rightT)
}

func (v *VM) execBinaryIntOp(op code.Opcode, left, right object.Object) error {
	leftV := left.(*object.Integer).Value
	rightV := right.(*object.Integer).Value

	var res int64

	switch op {
	case code.OpAdd:
		res = leftV + rightV
	case code.OpSub:
		res = leftV - rightV
	case code.OpMul:
		res = leftV * rightV
	case code.OpDiv:
		res = leftV / rightV
	default:
		return fmt.Errorf("unknown integer operation: %d", op)
	}

	return v.push(&object.Integer{
		Value: res,
	})
}

func (v *VM) execComparison(op code.Opcode) error {
	right := v.pop()
	left := v.pop()

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return v.execIntComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return v.push(nativeBoolToBoolObj(right == left))
	case code.OpNEQ:
		return v.push(nativeBoolToBoolObj(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)",
			op, left.Type(), right.Type())
	}
}

func (v *VM) execIntComparison(op code.Opcode, left, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return v.push(nativeBoolToBoolObj(leftVal == rightVal))
	case code.OpNEQ:
		return v.push(nativeBoolToBoolObj(leftVal != rightVal))
	case code.OpGT:
		return v.push(nativeBoolToBoolObj(leftVal > rightVal))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func nativeBoolToBoolObj(input bool) *object.Bool {
	if input {
		return True
	}
	return False
}

func (v *VM) push(obj object.Object) error {
	if v.sp >= STACK_SIZE {
		return fmt.Errorf("Stack Overflow")
	}

	v.stack[v.sp] = obj
	v.sp++

	return nil
}

func (v *VM) pop() object.Object {
	o := v.stack[v.sp-1]
	v.sp--

	return o
}

func (v *VM) StackTop() object.Object {
	if v.sp == 0 {
		return nil
	}

	return v.stack[v.sp-1]
}

func (v *VM) LastPoppedStackElem() object.Object {
	return v.stack[v.sp]
}
