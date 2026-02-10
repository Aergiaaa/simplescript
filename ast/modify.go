package ast

import "fmt"

type ModifierFunc func(Node) Node

func Modify(n Node, mf ModifierFunc) (Node, error) {
	var ok bool
	var err error
	var modified Node

	switch n := n.(type) {
	case *Program:
		for i, stmt := range n.Statements {
			modified, err = Modify(stmt, mf)
			if err != nil {
				return nil, err
			}
			if modified != nil {
				n.Statements[i], ok = modified.(Statement)
			}
		}
	case *ExpressionStatement:
		modified, err = Modify(n.Expression, mf)
		if err != nil {
			return nil, err
		}
		if modified != nil {
			n.Expression, ok = modified.(Expression)
		}
	case *InfixExpression:
		modified, err = Modify(n.Left, mf)
		if err != nil {
			return nil, err
		}
		if modified != nil {
			n.Left, ok = modified.(Expression)
		}
		modified, err = Modify(n.Right, mf)
		if err != nil {
			return nil, err
		}
		if modified != nil {
			n.Right, ok = modified.(Expression)
		}
	case *PrefixExpression:
		modified, err = Modify(n.Right, mf)
		if err != nil {
			return nil, err
		}
		if modified != nil {
			n.Right, ok = modified.(Expression)
		}
	case *IndexExpression:
		modified, err = Modify(n.Left, mf)
		if err != nil {
			return nil, err
		}
		if modified != nil {
			n.Left, ok = modified.(Expression)
		}
		modified, err = Modify(n.Index, mf)
		if err != nil {
			return nil, err
		}
		if modified != nil {
			n.Index, ok = modified.(Expression)
		}
	case *IfExpression:
		modified, err = Modify(n.Condition, mf)
		if err != nil {
			return nil, err
		}
		if modified != nil {
			n.Condition, ok = modified.(Expression)
		}

		modified, err = Modify(n.Consequence, mf)
		if err != nil {
			return nil, err
		}
		if modified != nil {
			n.Consequence, ok = modified.(*BlockStatement)
		}
		if n.Alternative != nil {
			modified, err = Modify(n.Alternative, mf)
			if err != nil {
				return nil, err
			}
			if modified != nil {
				n.Alternative, ok = modified.(*BlockStatement)
			}
		}
	case *BlockStatement:
		for i := range n.Statements {
			modified, err = Modify(n.Statements[i], mf)
			if err != nil {
				return nil, err
			}
			if modified != nil {
				n.Statements[i], ok = modified.(Statement)
			}
		}
	case *ReturnStatement:
		modified, err = Modify(n.ReturnValue, mf)
		if err != nil {
			return nil, err
		}
		if modified != nil {
			n.ReturnValue, ok = modified.(Expression)
		}
	case *LetStatement:
		modified, err = Modify(n.Value, mf)
		if err != nil {
			return nil, err
		}
		if modified != nil {
			n.Value, ok = modified.(Expression)
		}
	case *FunctionLiteral:
		for i := range n.Parameters {
			modified, err = Modify(n.Parameters[i], mf)
			if err != nil {
				return nil, err
			}
			if modified != nil {
				n.Parameters[i], ok = modified.(*Identifier)
			}
		}
		modified, err = Modify(n.Body, mf)
		if err != nil {
			return nil, err
		}
		if modified != nil {
			n.Body, ok = modified.(*BlockStatement)
		}
	case *ArrayLiteral:
		for i := range n.Elems {
			modified, err = Modify(n.Elems[i], mf)
			if err != nil {
				return nil, err
			}
			if modified != nil {
				n.Elems[i], ok = modified.(Expression)
			}
		}
	case *HashLiteral:
		newPairs := make(map[Expression]Expression)
		for key, val := range n.Pairs {
			var newKey, newVal Expression
			modified, err = Modify(key, mf)
			if err != nil {
				return nil, err
			}
			if modified != nil {
				newKey, _ = modified.(Expression)
			}
			modified, err = Modify(val, mf)
			if err != nil {
				return nil, err
			}
			if modified != nil {
				newVal, _ = modified.(Expression)
			}

			newPairs[newKey] = newVal
		}
		n.Pairs = newPairs
	case *CallExpression:
		modified, err = Modify(n.Func, mf)
		if err != nil {
			return nil, err
		}

		if modified != nil {
			n.Func, ok = modified.(Expression)
		}

		for i := range n.Args {
			modified, err = Modify(n.Args[i], mf)
			if err != nil {
				return nil, err
			}

			if modified != nil {
				n.Args[i], ok = modified.(Expression)
			}
		}
	}

	if !ok && modified != nil {
		return nil, fmt.Errorf("expected Statement, got %T", modified)
	}

	return mf(n), nil
}
