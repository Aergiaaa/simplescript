package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Aergiaaa/idiotic_interpreter/ast"
)

type ObjectType string

const (
	NULL_OBJ    = "NULL"
	ERR_OBJ     = "ERROR"
	FUNC_OBJ    = "FUNCTION"
	RET_VAL_OBJ = "RETURN_VALUE"
	INTEGER_OBJ = "INTEGER"
	BOOL_OBJ    = "BOOL"
	STRING_OBJ  = "STRING"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

type Error struct {
	Message string
}

func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Type() ObjectType { return ERR_OBJ }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Inspect() string {
	var output bytes.Buffer

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	output.WriteString("ft")
	output.WriteString("(")
	output.WriteString(strings.Join(params, ", "))
	output.WriteString(") {\n")
	output.WriteString(f.Body.String())
	output.WriteString("\n}")

	return output.String()
}
func (f *Function) Type() ObjectType {
	return FUNC_OBJ
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RET_VAL_OBJ }

type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Bool struct {
	Value bool
}

func (b *Bool) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Bool) Type() ObjectType { return BOOL_OBJ }
