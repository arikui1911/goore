package token

type Location struct {
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

type TokenTag int

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
