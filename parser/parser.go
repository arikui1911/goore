package parser

import (
	"fmt"
	"io"
	"strings"

	"github.com/arikui1911/goore/ast"
	"github.com/arikui1911/goore/lexer"
	"github.com/arikui1911/goore/token"
)

func ParseReader(src io.Reader, fileName string) (*ast.Program, error) {
	return New(lexer.New(src), fileName).Parse()
}

func ParseString(src string, fileName string) (*ast.Program, error) {
	return ParseReader(strings.NewReader(src), fileName)
}

type Parser struct {
	lexer         *lexer.Lexer
	fileName      string
	savedToken    token.Token
	hasSavedToken bool
	errs          []error
}

func New(l *lexer.Lexer, fileName string) *Parser {
	return &Parser{
		lexer:    l,
		fileName: fileName,
		errs:     []error{},
	}
}

func (p *Parser) Parse() (*ast.Program, error) {
	return parseProgram(p)
}

func (p *Parser) nextToken() (token.Token, error) {
	if p.hasSavedToken {
		p.hasSavedToken = false
		return p.savedToken, nil
	}
	t, err := p.lexer.NextToken()
	if err != nil {
		return token.Token{}, fmt.Errorf("%s:%w", p.fileName, err)
	}
	return t, nil
}

func (p *Parser) pushBack(t token.Token) {
	p.hasSavedToken = true
	p.savedToken = t
}

func (p *Parser) peekToken() (token.Token, error) {
	t, err := p.nextToken()
	if err != nil {
		return token.Token{}, err
	}
	p.pushBack(t)
	return t, nil
}

func (p *Parser) unexpected(t token.Token, ext string) error {
	return fmt.Errorf("%s:%s: unexpected token - %#v(%s) %s", p.fileName, t.Location, t.Value, t.Tag, ext)
}

func (p *Parser) addError(err error) {
	p.errs = append(p.errs, err)
}

func setLocation(loc *token.Location, beg *token.Location, end *token.Location) *token.Location {
	if loc == nil {
		loc = &token.Location{}
	}
	if beg != nil {
		loc.StartLine = beg.StartLine
		loc.StartColumn = beg.StartColumn
	}
	if end != nil {
		loc.EndLine = end.EndLine
		loc.EndColumn = end.EndColumn
	}
	return loc
}
