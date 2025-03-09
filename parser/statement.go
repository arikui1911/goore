package parser

import (
	"github.com/arikui1911/goore/ast"
	"github.com/arikui1911/goore/token"
)

func parseProgram(p *Parser) (*ast.Program, error) {
	stmts, _, err := parseStatements(p, token.EOF)
	if err != nil {
		return nil, err
	}
	return &ast.Program{FileName: p.fileName, Statements: stmts}, nil
}

func parseStatements(p *Parser, term token.TokenTag) ([]ast.Statement, token.Token, error) {
	buf := []ast.Statement{}
	for {
		t, err := p.nextToken()
		if err != nil {
			return nil, token.Token{}, err
		}
		if t.Tag == term {
			return buf, t, nil
		}
		p.pushBack(t)
		s, err := parseStatement(p)
		if err != nil {
			return nil, token.Token{}, err
		}
		if s != nil {
			buf = append(buf, s)
		}
	}
}

type statementParser func(*Parser) (ast.Statement, error)

var statementParsers = map[token.TokenTag]statementParser{
	token.If: parseIf,
}

func parseStatement(p *Parser) (ast.Statement, error) {
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	p.pushBack(t)
	if fn, ok := statementParsers[t.Tag]; ok {
		return fn(p)
	}
	return parseExpressionStatement(p)
}

func parseExpressionStatement(p *Parser) (*ast.ExpressionStatement, error) {
	x, err := parseExpression(p)
	if err != nil {
		return nil, err
	}
	return &ast.ExpressionStatement{Expression: x}, nil
}

func parseIf(p *Parser) (ast.Statement, error) {
	return nil, nil
}
