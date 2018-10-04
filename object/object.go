package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
	"wordbuilder/ast"
)

type Type string

const (
	IntegerObj     = "INTEGER"
	BooleanObj     = "BOOLEAN"
	NullObj        = "NULL"
	ReturnValueObj = "RETURN_VALUE"
	ErrorObj       = "ERROR"
	FunctionObj    = "FUNCTION"
	StringObj      = "STRING"
	BuiltinObj     = "BUILTIN"
	ArrayObj       = "ARRAY"
	HashObj        = "HASH"
	WordObj        = "WORD"
	ReferenceObj   = "REF"
	ConceptObj     = "CPT"
	TranslationObj = "TR"
	MeThoughtObj   = "ME"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() Type {
	return IntegerObj
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() Type {
	return BooleanObj
}

type Null struct{}

func (n *Null) Type() Type {
	return NullObj
}

func (n *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() Type {
	return ReturnValueObj
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() Type {
	return ErrorObj
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() Type {
	return FunctionObj
}

func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() Type {
	return StringObj
}

func (s *String) Inspect() string {
	return s.Value
}

type BuiltinFunction func(env *Environment, args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() Type {
	return BuiltinObj
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}

type Array struct {
	Elements []Object
}

func (ao *Array) Type() Type { return ArrayObj }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashKey struct {
	Type  Type
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *HashPair) Type() Type {
	return HashObj
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type Hashable interface {
	HashKey() HashKey
}

func (h *Hash) Type() Type {
	return HashObj
}

type Word struct {
	Word       string
	Definition string
}

func (w *Word) Type() Type {
	return WordObj
}

func (w *Word) Inspect() string {
	return fmt.Sprintf("%s->{%s}", w.Word, w.Definition)
}

type Reference struct {
	Ref        string
	Definition string
}

func (ref *Reference) Type() Type {
	return ReferenceObj
}

func (ref *Reference) Inspect() string {
	return fmt.Sprintf("%s->{%s}", ref.Ref, ref.Definition)
}

type Concept struct {
	Concept    string
	Definition string
}

func (cpt *Concept) Type() Type {
	return ConceptObj
}

func (cpt *Concept) Inspect() string {
	return fmt.Sprintf("%s->{%s}", cpt.Concept, cpt.Definition)
}

type Translation struct {
	Translation string
	Definition  string
}

func (tr *Translation) Type() Type {
	return TranslationObj
}

func (tr *Translation) Inspect() string {
	return fmt.Sprintf("%s->{%s}", tr.Translation, tr.Definition)
}

type MeThought struct {
	Thought string
}

func (me *MeThought) Type() Type {
	return MeThoughtObj
}

func (me *MeThought) Inspect() string {
	return fmt.Sprintf("'%s'", me.Thought)
}
