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

func TestParseNilLiteral(t *testing.T) {
	tree, err := parser.ParseString(`nil`, "test.goore")
	if err != nil {
		t.Error(err)
		return
	}
	testOneExpression(t, tree, func(t *testing.T, x ast.Expression) {
		testNilLiteral(t, x)
	})
}

func TestParseTrueLiteral(t *testing.T) {
	tree, err := parser.ParseString(`true`, "test.goore")
	if err != nil {
		t.Error(err)
		return
	}
	testOneExpression(t, tree, func(t *testing.T, x ast.Expression) {
		testBoolLiteral(t, x, true)
	})
}

func TestParseFalseLiteral(t *testing.T) {
	tree, err := parser.ParseString(`false`, "test.goore")
	if err != nil {
		t.Error(err)
		return
	}
	testOneExpression(t, tree, func(t *testing.T, x ast.Expression) {
		testBoolLiteral(t, x, false)
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

func TestParseStringLiteral(t *testing.T) {
	tree, err := parser.ParseString(`"Hello."`, "test.goore")
	if err != nil {
		t.Error(err)
		return
	}
	testOneExpression(t, tree, func(t *testing.T, x ast.Expression) {
		testStringLiteral(t, x, "Hello.")
	})
}

func TestParsePrefixExpressions(t *testing.T) {
	table := []struct {
		name string
		src  string
		op   ast.Operation
		val  int
	}{
		{"plus", `+123`, ast.Plus, 123},
		{"minus", `-123`, ast.Minus, 123},
		{"not", `!123`, ast.Not, 123},
	}

	for _, d := range table {
		t.Run(d.name, func(t *testing.T) {
			tree, err := parser.ParseString(d.src, "test.goore")
			if err != nil {
				t.Error(err)
				return
			}
			testOneExpression(t, tree, func(t *testing.T, x ast.Expression) {
				testPrefixExpression(t, x, d.op, d.val)
			})
		})
	}
}

func TestParseParened(t *testing.T) {
	tree, err := parser.ParseString(`(123)`, "test.goore")
	if err != nil {
		t.Error(err)
		return
	}
	testOneExpression(t, tree, func(t *testing.T, x ast.Expression) {
		testIntLiteral(t, x, 123)
	})
}

func TestParseInfixExpressions(t *testing.T) {
	table := []struct {
		name  string
		src   string
		op    ast.Operation
		left  int
		right int
	}{
		{"Eq", `1 == 2`, ast.Eq, 1, 2},
		{"Ne", `1 != 2`, ast.Ne, 1, 2},
		{"Le", `1 <= 2`, ast.Le, 1, 2},
		{"Ge", `1 >= 2`, ast.Ge, 1, 2},
		{"Lt", `1 < 2`, ast.Lt, 1, 2},
		{"Gt", `1 > 2`, ast.Gt, 1, 2},
		{"Add", `1 + 2`, ast.Add, 1, 2},
		{"Sub", `1 - 2`, ast.Sub, 1, 2},
		{"Mul", `1 * 2`, ast.Mul, 1, 2},
		{"Div", `1 / 2`, ast.Div, 1, 2},
		{"Mod", `1 % 2`, ast.Mod, 1, 2},
	}

	for _, d := range table {
		t.Run(d.name, func(t *testing.T) {
			tree, err := parser.ParseString(d.src, "test.goore")
			if err != nil {
				t.Error(err)
				return
			}
			testOneExpression(t, tree, func(t *testing.T, x ast.Expression) {
				testInfixExpression(t, x, d.op, d.left, d.right)
			})
		})
	}
}

func TestParseArrayLiterals(t *testing.T) {
	table := []struct {
		name string
		src  string
		vals []int
	}{
		{"empty", `[]`, []int{}},
		{"one element", `[123]`, []int{123}},
		{"several elements", `[1, 2, 3]`, []int{1, 2, 3}},
	}

	for _, d := range table {
		t.Run(d.name, func(t *testing.T) {
			tree, err := parser.ParseString(d.src, "test.goore")
			if err != nil {
				t.Error(err)
				return
			}
			testOneExpression(t, tree, func(t *testing.T, x ast.Expression) {
				testArrayLiteral(t, x, d.vals)
			})
		})
	}
}

func testArrayLiteral(t *testing.T, x ast.Expression, vals []int) {
	a, ok := x.(*ast.ArrayLiteral)
	if !ok {
		t.Errorf("want <%T> got <%T>", &ast.ArrayLiteral{}, x)
		return
	}
	if len(a.Elements) != len(vals) {
		t.Errorf("want <%d> got <%d>", len(vals), len(a.Elements))
		return
	}
	for i, e := range a.Elements {
		testIntLiteral(t, e, vals[i])
	}
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

func testNilLiteral(t *testing.T, expr ast.Expression) {
	_, ok := expr.(*ast.NilLiteral)
	if !ok {
		t.Errorf("want <%T> got <%T>", &ast.NilLiteral{}, expr)
		return
	}
}

func testBoolLiteral(t *testing.T, expr ast.Expression, val bool) {
	x, ok := expr.(*ast.BoolLiteral)
	if !ok {
		t.Errorf("want <%T> got <%T>", &ast.BoolLiteral{}, expr)
		return
	}
	if x.Value != val {
		t.Errorf("want <%v> got <%v>", val, x.Value)
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

func testStringLiteral(t *testing.T, expr ast.Expression, val string) {
	x, ok := expr.(*ast.StringLiteral)
	if !ok {
		t.Errorf("want <%T> got <%T>", &ast.StringLiteral{}, expr)
		return
	}
	if x.Value != val {
		t.Errorf("want <%#v> got <%#v>", val, x.Value)
		return
	}
}

func testPrefixExpression(t *testing.T, expr ast.Expression, op ast.Operation, val int) {
	x, ok := expr.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("want <%T> got <%T>", &ast.PrefixExpression{}, expr)
		return
	}
	if x.Operator != op {
		t.Errorf("want <%v> got <%v>", op, x.Operator)
		return
	}
	testIntLiteral(t, x.Right, val)
}

func testInfixExpression(t *testing.T, expr ast.Expression, op ast.Operation, l int, r int) {
	x, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Errorf("want <%T> got <%T>", &ast.InfixExpression{}, expr)
		return
	}
	if x.Operator != op {
		t.Errorf("want <%v> got <%v>", op, x.Operator)
		return
	}
	testIntLiteral(t, x.Left, l)
	testIntLiteral(t, x.Right, r)
}
