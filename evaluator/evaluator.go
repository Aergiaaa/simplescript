package evaluator

import (
	"fmt"

	"github.com/Aergiaaa/simplescript/ast"
	"github.com/Aergiaaa/simplescript/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Bool{Value: true}
	FALSE = &object.Bool{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
	case *ast.ArrayLiteral:
		elems := evalExprs(node.Elems, env)

		if len(elems) == 1 && isError(elems[0]) {
			return elems[0]
		}

		return &object.Array{
			Elems: elems,
		}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpr(left, index)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.CallExpression:
		f := Eval(node.Func, env)
		if isError(f) {
			return f
		}

		args := evalExprs(node.Args, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunc(f, args)
	case *ast.IfExpression:
		return evalIfExpr(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}

		return &object.ReturnValue{
			Value: val,
		}
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpr(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		if left == nil || right == nil {
			return NULL
		}

		return evalInfixExpr(node.Operator, left, right)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBoolObj(node.Value)
	}
	return nil
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERR_OBJ
	}
	return false
}

func applyFunc(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extEnv := extFuncEnv(fn, args)
		evaled := Eval(fn.Body, extEnv)

		return unwrapReturnVal(evaled)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extFuncEnv(f *object.Function, args []object.Object) *object.Environment {
	env := object.InitEnclosedEnv(f.Env)

	for i, param := range f.Parameters {
		env.Set(param.Value, args[i])
	}

	return env
}

func unwrapReturnVal(obj object.Object) object.Object {
	if retVal, ok := obj.(*object.ReturnValue); ok {
		return retVal.Value
	}

	return obj
}

func evalExprs(exprs []ast.Expression, env *object.Environment) []object.Object {
	var res []object.Object

	for _, e := range exprs {
		evaled := Eval(e, env)
		if isError(evaled) {
			return []object.Object{evaled}
		}

		res = append(res, evaled)
	}

	return res
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {

	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		val := Eval(valNode, env)
		if isError(val) {
			return val
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Val: val}
	}

	return &object.Hash{Pairs: pairs}
}

func evalIndexExpr(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARR_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpr(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpr(left, index)

	default:
		return newError("index operator is not supported: %s", left.Type())
	}
}

func evalHashIndexExpr(hash, index object.Object) object.Object {
	hashObj := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Val
}

func evalArrayIndexExpr(arr, index object.Object) object.Object {
	arrObj := arr.(*object.Array)
	i := index.(*object.Integer).Value

	maxIdx := int64(len(arrObj.Elems) - 1)

	if i < 0 {
		return newError("index out of bound: %d", i)
	}

	if i > maxIdx {
		return NULL
	}

	return arrObj.Elems[i]
}

func evalIfExpr(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}
	return NULL

}

func evalProgram(prog *ast.Program, env *object.Environment) object.Object {
	var res object.Object
	for _, stmt := range prog.Statements {
		res = Eval(stmt, env)

		switch res := res.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}
	}

	return res
}

func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) (res object.Object) {
	for _, stmt := range block.Statements {
		res = Eval(stmt, env)

		if res != nil {
			rt := res.Type()
			if rt == object.RET_VAL_OBJ || rt == object.ERR_OBJ {
				return res
			}
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
	case isSameObjType(l, r, object.STRING_OBJ):
		return evalStringInfixExpr(op, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
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
		return newError("unknown operator: %s%s", op, right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalStringInfixExpr(op string, left, right object.Object) object.Object {
	if op != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}

	lVal := left.(*object.String).Value
	rVal := right.(*object.String).Value

	return &object.String{
		Value: lVal + rVal,
	}
}

func evalBoolInfixExpr(op string, left, right object.Object) object.Object {
	lVal := left.(*object.Bool).Value
	rVal := right.(*object.Bool).Value

	switch op {
	case "==":
		return nativeBoolToBoolObj(lVal == rVal)
	case "!=":
		return nativeBoolToBoolObj(lVal != rVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}

}

func evalNegOpExpr(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
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
