package parser

import (
	"fmt"
	"strconv"

	"github.com/arikui1911/goore/ast"
	"github.com/arikui1911/goore/token"
)

type Precedence int

const (
	lowestPrecedence Precedence = iota
	callPrecedence
	highestPrecedence
)

type prefixedParser func(*Parser) (ast.Expression, error)

var prefixedParsers = map[token.TokenTag]prefixedParser{
	token.Identifier:   parseIdentifier,
	token.IntLiteral:   parseIntLiteral,
	token.FloatLiteral: parseFloatLiteral,
}

func parseExpression(p *Parser, prec Precedence) (ast.Expression, error) {
	t, err := p.peekToken()
	if err != nil {
		return nil, err
	}
	fn, ok := prefixedParsers[t.Tag]
	if !ok {
		return nil, p.unexpected(t, "for beginning of expression")
	}
	left, err := fn(p)
	if err != nil {
		return nil, err
	}
	return left, nil
}

func parseIdentifier(p *Parser) (ast.Expression, error) {
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	return &ast.Identifier{Loc: &t.Location, Name: t.Value}, nil
}

func parseIntLiteral(p *Parser) (ast.Expression, error) {
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	i64, err := strconv.ParseInt(t.Value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%s:%s: %w", p.fileName, t.Location, err)
	}
	return &ast.IntLiteral{Loc: &t.Location, Value: int(i64)}, nil
}

func parseFloatLiteral(p *Parser) (ast.Expression, error) {
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	f64, err := strconv.ParseFloat(t.Value, 64)
	if err != nil {
		return nil, fmt.Errorf("%s:%s: %w", p.fileName, t.Location, err)
	}
	return &ast.FloatLiteral{Loc: &t.Location, Value: f64}, nil
}
