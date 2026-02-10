package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/Aergiaaa/simplescript/ast"
)

type ObjectType string

const (
	NULL_OBJ    = "NULL"
	ERR_OBJ     = "ERROR"
	FUNC_OBJ    = "FUNCTION"
	RET_VAL_OBJ = "RETURN_VALUE"
	BREAK_OBJ   = "BREAK"
	INTEGER_OBJ = "INTEGER"
	BOOL_OBJ    = "BOOL"
	STRING_OBJ  = "STRING"
	ARR_OBJ     = "ARRAY"
	HASH_OBJ    = "HASH"
	QUOTE_OBJ   = "QUOTE"
	MACRO_OBJ   = "MACRO"
	BUILTIN_OBJ = "BUILTIN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type BuiltinFn func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFn
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Macro struct {
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	Env    *Environment
}

func (m *Macro) Type() ObjectType { return MACRO_OBJ }
func (m *Macro) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range m.Params {
		params = append(params, p.String())
	}

	out.WriteString("macro")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(m.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() ObjectType { return QUOTE_OBJ }
func (q *Quote) Inspect() string {

	return "QUOTE(" + q.Node.String() + ")"
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
	Params []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Inspect() string {
	var output bytes.Buffer

	var params []string
	for _, p := range f.Params {
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

type Break struct{}

func (b *Break) Inspect() string  { return "break" }
func (b *Break) Type() ObjectType { return BREAK_OBJ }

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

type Array struct {
	Elems []Object
}

func (a *Array) Type() ObjectType { return ARR_OBJ }
func (a *Array) Inspect() string {
	var output bytes.Buffer

	var elems []string
	for _, e := range a.Elems {
		elems = append(elems, e.Inspect())
	}

	output.WriteString("[")
	output.WriteString(strings.Join(elems, ", "))
	output.WriteString("]")

	return output.String()
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var output bytes.Buffer

	var pairs []string
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Val.Inspect()))
	}

	output.WriteString("{")
	output.WriteString(strings.Join(pairs, ", "))
	output.WriteString("}")

	return output.String()
}

type Hashable interface {
	HashKey() HashKey
}

type HashPair struct {
	Key Object
	Val Object
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (i *Integer) HashKey() HashKey {
	return HashKey{
		Type:  i.Type(),
		Value: uint64(i.Value),
	}
}

func (b *Bool) HashKey() HashKey {
	var val uint64

	if b.Value {
		val = 1
	} else {
		val = 0
	}

	return HashKey{
		Type:  b.Type(),
		Value: val,
	}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{
		Type:  s.Type(),
		Value: h.Sum64(),
	}
}
