package scm

import (
	"testing"
)

type NumberTestCase struct {
	a        Number
	b        Number
	expected Number
}

func TestAdd(t *testing.T) {
	var tests = []NumberTestCase{
		{a: IntNum(1), b: IntNum(2), expected: IntNum(3)},
		{a: IntNum(1), b: RealNum(1.1), expected: RealNum(2.1)},
		{a: IntNum(1), b: MakeRatnum(IntNum(1), IntNum(3)), expected: MakeRatnum(IntNum(4), IntNum(3))},
		{a: IntNum(1), b: CompNum(2 + 1.1i), expected: CompNum(3 + 1.1i)},
		{a: RealNum(2.2), b: RealNum(1.1), expected: RealNum(3.3)},
		{a: RealNum(1.1), b: IntNum(2), expected: RealNum(3.1)},
		{a: RealNum(1.1), b: MakeRatnum(IntNum(1), IntNum(2)), expected: RealNum(1.6)},
		{a: RealNum(1.1), b: CompNum(2 + 1.1i), expected: CompNum(3.1 + 1.1i)},
		{a: MakeRatnum(IntNum(1), IntNum(3)), b: IntNum(3), expected: MakeRatnum(IntNum(10), IntNum(3))},
		{a: MakeRatnum(IntNum(1), IntNum(2)), b: RealNum(0.5), expected: RealNum(1.0)},
		{a: MakeRatnum(IntNum(1), IntNum(3)), b: MakeRatnum(IntNum(1), IntNum(6)), expected: MakeRatnum(IntNum(1), IntNum(2))},
		{a: MakeRatnum(IntNum(1), IntNum(2)), b: CompNum(1 + 1i), expected: CompNum(1.5 + 1i)},
		{a: CompNum(1 + 1i), b: IntNum(2), expected: CompNum(3 + 1i)},
		{a: CompNum(1 + 1i), b: RealNum(2.2), expected: CompNum(3.2 + 1i)},
		{a: CompNum(1 + 1i), b: MakeRatnum(IntNum(1), IntNum(2)), expected: CompNum(1.5 + 1i)},
		{a: CompNum(1 + 1i), b: CompNum(2 + 1.1i), expected: CompNum(3 + 2.1i)},
	}

	for _, tt := range tests {
		actual := tt.a.Add(tt.b)
		if !EqNum(actual, tt.expected) {
			t.Errorf("%v + %v => actual %v, expected %v", tt.a, tt.b, actual, tt.expected)
		}
	}
}

func TestSub(t *testing.T) {
	var tests = []NumberTestCase{
		{a: IntNum(1), b: IntNum(2), expected: IntNum(-1)},
		{a: IntNum(1), b: MakeRatnum(IntNum(2), IntNum(3)), expected: MakeRatnum(IntNum(1), IntNum(3))},
		{a: IntNum(1), b: RealNum(1.1), expected: RealNum(-0.1)},
		{a: IntNum(1), b: CompNum(2 + 1.1i), expected: CompNum(-1 - 1.1i)},
		{a: RealNum(2.2), b: RealNum(1.1), expected: RealNum(1.1)},
		{a: RealNum(1.1), b: IntNum(2), expected: RealNum(-0.9)},
		{a: RealNum(2.0), b: MakeRatnum(IntNum(1), IntNum(2)), expected: RealNum(1.5)},
		{a: RealNum(1.1), b: CompNum(2 + 1.1i), expected: CompNum(-0.9 - 1.1i)},
		{a: MakeRatnum(IntNum(2), IntNum(3)), b: IntNum(1), expected: MakeRatnum(IntNum(-1), IntNum(3))},
		{a: MakeRatnum(IntNum(2), IntNum(3)), b: MakeRatnum(IntNum(1), IntNum(6)), expected: MakeRatnum(IntNum(1), IntNum(2))},
		{a: MakeRatnum(IntNum(2), IntNum(5)), b: RealNum(0.5), expected: RealNum(-0.1)},
		{a: MakeRatnum(IntNum(2), IntNum(5)), b: CompNum(1 + 1i), expected: CompNum(-0.6 - 1i)},
		{a: CompNum(1 + 1i), b: IntNum(2), expected: CompNum(-1 + 1i)},
		{a: CompNum(1 + 1i), b: MakeRatnum(IntNum(1), IntNum(2)), expected: CompNum(0.5 + 1i)},
		{a: CompNum(1 + 1i), b: RealNum(2.2), expected: CompNum(-1.2 + 1i)},
		{a: CompNum(1 + 1i), b: CompNum(2 + 1.1i), expected: CompNum(-1 - 0.1i)},
	}

	for _, tt := range tests {
		actual := tt.a.Sub(tt.b)
		if !eq(actual, tt.expected) {
			t.Errorf("%v - %v => actual %v, expected %v", tt.a, tt.b, actual, tt.expected)
		}
	}
}

func TestMul(t *testing.T) {
	var tests = []NumberTestCase{
		{a: IntNum(1), b: IntNum(2), expected: IntNum(2)},
		{a: IntNum(2), b: MakeRatnum(IntNum(5), IntNum(4)), expected: MakeRatnum(IntNum(5), IntNum(2))},
		{a: IntNum(2), b: RealNum(1.1), expected: RealNum(2.2)},
		{a: IntNum(2), b: CompNum(1 - 1i), expected: CompNum(2 - 2i)},
		{a: RealNum(1.0), b: RealNum(2.1), expected: RealNum(2.1)},
		{a: RealNum(1.1), b: IntNum(2), expected: RealNum(2.2)},
		{a: RealNum(1.0), b: MakeRatnum(IntNum(1), IntNum(4)), expected: RealNum(0.25)},
		{a: RealNum(1.2), b: CompNum(1 - 1i), expected: CompNum(1.2 - 1.2i)},
		{a: MakeRatnum(IntNum(1), IntNum(2)), b: IntNum(3), expected: MakeRatnum(IntNum(3), IntNum(2))},
		{a: MakeRatnum(IntNum(1), IntNum(2)), b: RealNum(1.2), expected: RealNum(0.6)},
		{a: MakeRatnum(IntNum(1), IntNum(2)), b: MakeRatnum(IntNum(2), IntNum(3)), expected: MakeRatnum(IntNum(1), IntNum(3))},
		{a: MakeRatnum(IntNum(1), IntNum(2)), b: CompNum(2.5 + 1i), expected: CompNum(1.25 + 0.5i)},
		{a: CompNum(1 + 1i), b: IntNum(2), expected: CompNum(2 + 2i)},
		{a: CompNum(1 + 1i), b: MakeRatnum(IntNum(1), IntNum(4)), expected: CompNum(0.25 + 0.25i)},
		{a: CompNum(1 + 1i), b: RealNum(1.1), expected: CompNum(1.1 + 1.1i)},
		{a: CompNum(1 + 1i), b: CompNum(1 - 1i), expected: CompNum(2 + 0i)},
	}

	for _, tt := range tests {
		actual := tt.a.Mul(tt.b)
		if !eq(actual, tt.expected) {
			t.Errorf("%v * %v => actual %v, expected %v", tt.a, tt.b, actual, tt.expected)
		}
	}
}

func TestDiv(t *testing.T) {
	var tests = []NumberTestCase{
		{a: IntNum(1), b: IntNum(2), expected: MakeRatnum(IntNum(1), IntNum(2))},
		{a: IntNum(1), b: RealNum(0.5), expected: RealNum(2.0)},
		{a: IntNum(2), b: MakeRatnum(IntNum(3), IntNum(2)), expected: MakeRatnum(IntNum(4), IntNum(3))},
		{a: IntNum(2), b: CompNum(1 + 1i), expected: CompNum(1 - 1i)},
		{a: MakeRatnum(IntNum(1), IntNum(2)), b: IntNum(3), expected: MakeRatnum(IntNum(1), IntNum(6))},
		{a: MakeRatnum(IntNum(1), IntNum(2)), b: RealNum(0.5), expected: RealNum(1.0)},
		{a: MakeRatnum(IntNum(1), IntNum(2)), b: MakeRatnum(IntNum(2), IntNum(3)), expected: MakeRatnum(IntNum(3), IntNum(4))},
		{a: MakeRatnum(IntNum(1), IntNum(2)), b: CompNum(1 + 1i), expected: CompNum(0.25 - 0.25i)},
		{a: RealNum(1.0), b: RealNum(0.5), expected: RealNum(2.0)},
		{a: RealNum(1.1), b: IntNum(2), expected: RealNum(0.55)},
		{a: RealNum(1.0), b: MakeRatnum(IntNum(1), IntNum(4)), expected: RealNum(4.0)},
		{a: RealNum(1.2), b: CompNum(1 - 1i), expected: CompNum(0.6 + 0.6i)},
		{a: CompNum(1 + 1i), b: IntNum(2), expected: CompNum(0.5 + 0.5i)},
		{a: CompNum(1 + 1i), b: MakeRatnum(IntNum(1), IntNum(4)), expected: CompNum(4 + 4i)},
		{a: CompNum(1 + 1i), b: RealNum(0.2), expected: CompNum(5.0 + 5.0i)},
		{a: CompNum(1 + 1i), b: CompNum(1 + 1i), expected: CompNum(1 + 0i)},
	}

	for _, tt := range tests {
		actual := tt.a.Div(tt.b)
		if !eq(actual, tt.expected) {
			t.Errorf("%v / %v => actual %v, expected %v", tt.a, tt.b, actual, tt.expected)
		}
	}
}

func TestMakeratnum(t *testing.T) {
	a := MakeRatnum(IntNum(1), IntNum(2))
	if !eqReal(RealNum(0.5), a.ToReal()) {
		t.Error("failed to make ratnum with positive integers")
	}
	b := MakeRatnum(IntNum(-1), IntNum(2))
	if !eqReal(RealNum(-0.5), b.ToReal()) {
		t.Error("failed to make ratnum with negative integer")
	}
	c := MakeRatnum(IntNum(1), IntNum(-2))
	if c.Denominator != IntNum(2) || c.Numerator != IntNum(-1) {
		t.Error("failed to make ratnum with negative denominator")
	}
}
