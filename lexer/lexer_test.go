package lexer_test

import (
	"strings"
	"testing"

	"github.com/arikui1911/goore/lexer"
	"github.com/arikui1911/goore/token"
)

func TestLexFirstToken(t *testing.T) {
	table := []struct {
		name string
		src  string
		tag  token.TokenTag
		val  string
	}{
		{"identifier", `hoge`, token.Identifier, "hoge"},
		{"zero", `0`, token.IntLiteral, "0"},
		{"int literal", `123`, token.IntLiteral, "123"},
		{"float literal", `4.56`, token.FloatLiteral, "4.56"},
		{"string literal", `"Hello."`, token.StringLiteral, "Hello."},
		{"true", `true`, token.True, "true"},
		{"false", `false`, token.False, "false"},
		{"nil", `nil`, token.Nil, "nil"},
		{"def", `def`, token.Def, "def"},
		{"if", `if`, token.If, "if"},
		{"elsif", `elsif`, token.Elsif, "elsif"},
		{"else", `else`, token.Else, "else"},
		{"while", `while`, token.While, "while"},
		{"break", `break`, token.Break, "break"},
		{"continue", `continue`, token.Continue, "continue"},
		{"return", `return`, token.Return, "return"},
		{"eq", `==`, token.Eq, "=="},
		{"ne", `!=`, token.Ne, "!="},
		{"le", `<=`, token.Le, "<="},
		{"ge", `>=`, token.Ge, ">="},
		{"lt", `<`, token.Lt, "<"},
		{"gt", `>`, token.Gt, ">"},
		{"add", `+`, token.Add, "+"},
		{"sub", `-`, token.Sub, "-"},
		{"mul", `*`, token.Mul, "*"},
		{"div", `/`, token.Div, "/"},
		{"mod", `%`, token.Mod, "%"},
		{"let", `=`, token.Let, "="},
		{"let add", `+=`, token.LetAdd, "+="},
		{"let sub", `-=`, token.LetSub, "-="},
		{"let mul", `*=`, token.LetMul, "*="},
		{"let div", `/=`, token.LetDiv, "/="},
		{"let mod", `%=`, token.LetMod, "%="},
		{"arrow", `->`, token.Arrow, "->"},
		{"comma", `,`, token.Comma, ","},
		{"colon", `:`, token.Colon, ":"},
		{"semi-colon", `;`, token.Semicolon, ";"},
		{"left paren", `(`, token.LeftParen, "("},
		{"right paren", `)`, token.RightParen, ")"},
		{"left brace", `{`, token.LeftBrace, "{"},
		{"right brace", `}`, token.RightBrace, "}"},
		{"left bracket", `[`, token.LeftBracket, "["},
		{"right  bracket", `]`, token.RightBracket, "]"},
	}

	for _, e := range table {
		t.Run(e.name, func(t *testing.T) {
			l := lexer.New(strings.NewReader(e.src))
			testNextTokenTagAndValue(t, l, e.tag, e.val)
		})
	}
}

func TestLexIdentifiers(t *testing.T) {
	src := `
		x
		foo
		foo123
		_foo
		foo_bar
		foo_bar_123
		foo_123_456_bar
	`
	seq := []struct {
		name string
		val  string
	}{
		{"one letter", "x"},
		{"several letters", "foo"},
		{"letters and number", "foo123"},
		{"underline started", "_foo"},
		{"underline joined", "foo_bar"},
		{"underline joined words and number", "foo_bar_123"},
		{"joined words and numbers in middle", "foo_123_456_bar"},
	}

	l := lexer.New(strings.NewReader(src))
	for _, e := range seq {
		t.Run(e.name, func(t *testing.T) {
			testNextTokenTagAndValue(t, l, token.Identifier, e.val)
			testNextTokenTagAndValue(t, l, token.Newline, "\n")
		})
	}
	t.Run("EOF", func(t *testing.T) {
		testNextTokenTagAndValue(t, l, token.EOF, "")
	})
}

func testNextTokenTagAndValue(t *testing.T, l *lexer.Lexer, tag token.TokenTag, val string) {
	r, err := l.NextToken()
	if err != nil {
		t.Error(err)
		return
	}
	if r.Tag != tag {
		t.Errorf("%s: want <%s> got <%s>", r, tag, r.Tag)
		return
	}
	if r.Value != val {
		t.Errorf("%s: want <%#v> got <%#v>", r, val, r.Value)
		return
	}
}
