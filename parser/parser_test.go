package parser

import (
	"reflect"
	"testing"

	"github.com/potix2/goscheme/scm"
)

type ParserTestCase struct {
	input    string
	expected []scm.Object
}

func TestParser(t *testing.T) {
	var tests = []ParserTestCase{
		{input: "a", expected: []scm.Object{scm.IdentExpr{Lit: "a"}}},
		{input: "10", expected: []scm.Object{scm.IntNum(10)}},
		{input: "10 20", expected: []scm.Object{scm.IntNum(10), scm.IntNum(20)}},
		{input: "-10", expected: []scm.Object{scm.IntNum(-10)}},
		{input: "+10", expected: []scm.Object{scm.IntNum(10)}},
		{input: "2/3", expected: []scm.Object{scm.MakeRatnum(scm.IntNum(2), scm.IntNum(3))}},
		{input: "#true", expected: []scm.Object{scm.BooleanExpr{Lit: true}}},
		{input: "#f", expected: []scm.Object{scm.BooleanExpr{Lit: false}}},
		{input: "'a", expected: []scm.Object{scm.QuoteExpr{Datum: scm.IdentExpr{Lit: "a"}}}},
		{input: "'()", expected: []scm.Object{scm.QuoteExpr{Datum: scm.PairExpr{}}}},
		{input: "'(a b)", expected: []scm.Object{scm.QuoteExpr{Datum: scm.PairExpr{Car: scm.IdentExpr{Lit: "a"}, Cdr: scm.PairExpr{Car: scm.IdentExpr{Lit: "b"}, Cdr: scm.PairExpr{}}}}}},
		{input: "'(a . b)", expected: []scm.Object{scm.QuoteExpr{Datum: scm.PairExpr{Car: scm.IdentExpr{Lit: "a"}, Cdr: scm.IdentExpr{Lit: "b"}}}}},
		{input: "(a)", expected: []scm.Object{scm.AppExpr{Objs: []scm.Object{scm.IdentExpr{Lit: "a"}}}}},
		{input: "(a b c)", expected: []scm.Object{scm.AppExpr{Objs: []scm.Object{scm.IdentExpr{Lit: "a"}, scm.IdentExpr{Lit: "b"}, scm.IdentExpr{Lit: "c"}}}}},
		{input: "(+ 1 2)", expected: []scm.Object{scm.AppExpr{Objs: []scm.Object{scm.IdentExpr{Lit: "+"}, scm.IntNum(1), scm.IntNum(2)}}}},
		{input: "(* 1 2)", expected: []scm.Object{scm.AppExpr{Objs: []scm.Object{scm.IdentExpr{Lit: "*"}, scm.IntNum(1), scm.IntNum(2)}}}},
		{input: "((a b) c)", expected: []scm.Object{scm.AppExpr{Objs: []scm.Object{
			scm.AppExpr{Objs: []scm.Object{scm.IdentExpr{Lit: "a"}, scm.IdentExpr{Lit: "b"}}},
			scm.IdentExpr{Lit: "c"}}}}},
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
