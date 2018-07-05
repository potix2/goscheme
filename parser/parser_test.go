package parser

import (
	"reflect"
	"testing"

	"github.com/potix2/goscheme/ast"
)

type ParserTestCase struct {
	input    string
	expected []ast.Expr
}

func TestParser(t *testing.T) {
	var tests = []ParserTestCase{
		{input: "a", expected: []ast.Expr{ast.IdentExpr{Lit: "a"}}},
		{input: "10", expected: []ast.Expr{ast.IntNum(10)}},
		{input: "10 20", expected: []ast.Expr{ast.IntNum(10), ast.IntNum(20)}},
		{input: "-10", expected: []ast.Expr{ast.IntNum(-10)}},
		{input: "+10", expected: []ast.Expr{ast.IntNum(10)}},
		{input: "2/3", expected: []ast.Expr{ast.MakeRatnum(ast.IntNum(2), ast.IntNum(3))}},
		{input: "#true", expected: []ast.Expr{ast.BooleanExpr{Lit: true}}},
		{input: "#f", expected: []ast.Expr{ast.BooleanExpr{Lit: false}}},
		{input: "'a", expected: []ast.Expr{ast.QuoteExpr{Datum: ast.IdentExpr{Lit: "a"}}}},
		{input: "'()", expected: []ast.Expr{ast.QuoteExpr{Datum: ast.PairExpr{}}}},
		{input: "'(a b)", expected: []ast.Expr{ast.QuoteExpr{Datum: ast.PairExpr{Car: ast.IdentExpr{Lit: "a"}, Cdr: ast.PairExpr{Car: ast.IdentExpr{Lit: "b"}, Cdr: ast.PairExpr{}}}}}},
		{input: "'(a . b)", expected: []ast.Expr{ast.QuoteExpr{Datum: ast.PairExpr{Car: ast.IdentExpr{Lit: "a"}, Cdr: ast.IdentExpr{Lit: "b"}}}}},
		{input: "(a)", expected: []ast.Expr{ast.AppExpr{Exprs: []ast.Expr{ast.IdentExpr{Lit: "a"}}}}},
		{input: "(a b c)", expected: []ast.Expr{ast.AppExpr{Exprs: []ast.Expr{ast.IdentExpr{Lit: "a"}, ast.IdentExpr{Lit: "b"}, ast.IdentExpr{Lit: "c"}}}}},
		{input: "(+ 1 2)", expected: []ast.Expr{ast.AppExpr{Exprs: []ast.Expr{ast.IdentExpr{Lit: "+"}, ast.IntNum(1), ast.IntNum(2)}}}},
		{input: "(* 1 2)", expected: []ast.Expr{ast.AppExpr{Exprs: []ast.Expr{ast.IdentExpr{Lit: "*"}, ast.IntNum(1), ast.IntNum(2)}}}},
		{input: "((a b) c)", expected: []ast.Expr{ast.AppExpr{Exprs: []ast.Expr{
			ast.AppExpr{Exprs: []ast.Expr{ast.IdentExpr{Lit: "a"}, ast.IdentExpr{Lit: "b"}}},
			ast.IdentExpr{Lit: "c"}}}}},
	}

	t.SkipNow()
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
