package parser

import (
	"errors"

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
	return &ast.Program{Loc: loc, FileName: p.fileName, Statements: stmts, Err: errors.Join(p.errs...)}, nil
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

func parseBlock(p *Parser) ([]ast.Statement, token.Token, error) {
	t, err := p.nextToken()
	if err != nil {
		return nil, token.Token{}, err
	}
	if t.Tag != token.LeftBrace {
		return nil, token.Token{}, p.unexpected(t, "expect left brace to begin block")
	}
	return parseStatements(p, token.RightBrace)
}

type statementParser func(*Parser) (ast.Statement, error)

var statementParsers map[token.TokenTag]statementParser

func init() {
	statementParsers = map[token.TokenTag]statementParser{
		token.Semicolon: parseEmpty,
		token.Def:       parseDef,
		token.While:     parseWhile,
		token.Break:     parseBreak,
		token.Continue:  parseContinue,
		token.Return:    parseReturn,
	}
}

func parseStatement(p *Parser) (ast.Statement, error) {
	t, err := p.peekToken()
	if err != nil {
		return nil, err
	}
	fn, ok := statementParsers[t.Tag]
	if !ok {
		fn = parseExpressionStatement
	}
	s, err := fn(p)
	if err != nil {
		p.addError(err)
		return parseInvalidStatement(p, t, err)
	}
	return s, nil
}

func parseInvalidStatement(p *Parser, first token.Token, cause error) (*ast.InvalidStatement, error) {
	for {
		t, err := p.nextToken()
		if err != nil {
			return nil, err
		}
		if t.Tag != token.Newline && t.Tag != token.Semicolon && t.Tag != token.EOF {
			continue
		}
		return &ast.InvalidStatement{
			Loc: setLocation(nil, &first.Location, &t.Location),
			Err: cause,
		}, nil
	}
}

func parseEmpty(*Parser) (ast.Statement, error) { return nil, nil }

func parseDef(p *Parser) (ast.Statement, error) {
	kw, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	x, err := parseIdentifier(p)
	if err != nil {
		return nil, err
	}
	name := x.(*ast.Identifier)
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	if t.Tag == token.Newline || t.Tag == token.Semicolon {
		return &ast.Def{
			Loc:  setLocation(nil, &kw.Location, &t.Location),
			Name: name,
		}, nil
	}
	if t.Tag != token.Let {
		return nil, p.unexpected(t, "expect '=', newline or semicolon")
	}
	x, err = parseExpression(p, lowestPrecedence)
	if err != nil {
		return nil, err
	}
	t, err = p.nextToken()
	if err != nil {
		return nil, err
	}
	if t.Tag != token.Newline && t.Tag != token.Semicolon {
		return nil, p.unexpected(t, "expect newline or semicolon to terminate def statement")
	}
	return &ast.Def{
		Loc:  setLocation(nil, &kw.Location, &t.Location),
		Name: name,
		Init: x,
	}, nil
}

func parseWhile(p *Parser) (ast.Statement, error) {
	kw, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	cond, err := parseExpression(p, lowestPrecedence)
	if err != nil {
		return nil, err
	}
	stmts, rb, err := parseBlock(p)
	if err != nil {
		return nil, err
	}
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	if t.Tag != token.Newline {
		p.pushBack(t)
	}
	return &ast.While{
		Loc:  setLocation(nil, &kw.Location, &rb.Location),
		Cond: cond,
		Body: stmts,
	}, nil
}

func parseBreak(p *Parser) (ast.Statement, error) {
	kw, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	if t.Tag != token.Newline && t.Tag != token.Semicolon {
		return nil, p.unexpected(t, "expect newline or semicolon to terminate break statement")
	}
	return &ast.Break{Loc: setLocation(nil, &kw.Location, &t.Location)}, nil
}

func parseContinue(p *Parser) (ast.Statement, error) {
	kw, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	if t.Tag != token.Newline && t.Tag != token.Semicolon {
		return nil, p.unexpected(t, "expect newline or semicolon to terminate continue statement")
	}
	return &ast.Continue{Loc: setLocation(nil, &kw.Location, &t.Location)}, nil
}

func parseReturn(p *Parser) (ast.Statement, error) {
	kw, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	if t.Tag == token.Newline || t.Tag == token.Semicolon {
		return &ast.Return{Loc: setLocation(nil, &kw.Location, &t.Location)}, nil
	}
	p.pushBack(t)
	x, err := parseExpression(p, lowestPrecedence)
	if err != nil {
		return nil, err
	}
	t, err = p.nextToken()
	if err != nil {
		return nil, err
	}
	if t.Tag != token.Newline && t.Tag != token.Semicolon {
		return nil, p.unexpected(t, "expect newline or semicolon to terminate return statement")
	}
	return &ast.Return{Loc: setLocation(nil, &kw.Location, &t.Location), Expression: x}, nil
}

func parseExpressionStatement(p *Parser) (ast.Statement, error) {
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
