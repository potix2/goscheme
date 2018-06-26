package parser

import (
	"reflect"
	"testing"

	"github.com/potix2/goscheme/ast"
)

type ParserTestCase struct {
	input    string
	expected ast.Expr
}

func TestParser(t *testing.T) {
	var tests = []ParserTestCase{
		{input: "a", expected: ast.IdentExpr{Lit: "a"}},
		{input: "10", expected: ast.Uint10Expr{Lit: 10}},
		{input: "-10", expected: ast.Uint10Expr{Lit: -10}},
		{input: "+10", expected: ast.Uint10Expr{Lit: 10}},
		{input: "(a)", expected: ast.AppExpr{Operator: ast.IdentExpr{Lit: "a"},
			Operands: []ast.Expr{}}},
		{input: "(a b c)", expected: ast.AppExpr{Operator: ast.IdentExpr{Lit: "a"},
			Operands: []ast.Expr{ast.IdentExpr{Lit: "b"}, ast.IdentExpr{Lit: "c"}}}},
		{input: "(+ 1 2)", expected: ast.AppExpr{Operator: ast.IdentExpr{Lit: "+"},
			Operands: []ast.Expr{ast.Uint10Expr{Lit: 1}, ast.Uint10Expr{Lit: 2}}}},
		{input: "(* 1 2)", expected: ast.AppExpr{Operator: ast.IdentExpr{Lit: "*"},
			Operands: []ast.Expr{ast.Uint10Expr{Lit: 1}, ast.Uint10Expr{Lit: 2}}}},
		{input: "((a b) c)", expected: ast.AppExpr{Operator: ast.AppExpr{Operator: ast.IdentExpr{Lit: "a"},
			Operands: []ast.Expr{ast.IdentExpr{Lit: "b"}}},
			Operands: []ast.Expr{ast.IdentExpr{Lit: "c"}}}},
	}

	for _, tt := range tests {
		s := Scanner{}
		s.Init(tt.input)
		actual, err := Parse(&s)
		if err != nil {
			t.Errorf("parse error: %v\n", err)
		}
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("actual: %+#v expected: %+#v\n", actual, tt.expected)
		}
	}
}
