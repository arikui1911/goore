package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/arikui1911/goore/token"
)

type Lexer struct {
	src          *bufio.Reader
	line         int
	col          int
	lastNlCol    int
	savedRune    rune
	hasSavedRune bool
	lastTag      token.TokenTag
}

func New(src io.Reader) *Lexer {
	return &Lexer{
		src: bufio.NewReader(src),
	}
}

var keywords = map[string]token.TokenTag{
	"true":     token.True,
	"false":    token.False,
	"nil":      token.Nil,
	"def":      token.Def,
	"if":       token.If,
	"elsif":    token.Elsif,
	"else":     token.Else,
	"while":    token.While,
	"break":    token.Break,
	"continue": token.Continue,
	"return":   token.Return,
}

var operators = map[string]token.TokenTag{
	"==": token.Eq,
	"!=": token.Ne,
	"<=": token.Le,
	">=": token.Ge,
	"<":  token.Lt,
	">":  token.Gt,
	"+":  token.Add,
	"-":  token.Sub,
	"*":  token.Mul,
	"/":  token.Div,
	"%":  token.Mod,
	"=":  token.Let,
	"+=": token.LetAdd,
	"-=": token.LetSub,
	"*=": token.LetMul,
	"/=": token.LetDiv,
	"%=": token.LetMod,
	"->": token.Arrow,
	",":  token.Comma,
	":":  token.Colon,
	";":  token.Semicolon,
	"(":  token.LeftParen,
	")":  token.RightParen,
	"{":  token.LeftBrace,
	"[":  token.LeftBracket,
	"]":  token.RightBracket,
}

func isOperatorCandidate(s string) bool {
	for k := range operators {
		if strings.HasPrefix(k, s) {
			return true
		}
	}
	return false
}

type lexState int

const (
	initialState lexState = iota
	commentState
	zeroState
	intState
	floatState
	stringState
	stringEscState
	identState
	operatorState
)

func (l *Lexer) NextToken() (token.Token, error) {
	var buf []rune
	t := token.Token{Tag: token.EOF}
	state := initialState
	err := func() error {
		for {
			line := l.line
			col := l.col
			c, err := l.getc()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			switch state {
			case initialState:
				if c != '\n' && unicode.IsSpace(c) {
					continue
				}
				t.Location.StartLine = line
				t.Location.StartColumn = col
				t.Location.EndLine = line
				t.Location.EndColumn = col
				switch c {
				case '\n':
					if l.newlineRequired() {
						t.Tag = token.Newline
						return nil
					}
				case '}':
					t.Tag = token.RightBrace
					if l.newlineRequired() {
						l.ungetc(c)
						t.Tag = token.Newline
					}
					return nil
				case '#':
					state = commentState
				case '0':
					t.Tag = token.IntLiteral
					state = zeroState
				case '"':
					t.Tag = token.StringLiteral
					state = stringState
					buf = []rune{}
				default:
					buf = []rune{c}
					if unicode.IsDigit(c) {
						t.Tag = token.IntLiteral
						state = intState
					} else if c == '_' || unicode.IsLetter(c) {
						t.Tag = token.Identifier
						state = identState
					} else {
						state = operatorState
					}
				}
			case commentState:
				if c == '\n' {
					l.ungetc(c)
					state = initialState
				}
			case zeroState:
				if c != '.' {
					l.ungetc(c)
					t.Value = "0"
					return nil
				}
				buf = append(buf, '0')
				buf = append(buf, '.')
				t.Tag = token.FloatLiteral
				state = floatState
			case intState:
			case floatState:
			case stringState:
				switch c {
				case '\\':
					state = stringEscState
				case '"':
					t.Location.EndLine = line
					t.Location.EndColumn = col
					state = initialState
					return nil
				default:
					t.Location.EndLine = line
					t.Location.EndColumn = col
					buf = append(buf, c)
				}
			case stringEscState:
				switch c {
				case 'n':
					buf = append(buf, '\n')
				default:
					buf = append(buf, c)
				}
				t.Location.EndLine = line
				t.Location.EndColumn = col
				state = stringState
			case identState:
				if c != '_' && !unicode.IsLetter(c) && !unicode.IsDigit(c) {
					l.ungetc(c)
					return nil
				}
				buf = append(buf, c)
				t.Location.EndLine = line
				t.Location.EndColumn = col
			case operatorState:
				buf = append(buf, c)
				if !isOperatorCandidate(string(buf)) {
					l.ungetc(c)
					buf = buf[:len(buf)-1]
					return nil
				}
				t.Location.EndLine = line
				t.Location.EndColumn = col
			default:
				panic("must not happen")
			}
		}
	}()
	if err != nil {
		return token.Token{}, err
	}
	if buf != nil {
		t.Value = string(buf)
	}
	switch state {
	case stringState, stringEscState:
		return token.Token{}, fmt.Errorf("%d:%d: unterminated string literal", t.Location.StartLine, t.Location.StartColumn)
	case identState:
		if v, ok := keywords[t.Value]; ok {
			t.Tag = v
		}
	case operatorState:
		v, ok := operators[t.Value]
		if !ok {
			return token.Token{}, fmt.Errorf("%d:%d: invalid character - '%c'", t.Location.StartLine, t.Location.StartColumn, buf[0])
		}
		t.Tag = v
	}
	l.lastTag = t.Tag
	return t, nil
}

func (l *Lexer) newlineRequired() bool {
	return false
}

func (l *Lexer) getc() (c rune, err error) {
	if l.hasSavedRune {
		l.hasSavedRune = false
		c = l.savedRune
	} else {
		c, _, err = l.src.ReadRune()
	}
	if err != nil {
		return
	}
	if c == '\n' {
		l.lastNlCol = l.col
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return
}

func (l *Lexer) ungetc(c rune) {
	l.hasSavedRune = true
	l.savedRune = c
	l.col--
	if c == '\n' {
		l.line--
		l.col = l.lastNlCol
	}
}
