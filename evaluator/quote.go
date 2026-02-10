package evaluator

import (
	"fmt"

	"github.com/Aergiaaa/simplescript/ast"
	"github.com/Aergiaaa/simplescript/object"
	"github.com/Aergiaaa/simplescript/token"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	var err error
	node, err = evalUnquote(node, env)
	if err != nil {
		msg := fmt.Sprintf("error evaluating quote: %v", err)

		return &object.Error{
			Message: msg,
		}
	}

	return &object.Quote{
		Node: node,
	}
}

func evalUnquote(quoted ast.Node, env *object.Environment) (ast.Node, error) {
	return ast.Modify(quoted, func(n ast.Node) ast.Node {
		if !isUnquote(n) {
			return n
		}

		call, ok := n.(*ast.CallExpression)
		if !ok {
			return n
		}

		if len(call.Args) != 1 {
			return n
		}

		unquoted := Eval(call.Args[0], env)
		return convertObjToAstNode(unquoted)
	})
}

func isUnquote(n ast.Node) bool {
	callExpr, ok := n.(*ast.CallExpression)
	if !ok {
		return false
	}
	return callExpr.Func.TokenLiteral() == "unquote"
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

func convertObjToAstNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}

	case *object.Bool:
		var t token.Token
		if obj.Value {
			t = token.Token{Type: token.TRUE, Literal: "true"}
		} else {
			t = token.Token{Type: token.FALSE, Literal: "false"}
		}
		return &ast.Boolean{Token: t, Value: obj.Value}

	case *object.Quote:
		return obj.Node

	default:
		return nil
	}
}
