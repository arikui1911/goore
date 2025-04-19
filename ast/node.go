package ast

import (
	"fmt"
	"io"

	"github.com/arikui1911/goore/token"
)

type Node interface {
	Location() *token.Location
	dump(io.Writer, int)
}

func Dump(tree Node, dest io.Writer) {
	tree.dump(dest, 0)
}

type Statement interface {
	Node
	statement()
}

type Program struct {
	Loc        *token.Location
	FileName   string
	Statements []Statement
	Err        error
}

func (*Program) statement() {}

func (n *Program) Location() *token.Location {
	return n.Loc
}

func indent(w io.Writer, lv int) {
	for i := 0; i < lv; i++ {
		fmt.Fprint(w, "  ")
	}
}

func dumpHeader(n Node, w io.Writer, lv int) {
	indent(w, lv)
	fmt.Fprintf(w, "%s:%T", n.Location(), n)
}

func attrHeader(s string, w io.Writer, lv int) {
	indent(w, lv)
	fmt.Fprintf(w, "@%s:\n", s)
}

func (n *Program) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintf(w, ": %s\n", n.FileName)
	for _, s := range n.Statements {
		s.dump(w, lv+1)
	}
}

type InvalidStatement struct {
	Loc *token.Location
	Err error
}

func (*InvalidStatement) statement() {}

func (n *InvalidStatement) Location() *token.Location {
	return n.Loc
}

func (n *InvalidStatement) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintf(w, ": %v\n", n.Err)
}

type Def struct {
	Loc  *token.Location
	Name *Identifier
	Init Expression
}

func (*Def) statement() {}

func (n *Def) Location() *token.Location {
	return n.Loc
}

func (n *Def) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintln(w, ":")
	attrHeader("Name", w, lv+1)
	n.Name.dump(w, lv+1)
	if n.Init == nil {
		return
	}
	attrHeader("Init", w, lv+1)
	n.Init.dump(w, lv+1)
}

type While struct {
	Loc  *token.Location
	Cond Expression
	Body []Statement
}

func (*While) statement() {}

func (n *While) Location() *token.Location {
	return n.Loc
}

func (n *While) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintln(w, ":")
	attrHeader("Cond", w, lv+1)
	n.Cond.dump(w, lv+1)
	attrHeader("Body", w, lv+1)
	for _, s := range n.Body {
		s.dump(w, lv+1)
	}
}

type Break struct {
	Loc *token.Location
}

func (*Break) statement() {}

func (n *Break) Location() *token.Location {
	return n.Loc
}

func (n *Break) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintln(w, "")
}

type Continue struct {
	Loc *token.Location
}

func (*Continue) statement() {}

func (n *Continue) Location() *token.Location {
	return n.Loc
}

func (n *Continue) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintln(w, "")
}

type Return struct {
	Loc        *token.Location
	Expression Expression
}

func (*Return) statement() {}

func (n *Return) Location() *token.Location {
	return n.Loc
}

func (n *Return) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	if n.Expression == nil {
		fmt.Fprintln(w, "")
		return
	}
	fmt.Fprintln(w, ":")
	n.Expression.dump(w, lv+1)
}

type If struct {
	Loc  *token.Location
	Test Expression
	Body []Statement
	Alt  Expression
}

func (*If) expression() {}

func (n *If) Location() *token.Location {
	return n.Loc
}

func (n *If) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintln(w, ":")
	attrHeader("Test", w, lv+1)
	n.Test.dump(w, lv+1)
	attrHeader("Body", w, lv+1)
	for _, s := range n.Body {
		s.dump(w, lv+1)
	}
	attrHeader("Alt", w, lv+1)
	n.Alt.dump(w, lv+1)
}

type Else struct {
	Loc  *token.Location
	Body []Statement
}

func (*Else) expression() {}

func (n *Else) Location() *token.Location {
	return n.Loc
}

func (n *Else) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintln(w, ":")
	for _, s := range n.Body {
		s.dump(w, lv+1)
	}
}

type ExpressionStatement struct {
	Loc        *token.Location
	Expression Expression
}

func (*ExpressionStatement) statement() {}

func (n *ExpressionStatement) Location() *token.Location {
	return n.Loc
}

func (n *ExpressionStatement) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintln(w, ":")
	n.Expression.dump(w, lv+1)
}

type Expression interface {
	Node
	expression()
}

type Identifier struct {
	Loc  *token.Location
	Name string
}

func (*Identifier) expression() {}

func (n *Identifier) Location() *token.Location {
	return n.Loc
}

func (n *Identifier) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintf(w, ": %s\n", n.Name)
}

type NilLiteral struct {
	Loc *token.Location
}

func (*NilLiteral) expression() {}

func (n *NilLiteral) Location() *token.Location {
	return n.Loc
}

func (n *NilLiteral) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintf(w, "\n")
}

type BoolLiteral struct {
	Loc   *token.Location
	Value bool
}

func (*BoolLiteral) expression() {}

func (n *BoolLiteral) Location() *token.Location {
	return n.Loc
}

func (n *BoolLiteral) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintf(w, ": %v\n", n.Value)
}

type IntLiteral struct {
	Loc   *token.Location
	Value int
}

func (*IntLiteral) expression() {}

func (n *IntLiteral) Location() *token.Location {
	return n.Loc
}

func (n *IntLiteral) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintf(w, ": %d\n", n.Value)
}

type FloatLiteral struct {
	Loc   *token.Location
	Value float64
}

func (*FloatLiteral) expression() {}

func (n *FloatLiteral) Location() *token.Location {
	return n.Loc
}

func (n *FloatLiteral) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintf(w, ": %f\n", n.Value)
}

type StringLiteral struct {
	Loc   *token.Location
	Value string
}

func (*StringLiteral) expression() {}

func (n *StringLiteral) Location() *token.Location {
	return n.Loc
}

func (n *StringLiteral) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintf(w, ": %#v\n", n.Value)
}

//go:generate stringer -type=Operation node.go
type Operation int

const (
	Plus Operation = iota
	Minus
	Not

	Eq
	Ne
	Le
	Ge
	Lt
	Gt
	Add
	Sub
	Mul
	Div
	Mod
)

type PrefixExpression struct {
	Loc      *token.Location
	Operator Operation
	Right    Expression
}

func (*PrefixExpression) expression() {}

func (n *PrefixExpression) Location() *token.Location {
	return n.Loc
}

func (n *PrefixExpression) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintf(w, ": %v\n", n.Operator)
	n.Right.dump(w, lv+1)
}

type ArrayLiteral struct {
	Loc      *token.Location
	Elements []Expression
}

func (*ArrayLiteral) expression() {}

func (n *ArrayLiteral) Location() *token.Location {
	return n.Loc
}

func (n *ArrayLiteral) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintln(w, ":")
	for _, e := range n.Elements {
		e.dump(w, lv+1)
	}
}

type HashLiteral struct {
	Loc   *token.Location
	Pairs []*HashEntry
}

func (*HashLiteral) expression() {}

func (n *HashLiteral) Location() *token.Location {
	return n.Loc
}

func (n *HashLiteral) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintln(w, ":")
	for _, p := range n.Pairs {
		p.dump(w, lv+1)
	}
}

type HashEntry struct {
	Loc   *token.Location
	Key   Expression
	Value Expression
}

func (*HashEntry) expression() {}

func (n *HashEntry) Location() *token.Location {
	return n.Loc
}

func (n *HashEntry) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintln(w, ":")
	attrHeader("Key", w, lv+1)
	n.Key.dump(w, lv+1)
	attrHeader("Value", w, lv+1)
	n.Value.dump(w, lv+1)
}

type FunctionLiteral struct {
	Loc        *token.Location
	Parameters []*Identifier
	Statements []Statement
}

func (*FunctionLiteral) expression() {}

func (n *FunctionLiteral) Location() *token.Location {
	return n.Loc
}

func (n *FunctionLiteral) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintln(w, ":")
	attrHeader("Parameters", w, lv+1)
	for _, p := range n.Parameters {
		p.dump(w, lv+1)
	}
	attrHeader("Statements", w, lv+1)
	for _, s := range n.Statements {
		s.dump(w, lv+1)
	}
}

type InfixExpression struct {
	Loc      *token.Location
	Operator Operation
	Left     Expression
	Right    Expression
}

func (*InfixExpression) expression() {}

func (n *InfixExpression) Location() *token.Location {
	return n.Loc
}

func (n *InfixExpression) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintf(w, ": %v\n", n.Operator)
	attrHeader("Left", w, lv+1)
	n.Left.dump(w, lv+1)
	attrHeader("Right", w, lv+1)
	n.Right.dump(w, lv+1)
}

type Call struct {
	Loc       *token.Location
	Function  Expression
	Arguments []Expression
}

func (*Call) expression() {}

func (n *Call) Location() *token.Location {
	return n.Loc
}

func (n *Call) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintln(w, ":")
	attrHeader("Function", w, lv+1)
	n.Function.dump(w, lv+1)
	attrHeader("Arguments", w, lv+1)
	for _, a := range n.Arguments {
		a.dump(w, lv+1)
	}
}

type KeyAccess struct {
	Loc       *token.Location
	Container Expression
	Key       Expression
}

func (*KeyAccess) expression() {}

func (n *KeyAccess) Location() *token.Location {
	return n.Loc
}

func (n *KeyAccess) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintln(w, ":")
	attrHeader("Container", w, lv+1)
	n.Container.dump(w, lv+1)
	attrHeader("Key", w, lv+1)
	n.Key.dump(w, lv+1)
}

type Let struct {
	Loc   *token.Location
	Left  *Identifier
	Right Expression
}

func (*Let) expression() {}

func (n *Let) Location() *token.Location {
	return n.Loc
}

func (n *Let) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintln(w, ":")
	attrHeader("Left", w, lv+1)
	n.Left.dump(w, lv+1)
	attrHeader("Right", w, lv+1)
	n.Right.dump(w, lv+1)
}

type KeyAssign struct {
	Loc   *token.Location
	Left  *KeyAccess
	Right Expression
}

func (*KeyAssign) expression() {}

func (n *KeyAssign) Location() *token.Location {
	return n.Loc
}

func (n *KeyAssign) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintln(w, ":")
	attrHeader("Left", w, lv+1)
	n.Left.dump(w, lv+1)
	attrHeader("Right", w, lv+1)
	n.Right.dump(w, lv+1)
}
