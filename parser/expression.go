package parser

import (
	"fmt"
	"strconv"

	"github.com/arikui1911/goore/ast"
	"github.com/arikui1911/goore/token"
)

type precedence int

const (
	lowestPrecedence precedence = iota
	equalityPrecedence
	comparePrecedence
	additivePrecedence
	multivePrecedence
	prefixPrecedence
	callPrecedence
	highestPrecedence
)

type prefixedParser func(*Parser) (ast.Expression, error)

var prefixedParsers map[token.TokenTag]prefixedParser

type infixedParser func(*Parser, ast.Expression) (ast.Expression, error)

var infixedParsers map[token.TokenTag]infixedParser

func init() {
	prefixedParsers = map[token.TokenTag]prefixedParser{
		token.Identifier:    parseIdentifier,
		token.Nil:           parseNilLiteral,
		token.True:          parseBoolLiteral,
		token.False:         parseBoolLiteral,
		token.IntLiteral:    parseIntLiteral,
		token.FloatLiteral:  parseFloatLiteral,
		token.StringLiteral: parseStringLiteral,
		token.Add:           parsePrefixed,
		token.Sub:           parsePrefixed,
		token.Bang:          parsePrefixed,
		token.LeftParen:     parseParenExpr,
		token.LeftBracket:   parseArrayLiteral,
		token.LeftBrace:     parseHashLiteral,
		token.Arrow:         parseFunctionLiteral,
		token.If:            parseIf,
	}
	infixedParsers = map[token.TokenTag]infixedParser{
		token.Eq:          parseInfixed,
		token.Ne:          parseInfixed,
		token.Ge:          parseInfixed,
		token.Le:          parseInfixed,
		token.Gt:          parseInfixed,
		token.Lt:          parseInfixed,
		token.Add:         parseInfixed,
		token.Sub:         parseInfixed,
		token.Mul:         parseInfixed,
		token.Div:         parseInfixed,
		token.Mod:         parseInfixed,
		token.LeftParen:   parseCall,
		token.LeftBracket: parseKeyAccess,
		token.Let:         parseLet,
		token.LetAdd:      parseLet,
		token.LetSub:      parseLet,
		token.LetMul:      parseLet,
		token.LetDiv:      parseLet,
	}
}

var precedences = map[token.TokenTag]precedence{
	token.Eq:          equalityPrecedence,
	token.Ne:          equalityPrecedence,
	token.Ge:          comparePrecedence,
	token.Le:          comparePrecedence,
	token.Gt:          comparePrecedence,
	token.Lt:          comparePrecedence,
	token.Add:         additivePrecedence,
	token.Sub:         additivePrecedence,
	token.Mul:         multivePrecedence,
	token.Div:         multivePrecedence,
	token.Mod:         multivePrecedence,
	token.LeftParen:   callPrecedence,
	token.LeftBracket: callPrecedence,
	token.Let:         highestPrecedence,
	token.LetAdd:      highestPrecedence,
	token.LetSub:      highestPrecedence,
	token.LetMul:      highestPrecedence,
	token.LetDiv:      highestPrecedence,
	token.LetMod:      highestPrecedence,
}

func parseExpression(p *Parser, prec precedence) (ast.Expression, error) {
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

	for {
		t, err = p.peekToken()
		if err != nil {
			return nil, err
		}
		p.pushBack(t)

		np, ok := precedences[t.Tag]
		if !ok {
			break
		}
		if prec != highestPrecedence && prec >= np {
			break
		}
		left, err = infixedParsers[t.Tag](p, left)
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

func parseIdentifier(p *Parser) (ast.Expression, error) {
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	if t.Tag != token.Identifier {
		return nil, p.unexpected(t, "expect identifier")
	}
	return &ast.Identifier{Loc: &t.Location, Name: t.Value}, nil
}

func parseNilLiteral(p *Parser) (ast.Expression, error) {
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	return &ast.NilLiteral{Loc: &t.Location}, nil
}

func parseBoolLiteral(p *Parser) (ast.Expression, error) {
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	return &ast.BoolLiteral{Loc: &t.Location, Value: t.Tag == token.True}, nil
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

func parseStringLiteral(p *Parser) (ast.Expression, error) {
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	return &ast.StringLiteral{Loc: &t.Location, Value: t.Value}, nil
}

var prefixOperators = map[token.TokenTag]ast.Operation{
	token.Add:  ast.Plus,
	token.Sub:  ast.Minus,
	token.Bang: ast.Not,
}

func parsePrefixed(p *Parser) (ast.Expression, error) {
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	x, err := parseExpression(p, prefixPrecedence)
	if err != nil {
		return nil, err
	}
	return &ast.PrefixExpression{
		Loc:      setLocation(nil, &t.Location, x.Location()),
		Operator: prefixOperators[t.Tag],
		Right:    x,
	}, nil
}

func parseParenExpr(p *Parser) (ast.Expression, error) {
	_, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	x, err := parseExpression(p, lowestPrecedence)
	if err != nil {
		return nil, err
	}
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	if t.Tag != token.RightParen {
		return nil, p.unexpected(t, "expect right paren")
	}
	return x, nil
}

func parseArrayLiteral(p *Parser) (ast.Expression, error) {
	lb, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	elems, rb, err := parseCommaList(p, token.RightBracket, parseArrayElement)
	if err != nil {
		return nil, err
	}
	return &ast.ArrayLiteral{
		Loc:      setLocation(nil, &lb.Location, &rb.Location),
		Elements: elems,
	}, nil
}

func parseArrayElement(p *Parser) (ast.Expression, error) {
	return parseExpression(p, lowestPrecedence)
}

func parseHashLiteral(p *Parser) (ast.Expression, error) {
	lb, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	pairs, rb, err := parseCommaList(p, token.RightBrace, parseHashEntry)
	if err != nil {
		return nil, err
	}
	return &ast.HashLiteral{
		Loc:   setLocation(nil, &lb.Location, &rb.Location),
		Pairs: pairs,
	}, nil
}

func parseHashEntry(p *Parser) (*ast.HashEntry, error) {
	k, err := parseExpression(p, lowestPrecedence)
	if err != nil {
		return nil, err
	}
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	if t.Tag != token.Colon {
		return nil, p.unexpected(t, "expect colon to delimitting key and value")
	}
	v, err := parseExpression(p, lowestPrecedence)
	if err != nil {
		return nil, err
	}
	return &ast.HashEntry{
		Loc:   setLocation(nil, k.Location(), v.Location()),
		Key:   k,
		Value: v,
	}, nil
}

func parseFunctionLiteral(p *Parser) (ast.Expression, error) {
	arrow, err := p.nextToken()
	if err != nil {
		return nil, err
	}

	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	var params []*ast.Identifier
	if t.Tag == token.LeftParen {
		params, _, err = parseCommaList(p, token.RightParen, parseParameter)
		if err != nil {
			return nil, err
		}
	} else {
		p.pushBack(t)
		params = []*ast.Identifier{}
	}

	stmts, rb, err := parseBlock(p)
	if err != nil {
		return nil, err
	}

	return &ast.FunctionLiteral{
		Loc:        setLocation(nil, &arrow.Location, &rb.Location),
		Parameters: params,
		Statements: stmts,
	}, nil
}

func parseParameter(p *Parser) (*ast.Identifier, error) {
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	if t.Tag != token.Identifier {
		return nil, p.unexpected(t, "expect identifier as parameter name")
	}
	return &ast.Identifier{Loc: &t.Location, Name: t.Value}, nil
}

func parseIf(p *Parser) (ast.Expression, error) {
	kw, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	switch kw.Tag {
	case token.If, token.Elsif:
		test, err := parseExpression(p, lowestPrecedence)
		if err != nil {
			return nil, err
		}
		body, rb, err := parseBlock(p)
		if err != nil {
			return nil, err
		}
		alt, err := parseIf(p)
		if err != nil {
			return nil, err
		}
		loc := setLocation(nil, &kw.Location, &rb.Location)
		if alt != nil {
			loc = setLocation(loc, nil, alt.Location())
		}
		return &ast.If{Loc: loc, Test: test, Body: body, Alt: alt}, nil
	case token.Else:
		body, rb, err := parseBlock(p)
		if err != nil {
			return nil, err
		}
		return &ast.Else{Loc: setLocation(nil, &kw.Location, &rb.Location), Body: body}, nil
	default:
		p.pushBack(kw)
		return nil, nil
	}
}

var infixOperators = map[token.TokenTag]ast.Operation{
	token.Eq:  ast.Eq,
	token.Ne:  ast.Ne,
	token.Le:  ast.Le,
	token.Ge:  ast.Ge,
	token.Lt:  ast.Lt,
	token.Gt:  ast.Gt,
	token.Add: ast.Add,
	token.Sub: ast.Sub,
	token.Mul: ast.Mul,
	token.Div: ast.Div,
	token.Mod: ast.Mod,
}

func parseInfixed(p *Parser, left ast.Expression) (ast.Expression, error) {
	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	right, err := parseExpression(p, precedences[t.Tag])
	if err != nil {
		return nil, err
	}
	return &ast.InfixExpression{
		Loc:      setLocation(nil, left.Location(), right.Location()),
		Operator: infixOperators[t.Tag],
		Left:     left,
		Right:    right,
	}, nil
}

func parseCall(p *Parser, fn ast.Expression) (ast.Expression, error) {
	_, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	args, rp, err := parseCommaList(p, token.RightParen, parseArgument)
	if err != nil {
		return nil, err
	}
	return &ast.Call{
		Loc:       setLocation(nil, fn.Location(), &rp.Location),
		Function:  fn,
		Arguments: args,
	}, nil
}

func parseArgument(p *Parser) (ast.Expression, error) {
	return parseExpression(p, lowestPrecedence)
}

func parseKeyAccess(p *Parser, c ast.Expression) (ast.Expression, error) {
	_, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	k, err := parseExpression(p, lowestPrecedence)
	if err != nil {
		return nil, err
	}

	t, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	if t.Tag == token.Comma || t.Tag == token.Newline {
		t, err = p.nextToken()
		if err != nil {
			return nil, err
		}
	}
	if t.Tag != token.RightBracket {
		return nil, p.unexpected(t, "expect right bracket")
	}

	return &ast.KeyAccess{
		Loc:       setLocation(nil, c.Location(), &t.Location),
		Container: c,
		Key:       k,
	}, nil
}

var selfLetOperators = map[token.TokenTag]ast.Operation{
	token.LetAdd: ast.Add,
	token.LetSub: ast.Sub,
	token.LetMul: ast.Mul,
	token.LetDiv: ast.Div,
	token.LetMod: ast.Mod,
}

func parseLet(p *Parser, left ast.Expression) (ast.Expression, error) {
	var isLet bool
	switch left.(type) {
	case *ast.Identifier:
		isLet = true
	case *ast.KeyAccess:
		isLet = false
	default:
		return nil, fmt.Errorf("%s:%s: invalid let left part", p.fileName, left.Location())
	}

	let, err := p.nextToken()
	if err != nil {
		return nil, err
	}

	right, err := parseExpression(p, highestPrecedence)
	if err != nil {
		return nil, err
	}

	if op, ok := selfLetOperators[let.Tag]; ok {
		right = &ast.InfixExpression{
			Loc:      setLocation(nil, left.Location(), right.Location()),
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}

	if isLet {
		return &ast.Let{
			Loc:   setLocation(nil, left.Location(), right.Location()),
			Left:  left.(*ast.Identifier),
			Right: right,
		}, nil
	}

	return &ast.KeyAssign{
		Loc:   setLocation(nil, left.Location(), right.Location()),
		Left:  left.(*ast.KeyAccess),
		Right: right,
	}, nil
}

func parseCommaList[T ast.Expression](p *Parser, term token.TokenTag, elementParser func(*Parser) (T, error)) ([]T, token.Token, error) {
	// empty?
	t, err := p.nextToken()
	if err != nil {
		return nil, token.Token{}, err
	}
	if t.Tag == term {
		return []T{}, t, nil
	}
	p.pushBack(t)

	// first element
	e, err := elementParser(p)
	if err != nil {
		return nil, token.Token{}, err
	}
	list := []T{e}

	for {
		t, err := p.nextToken()
		if err != nil {
			return nil, token.Token{}, err
		}

		switch t.Tag {
		case token.Newline:
			// auto newline and term
			nt, err := p.nextToken()
			if err != nil {
				return nil, token.Token{}, err
			}
			if nt.Tag == term {
				return list, nt, nil
			}

			// カンマが欠けていたので改行トークンが入ってしまったとして、次の要素へ
			p.pushBack(nt)
			p.addError(fmt.Errorf("%s:%s: missing comma for delimiter", p.fileName, nt.Location))
		case term:
			// term without auto newline (e.g. [123])
			return list, t, nil
		case token.Comma:
			// consume last extra comma
			t, err := p.nextToken()
			if err != nil {
				return nil, token.Token{}, err
			}
			if t.Tag == term {
				return list, t, nil
			}
			p.pushBack(t)
		default:
			// カンマが欠けてると仮定して、次の要素を読みにいく
			p.pushBack(t)
			p.addError(fmt.Errorf("%s:%s: missing comma for delimiter", p.fileName, t.Location))
		}

		// next element
		e, err := elementParser(p)
		if err != nil {
			return nil, token.Token{}, err
		}
		list = append(list, e)
	}
}
