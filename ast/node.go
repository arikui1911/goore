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
}

func (*Program) statement() {}

func (n *Program) Location() *token.Location {
	return n.Loc
}

func dumpHeader(n Node, w io.Writer, lv int) {
	for i := 0; i < lv; i++ {
		fmt.Fprint(w, "  ")
	}
	fmt.Fprintf(w, "%s:%T", n.Location(), n)
}

func (n *Program) dump(w io.Writer, lv int) {
	dumpHeader(n, w, lv)
	fmt.Fprintf(w, ": %s\n", n.FileName)
	for _, s := range n.Statements {
		s.dump(w, lv+1)
	}
}

type If struct {
	Loc  *token.Location
	Test Expression
	Body []Statement
	Alt  []Statement
}

func (*If) statement() {}

func (n *If) Location() *token.Location {
	return n.Loc
}

func (n *If) dump(w io.Writer, lv int) {}

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
