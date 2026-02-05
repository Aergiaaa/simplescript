package evaluator

import (
	"github.com/Aergiaaa/idiotic_interpreter/ast"
	"github.com/Aergiaaa/idiotic_interpreter/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Bool{Value: true}
	FALSE = &object.Bool{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBoolObj(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpr(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpr(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatements(node)
	case *ast.IfExpression:
		return evalIfExpr(node)
	case *ast.ReturnStatement:
		return &object.ReturnValue{
			Value: Eval(node.ReturnValue),
		}
	}
	return nil
}

func evalIfExpr(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	}
	return NULL

}

func evalStatements(prog *ast.Program) object.Object {
	var res object.Object
	for _, statement := range prog.Statements {
		res = Eval(statement)

		if rv, ok := res.(*object.ReturnValue); ok {
			return rv.Value
		}
	}
	return res
}

func evalBlockStatements(block *ast.BlockStatement) (res object.Object) {
	for _, stmt := range block.Statements {
		res = Eval(stmt)

		if res != nil && res.Type() == object.RET_VAL_OBJ {
			return res
		}
	}
	return res
}

func evalInfixExpr(op string, left, right object.Object) object.Object {
	l, r := left.Type(), right.Type()
	switch {
	case isSameObjType(l, r, object.INTEGER_OBJ):
		return evalIntegInfixExpr(op, left, right)
	case isSameObjType(l, r, object.BOOL_OBJ):
		return evalBoolInfixExpr(op, left, right)
	default:
		return NULL
	}
}

func isSameObjType(left, right object.ObjectType, obj object.ObjectType) bool {
	return left == obj && right == obj
}

func evalPrefixExpr(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalNotOpExpr(right)
	case "-":
		return evalNegOpExpr(right)
	default:
		return NULL
	}
}

func evalIntegInfixExpr(op string, left, right object.Object) object.Object {
	lVal := left.(*object.Integer).Value
	rVal := right.(*object.Integer).Value

	switch op {
	case "+":
		return &object.Integer{Value: lVal + rVal}
	case "-":
		return &object.Integer{Value: lVal - rVal}
	case "*":
		return &object.Integer{Value: lVal * rVal}
	case "/":
		return &object.Integer{Value: lVal / rVal}
	case "==":
		return nativeBoolToBoolObj(lVal == rVal)
	case "!=":
		return nativeBoolToBoolObj(lVal != rVal)
	case ">=":
		return nativeBoolToBoolObj(lVal >= rVal)
	case "<=":
		return nativeBoolToBoolObj(lVal <= rVal)
	case ">":
		return nativeBoolToBoolObj(lVal > rVal)
	case "<":
		return nativeBoolToBoolObj(lVal < rVal)
	default:
		return NULL
	}
}

func evalBoolInfixExpr(op string, left, right object.Object) object.Object {
	lVal := left.(*object.Integer).Value
	rVal := right.(*object.Integer).Value

	switch op {
	case "==":
		return nativeBoolToBoolObj(lVal == rVal)
	case "!=":
		return nativeBoolToBoolObj(lVal != rVal)
	default:
		return NULL
	}

}

func evalNegOpExpr(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	val := right.(*object.Integer).Value
	return &object.Integer{
		Value: -val,
	}
}

func evalNotOpExpr(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL: // pretty sloppy but after this ill think this
		return TRUE
	default:
		return FALSE
	}
}

func nativeBoolToBoolObj(in bool) *object.Bool {
	if in {
		return TRUE
	}

	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default: // sloppy here too
		return true
	}
}
