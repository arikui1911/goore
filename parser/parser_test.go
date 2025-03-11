package parser_test

import (
	"testing"

	"github.com/arikui1911/goore/ast"
	"github.com/arikui1911/goore/parser"
)

func TestParseIdentifier(t *testing.T) {
	tree, err := parser.ParseString(`hoge`, "test.goore")
	if err != nil {
		t.Error(err)
		return
	}
	testOneExpression(t, tree, func(t *testing.T, x ast.Expression) {
		testIdentifier(t, x, "hoge")
	})
}

func TestParseIntLiteral(t *testing.T) {
	tree, err := parser.ParseString(`123`, "test.goore")
	if err != nil {
		t.Error(err)
		return
	}
	testOneExpression(t, tree, func(t *testing.T, x ast.Expression) {
		testIntLiteral(t, x, 123)
	})
}

func TestParseFloatLiteral(t *testing.T) {
	tree, err := parser.ParseString(`1.23`, "test.goore")
	if err != nil {
		t.Error(err)
		return
	}
	testOneExpression(t, tree, func(t *testing.T, x ast.Expression) {
		testFloatLiteral(t, x, 1.23)
	})
}

func testProgram(t *testing.T, tree ast.Node, nStmts int, inner func(*testing.T, []ast.Statement)) {
	pg, ok := tree.(*ast.Program)
	if !ok {
		t.Errorf("want <%T> got <%T>", &ast.Program{}, tree)
		return
	}
	if len(pg.Statements) != nStmts {
		t.Errorf("want <%d> got <%d>", nStmts, len(pg.Statements))
		return
	}
	inner(t, pg.Statements)
}

func testExpressionStatement(t *testing.T, s ast.Statement, inner func(*testing.T, ast.Expression)) {
	xs, ok := s.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("want <%T> got <%T>", &ast.ExpressionStatement{}, s)
		return
	}
	inner(t, xs.Expression)
}

func testOneExpression(t *testing.T, tree ast.Node, inner func(*testing.T, ast.Expression)) {
	testProgram(t, tree, 1, func(t *testing.T, stmts []ast.Statement) {
		testExpressionStatement(t, stmts[0], inner)
	})
}

func testIdentifier(t *testing.T, expr ast.Expression, name string) {
	x, ok := expr.(*ast.Identifier)
	if !ok {
		t.Errorf("want <%T> got <%T>", &ast.Identifier{}, expr)
		return
	}
	if x.Name != name {
		t.Errorf("want <%#v> got <%#v>", name, x.Name)
		return
	}
}

func testIntLiteral(t *testing.T, expr ast.Expression, val int) {
	x, ok := expr.(*ast.IntLiteral)
	if !ok {
		t.Errorf("want <%T> got <%T>", &ast.IntLiteral{}, expr)
		return
	}
	if x.Value != val {
		t.Errorf("want <%d> got <%d>", val, x.Value)
		return
	}
}

func testFloatLiteral(t *testing.T, expr ast.Expression, val float64) {
	x, ok := expr.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("want <%T> got <%T>", &ast.FloatLiteral{}, expr)
		return
	}
	if x.Value != val {
		t.Errorf("want <%f> got <%f>", val, x.Value)
		return
	}
}
