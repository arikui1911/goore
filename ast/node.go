package ast

import "io"

type Node interface {
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
	FileName   string
	Statements []Statement
}

func (*Program) statement() {}

func (n *Program) dump(w io.Writer, lv int) {}

type If struct {
	Test Expression
	Body []Statement
	Alt  []Statement
}

func (*If) statement() {}

func (n *If) dump(w io.Writer, lv int) {}

type ExpressionStatement struct {
	Expression Expression
}

func (*ExpressionStatement) statement() {}

func (n *ExpressionStatement) dump(w io.Writer, lv int) {}

type Expression interface {
	Node
	expression()
}
