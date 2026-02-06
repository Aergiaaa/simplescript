package ast

import (
	"bytes"
	"strings"

	"github.com/Aergiaaa/idiotic_interpreter/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var output bytes.Buffer

	for _, s := range p.Statements {
		output.WriteString(s.String())
	}

	return output.String()
}

type ExpressionStatement struct {
	Token token.Token // just first token of expr
	Expression
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

type LetStatement struct {
	Token token.Token // should always be token.LET
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) String() string {
	var output bytes.Buffer

	output.WriteString(ls.TokenLiteral() + " ")
	output.WriteString(ls.Name.String())
	output.WriteString(" = ")

	if ls.Value != nil {
		output.WriteString(ls.Value.String())
	}

	output.WriteString(";")

	return output.String()

}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type ReturnStatement struct {
	Token       token.Token // should always be token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) String() string {
	var output bytes.Buffer

	output.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		output.WriteString(rs.ReturnValue.String())
	}

	output.WriteString(";")

	return output.String()
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var output bytes.Buffer

	output.WriteString("if")
	output.WriteString(ie.Condition.String())
	output.WriteString(" ")
	output.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		output.WriteString("else ")
		output.WriteString(ie.Alternative.String())
	}

	return output.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) expressionNode()      {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var output bytes.Buffer

	for _, stmt := range bs.Statements {
		output.WriteString(stmt.String())
	}

	return output.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type CallExpression struct {
	Token token.Token // should be '(' because at i.e add(1,2), the 'add' identifier already consumed
	Func  Expression  // Identifier or Func Literal
	Args  []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var output bytes.Buffer

	var args []string
	for _, a := range ce.Args {
		args = append(args, a.String())
	}

	output.WriteString(ce.Func.String())
	output.WriteString("(")
	output.WriteString(strings.Join(args, ", "))
	output.WriteString(")")

	return output.String()
}

type FunctionLiteral struct {
	Token      token.Token // should be 'ft'
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var output bytes.Buffer

	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	output.WriteString(fl.TokenLiteral())
	output.WriteString("(")
	output.WriteString(strings.Join(params, ", "))
	output.WriteString(") ")
	output.WriteString(fl.Body.String())

	return output.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var output bytes.Buffer

	output.WriteString("(")
	output.WriteString(pe.Operator)
	output.WriteString(pe.Right.String())
	output.WriteString(")")

	return output.String()
}

type InfixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
	Left     Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var output bytes.Buffer

	output.WriteString("(")
	output.WriteString(ie.Left.String())
	output.WriteString(" " + ie.Operator + " ")
	output.WriteString(ie.Right.String())
	output.WriteString(")")

	return output.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) String() string       { return i.Value }
func (i *Identifier) statementNode()       {}
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
