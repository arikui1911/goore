package token

import "fmt"

type Location struct {
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

func (l Location) String() string {
	return fmt.Sprintf("(%d:%d):(%d:%d)", l.StartLine, l.StartColumn, l.EndLine, l.EndColumn)
}

type TokenTag int

//go:generate stringer -Type=TokenTag token.go
const (
	Invalid TokenTag = iota
	EOF

	IntLiteral
	FloatLiteral
	StringLiteral
	Identifier

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
	Let
	LetAdd
	LetSub
	LetMul
	LetDiv
	LetMod
	Arrow
	Comma
	Colon
	Semicolon
	Newline
	LeftParen
	RightParen
	LeftBrace
	RightBrace
	LeftBracket
	RightBracket

	True
	False
	Nil
	Def
	If
	Elsif
	Else
	While
	Break
	Continue
	Return
)

type Token struct {
	Tag      TokenTag
	Value    string
	Location Location
}

func (t Token) String() string {
	return fmt.Sprintf("#<%T:%s>:%s:%#v", t, t.Location, t.Tag, t.Value)
}
