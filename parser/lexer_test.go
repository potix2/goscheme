package parser

import (
	"testing"
)

type LexerTestCase struct {
	input    string
	expected []Token
}

func TestSingleLineScan(t *testing.T) {
	var tests = []LexerTestCase{
		{input: "(exit )", expected: []Token{
			Token{int('('), "(", Position{1, 1}},
			Token{IDENT, "exit", Position{1, 2}},
			Token{int(')'), ")", Position{1, 7}},
			Token{-1, "", Position{1, 8}},
		},
		},
		{input: "(+ b 10)", expected: []Token{
			Token{int('('), "(", Position{1, 1}},
			Token{IDENT, "+", Position{1, 2}},
			Token{IDENT, "b", Position{1, 4}},
			Token{NUMBER, "10", Position{1, 6}},
			Token{int(')'), ")", Position{1, 8}},
			Token{-1, "", Position{1, 9}},
		},
		},
		{input: "(++ ..> ->)", expected: []Token{
			Token{int('('), "(", Position{1, 1}},
			Token{IDENT, "++", Position{1, 2}},
			Token{IDENT, "..>", Position{1, 5}},
			Token{IDENT, "->", Position{1, 9}},
			Token{int(')'), ")", Position{1, 11}},
			Token{-1, "", Position{1, 12}},
		},
		},
		{input: "#t", expected: []Token{
			Token{BOOLEAN, "true", Position{1, 1}},
			Token{-1, "", Position{1, 3}},
		},
		},
		{input: "'a", expected: []Token{
			Token{int('\''), "'", Position{1, 1}},
			Token{IDENT, "a", Position{1, 2}},
			Token{-1, "", Position{1, 3}},
		},
		},
		{input: "#false", expected: []Token{
			Token{BOOLEAN, "false", Position{1, 1}},
			Token{-1, "", Position{1, 7}},
		},
		},
		{input: "(+ b ;'a\n10)", expected: []Token{
			Token{int('('), "(", Position{1, 1}},
			Token{IDENT, "+", Position{1, 2}},
			Token{IDENT, "b", Position{1, 4}},
			Token{NUMBER, "10", Position{2, 1}},
			Token{int(')'), ")", Position{2, 3}},
			Token{-1, "", Position{2, 4}},
		},
		},
		{input: ".12", expected: []Token{
			Token{NUMBER, ".12", Position{1, 1}},
			Token{-1, "", Position{1, 4}},
		},
		},
		{input: "1.2", expected: []Token{
			Token{NUMBER, "1.2", Position{1, 1}},
			Token{-1, "", Position{1, 4}},
		},
		},
		{input: "1/2", expected: []Token{
			Token{NUMBER, "1/2", Position{1, 1}},
			Token{-1, "", Position{1, 4}},
		},
		},
		{input: "\"abc\"", expected: []Token{
			Token{STRING, "abc", Position{1, 1}},
			Token{-1, "", Position{1, 6}},
		},
		},
		{input: "\"ab\\\"c\"", expected: []Token{
			Token{STRING, "ab\"c", Position{1, 1}},
			Token{-1, "", Position{1, 8}},
		},
		},
		{input: "\"\\x60;\"", expected: []Token{
			Token{STRING, "`", Position{1, 1}},
			Token{-1, "", Position{1, 8}},
		},
		},
		{input: "1\n2", expected: []Token{
			Token{NUMBER, "1", Position{1, 1}},
			Token{NUMBER, "2", Position{2, 1}},
			Token{-1, "", Position{2, 2}},
		},
		},
	}

	for _, tt := range tests {
		s := Scanner{}
		s.Init(tt.input)
		for _, e := range tt.expected {
			tok, lit, pos, err := s.Scan()
			if err != nil {
				t.Errorf("error: %v", err)
			}
			if tok != e.Tok {
				t.Errorf("token: actual %v, expected %v (line=%d col=%d)\n", tok, e.Tok, pos.Line, pos.Column)
			}
			if lit != e.Lit {
				t.Errorf("literal: actual %v, expected %v (line=%d col=%d)\n", lit, e.Lit, pos.Line, pos.Column)
			}
			if pos != e.Pos {
				t.Errorf("pos: actual %v, expected %v\n", pos, e.Pos)
			}
		}
	}
}
