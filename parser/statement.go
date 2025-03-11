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
	loc := &token.Location{StartLine: 1, StartColumn: 1, EndLine: 1, EndColumn: 1}
	if len(stmts) > 0 {
		setLocation(loc, stmts[0].Location(), stmts[len(stmts)-1].Location())
	}
	return &ast.Program{Loc: loc, FileName: p.fileName, Statements: stmts}, nil
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
	t, err := p.peekToken()
	if err != nil {
		return nil, err
	}
	if fn, ok := statementParsers[t.Tag]; ok {
		return fn(p)
	}
	return parseExpressionStatement(p)
}

func parseExpressionStatement(p *Parser) (*ast.ExpressionStatement, error) {
	x, err := parseExpression(p, lowestPrecedence)
	if err != nil {
		return nil, err
	}
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	if t.Tag != token.Semicolon && t.Tag != token.Newline {
		return nil, p.unexpected(t, "to terminate expression statement")
	}
	return &ast.ExpressionStatement{Loc: setLocation(&token.Location{}, x.Location(), &t.Location), Expression: x}, nil
}

func parseIf(p *Parser) (ast.Statement, error) {
	return nil, nil
}
