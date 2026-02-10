package evaluator

import (
	"github.com/Aergiaaa/simplescript/ast"
	"github.com/Aergiaaa/simplescript/object"
)

func DefineMacros(p *ast.Program, env *object.Environment) {
	var defs []int

	for i, stmt := range p.Statements {
		if isMacroDefinition(stmt) {
			addMacro(stmt, env)
			defs = append(defs, i)
		}
	}

	for i := len(defs) - 1; i >= 0; i-- {
		defIdx := defs[i]
		p.Statements = append(
			p.Statements[:defIdx],
			p.Statements[defIdx+1:]...,
		)
	}
}

func addMacro(stmt ast.Statement, env *object.Environment) {
	letStmt, _ := stmt.(*ast.LetStatement)
	macroLit, _ := letStmt.Value.(*ast.MacroLiteral)

	macro := &object.Macro{
		Params: macroLit.Parameters,
		Body:   macroLit.Body,
		Env:    env,
	}

	env.Set(letStmt.Name.Value, macro)
}

func ExpandMacros(program ast.Node, env *object.Environment) (ast.Node, error) {
	return ast.Modify(program, func(n ast.Node) ast.Node {
		callExpr, ok := n.(*ast.CallExpression)
		if !ok {
			return n
		}

		macro, ok := isMacroCall(callExpr, env)
		if !ok {
			return n
		}

		args := quoteArgs(callExpr)
		evalEnv := extendMacroEnv(macro, args)

		evaled := Eval(macro.Body, evalEnv)

		quote, ok := evaled.(*object.Quote)
		if !ok {
			panic("only support returning AST-nodes from macros")
		}

		return quote.Node
	})
}

func extendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
	extended := object.InitEnclosedEnv(macro.Env)

	for i, p := range macro.Params {
		extended.Set(p.Value, args[i])
	}

	return extended
}

func quoteArgs(expr *ast.CallExpression) []*object.Quote {
	var args []*object.Quote

	for _, a := range expr.Args {
		args = append(args, &object.Quote{Node: a})
	}
	return args
}

func isMacroCall(expr *ast.CallExpression, env *object.Environment) (*object.Macro, bool) {
	ident, ok := expr.Func.(*ast.Identifier)
	if !ok {
		return nil, false
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		return nil, false
	}

	macro, ok := val.(*object.Macro)
	if !ok {
		return nil, false
	}

	return macro, true
}

func isMacroDefinition(stmt ast.Statement) bool {
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		return false
	}

	_, ok = letStmt.Value.(*ast.MacroLiteral)
	if !ok {
		return false
	}

	return true
}
