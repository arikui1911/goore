package parser

import (
	"fmt"
	"io"
	"strings"

	"github.com/arikui1911/goore/ast"
	"github.com/arikui1911/goore/lexer"
	"github.com/arikui1911/goore/token"
)

func ParseReader(src io.Reader, fileName string) (ast.Node, error) {
	return New(lexer.New(src), fileName).Parse()
}

func ParseString(src string, fileName string) (ast.Node, error) {
	return ParseReader(strings.NewReader(src), fileName)
}

type Parser struct {
	lexer         *lexer.Lexer
	fileName      string
	savedToken    token.Token
	hasSavedToken bool
}

func New(l *lexer.Lexer, fileName string) *Parser {
	return &Parser{
		lexer:    l,
		fileName: fileName,
	}
}

func (p *Parser) Parse() (ast.Node, error) {
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
