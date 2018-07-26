package scm

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	EPSILON RealNum = 0.00000001
)

type Number interface {
	Add(Number) Number
	Sub(Number) Number
	Mul(Number) Number
	Div(Number) Number
	Expr
}

type IntNum int

type RealNum float64

type RatNum struct {
	Numerator   Number
	Denominator Number
}

type CompNum complex128

type NAN struct{}

var NanValue = NAN{}

func (x NAN) Print(output io.Writer) {
	output.Write([]byte("NAN"))
}

func (x NAN) Add(y Number) Number {
	return NanValue
}

func (x NAN) Sub(y Number) Number {
	return NanValue
}

func (x NAN) Mul(y Number) Number {
	return NanValue
}

func (x NAN) Div(y Number) Number {
	return NanValue
}

/*
type BigNum struct {
	Sign   bool
	Size   int
	Values []int
}
*/

//IntNum
func (x IntNum) Print(output io.Writer) {
	output.Write([]byte(fmt.Sprintf("%d", x)))
}

func intToComplex(a IntNum) CompNum {
	return CompNum(complex(float64(a), 0))
}

func (x IntNum) Add(y Number) Number {
	switch y.(type) {
	case IntNum:
		return x + y.(IntNum)
	case RatNum:
		return MakeRatnum(x, IntNum(1)).Add(y.(RatNum))
	case RealNum:
		return RealNum(x) + y.(RealNum)
	case CompNum:
		return intToComplex(x) + y.(CompNum)
	}
	return NanValue
}

func (x IntNum) Sub(y Number) Number {
	return x.Add(Negate(y))
}

func (x IntNum) Mul(y Number) Number {
	switch y.(type) {
	case IntNum:
		return x * y.(IntNum)
	case RatNum:
		return MakeRatnum(x, IntNum(1)).Mul(y)
	case RealNum:
		return RealNum(x) * y.(RealNum)
	case CompNum:
		return intToComplex(x) * y.(CompNum)
	default:
		return IntNum(1) //TODO: return NAN
	}
}

func (x IntNum) Div(y Number) Number {
	return x.Mul(inverse(y))
}

//RealNum
func (x RealNum) Print(output io.Writer) {
	output.Write([]byte(fmt.Sprintf("%f", x)))
}

func realToComplex(a RealNum) CompNum {
	return CompNum(complex(a, 0))
}

func (x RealNum) Add(y Number) Number {
	switch y.(type) {
	case IntNum:
		return x + RealNum(y.(IntNum))
	case RatNum:
		return x + y.(RatNum).ToReal()
	case RealNum:
		return x + y.(RealNum)
	case CompNum:
		return realToComplex(x) + y.(CompNum)
	default:
		return IntNum(0) //TODO: return NAN
	}
}

func (x RealNum) Sub(y Number) Number {
	return x.Add(Negate(y))
}

func (x RealNum) Mul(y Number) Number {
	switch y.(type) {
	case IntNum:
		return x * RealNum(y.(IntNum))
	case RatNum:
		return x * y.(RatNum).ToReal()
	case RealNum:
		return x * y.(RealNum)
	case CompNum:
		return realToComplex(x) * y.(CompNum)
	default:
		return IntNum(1) //TODO: return NAN
	}
}

func (x RealNum) Div(y Number) Number {
	return x.Mul(inverse(y))
}

//RatNum
func (x RatNum) Print(output io.Writer) {
	output.Write([]byte(fmt.Sprintf("%d/%d", x.Numerator, x.Denominator)))
}

func ratnumAdd(a RatNum, b RatNum) Number {
	return MakeRatnum(a.Numerator.Mul(b.Denominator).Add(b.Numerator.Mul(a.Denominator)), a.Denominator.Mul(b.Denominator))
}

func ratnumMul(a RatNum, b RatNum) Number {
	return MakeRatnum(a.Numerator.Mul(b.Numerator), a.Denominator.Mul(b.Denominator))
}

func (x RatNum) ToReal() RealNum {
	return RealNum(float64(x.Numerator.(IntNum)) / float64(x.Denominator.(IntNum)))
}

func (x RatNum) Add(y Number) Number {
	switch y.(type) {
	case IntNum:
		return ratnumAdd(x, MakeRatnum(y.(IntNum), IntNum(1)))
	case RealNum:
		return x.ToReal().Add(y.(RealNum))
	case RatNum:
		return ratnumAdd(x, y.(RatNum))
	case CompNum:
		return y.(CompNum).Add(x.ToReal())
	}
	return MakeRatnum(IntNum(0), IntNum(0))
}

func (x RatNum) Sub(y Number) Number {
	return x.Add(Negate(y))
}

func (x RatNum) Mul(y Number) Number {
	switch y.(type) {
	case IntNum:
		return ratnumMul(x, MakeRatnum(y.(IntNum), IntNum(1)))
	case RealNum:
		return x.ToReal().Mul(y.(RealNum))
	case RatNum:
		return ratnumMul(x, y.(RatNum))
	case CompNum:
		return y.(CompNum).Mul(x.ToReal())
	}
	return MakeRatnum(IntNum(0), IntNum(0))
}

func (x RatNum) Div(y Number) Number {
	return x.Mul(inverse(y))
}

func gcd(a int64, b int64) int64 {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}

	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	if b < a {
		t := a
		a = b
		b = t
	}
	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}
	return a
}

func Negate(number Number) Number {
	return number.Mul(IntNum(-1))
}

func inverse(a Number) Number {
	switch a.(type) {
	case IntNum:
		return MakeRatnum(IntNum(1), a.(IntNum))
	case RealNum:
		return RealNum(1.0) / a.(RealNum)
	case RatNum:
		ra := a.(RatNum)
		return MakeRatnum(ra.Denominator, ra.Numerator)
	case CompNum:
		return CompNum(1) / a.(CompNum)
	}
	return NanValue
}

func MakeRatnum(numer Number, denom Number) RatNum {
	n := numer.(IntNum)
	d := denom.(IntNum)
	if d == 0 {
		return RatNum{IntNum(0), IntNum(0)} //TODO: return NAN
	}

	negated := false
	if d < 0 {
		n = Negate(n).(IntNum)
		d = Negate(d).(IntNum)
		negated = true
	}
	g := gcd(int64(n), int64(d))
	if g == 1 {
		if negated {
			return RatNum{n, d}
		} else {
			return RatNum{numer, denom}
		}
	} else {
		n = n / IntNum(g)
		d = d / IntNum(g)
		//TODO: when denom is 1, we should return IntNum
		return RatNum{n, d}
	}
}

//CompNum
func (x CompNum) Print(output io.Writer) {
	output.Write([]byte(fmt.Sprintf("%v", x)))
}

func (x CompNum) Add(y Number) Number {
	switch y.(type) {
	case IntNum:
		return x + intToComplex(y.(IntNum))
	case RatNum:
		return x + realToComplex(y.(RatNum).ToReal())
	case RealNum:
		return x + realToComplex(y.(RealNum))
	case CompNum:
		return x + y.(CompNum)
	}
	return NanValue
}

func (x CompNum) Sub(y Number) Number {
	return x.Add(Negate(y))
}

func (x CompNum) Mul(y Number) Number {
	switch y.(type) {
	case IntNum:
		return x * intToComplex(y.(IntNum))
	case RatNum:
		return x * realToComplex(y.(RatNum).ToReal())
	case RealNum:
		return x * realToComplex(y.(RealNum))
	case CompNum:
		return x * y.(CompNum)
	}
	return NanValue
}

func (x CompNum) Div(y Number) Number {
	return x.Mul(inverse(y))
}

/*
//BigNum
func (x BigNum) Print(output io.Writer) {
	if !x.Sign {
		output.Write([]byte("-"))
	}
	for i := len(x.Values) - 1; i >= 0; i-- {
		output.Write([]byte(fmt.Sprintf("%08lx", x.Values[i])))
	}
}

func (x BigNum) Add(args []Number) Number {
}

func (x BigNum) Sub(args []Number) Number {
}

func (x BigNum) Mul(args []Number) Number {
}

func (x BigNum) Mul(args []Number) Number {
}
*/

func StringToNumber(s string) Number {
	//TODO: check radix
	radix := 10
	//TODO: check exact
	exact := true
	return stringToNumberImpl(s, radix, exact)
}

func NumberToString(e Expr) StringExpr {
	var buf bytes.Buffer
	e.Print(&buf)
	return StringExpr(buf.String())
}

func stringToNumberImpl(s string, radix int, exact bool) Number {
	// -> complex (~i or @)
	//case strings.ContainsRune(s, '@'):
	//case strings.HasSuffix(s, "i"):
	if strings.ContainsRune(s, '/') { //ratnum
		nums := strings.Split(s, "/")
		if len(nums) != 2 {
			return NanValue
		}
		n1 := stringToNumberImpl(nums[0], radix, exact).(IntNum)
		n2 := stringToNumberImpl(nums[1], radix, exact).(IntNum)

		r := MakeRatnum(n1, n2)
		if exact {
			return r
		} else {
			return r.ToReal()
		}
	} else if strings.ContainsRune(s, '.') {
		//realnum
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return NanValue
		}
		return RealNum(f)
	} else {
		//integer
		i, err := strconv.ParseInt(s, radix, 64)
		if err != nil {
			return NanValue
		}
		return IntNum(i)
	}
}

func eqInt(a IntNum, b IntNum) bool {
	return a == b
}

func eqReal(a RealNum, b RealNum) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

func eqRatnum(a RatNum, b RatNum) bool {
	return a.Numerator == b.Numerator && a.Denominator == b.Denominator
}

func eqComp(a CompNum, b CompNum) bool {
	ra := real(complex128(a))
	rb := real(complex128(b))
	ia := imag(complex128(a))
	ib := imag(complex128(b))

	if RealNum(ra-rb) < EPSILON && RealNum(rb-ra) < EPSILON &&
		RealNum(ia-ib) < EPSILON && RealNum(ib-ia) < EPSILON {
		return true
	}
	return false
}

func EqNum(a Number, b Number) bool {
	switch a.(type) {
	case IntNum:
		if bi, ok := b.(IntNum); ok {
			return eqInt(a.(IntNum), bi)
		} else if br, ok := b.(RealNum); ok {
			return eqReal(RealNum(a.(IntNum)), br)
		} else if br, ok := b.(RatNum); ok {
			return eqReal(RealNum(a.(IntNum)), br.ToReal())
		} else {
			return eqComp(intToComplex(a.(IntNum)), b.(CompNum))
		}
	case RealNum:
		ar := a.(RealNum)
		if bi, ok := b.(IntNum); ok {
			return eqReal(ar, RealNum(bi))
		} else if br, ok := b.(RealNum); ok {
			return eqReal(ar, br)
		} else if br, ok := b.(RatNum); ok {
			return eqReal(ar, br.ToReal())
		} else {
			return eqComp(realToComplex(ar), b.(CompNum))
		}
	case RatNum:
		ar := a.(RatNum)
		if bi, ok := b.(IntNum); ok {
			return eqRatnum(ar, MakeRatnum(bi, IntNum(1)))
		} else if br, ok := b.(RealNum); ok {
			return eqReal(ar.ToReal(), br)
		} else if br, ok := b.(RatNum); ok {
			return eqRatnum(ar, br)
		} else {
			return eqComp(realToComplex(ar.ToReal()), b.(CompNum))
		}
	case CompNum:
		ac := a.(CompNum)
		if bi, ok := b.(IntNum); ok {
			return eqComp(ac, intToComplex(bi))
		} else if br, ok := b.(RealNum); ok {
			return eqComp(ac, realToComplex(br))
		} else if br, ok := b.(RatNum); ok {
			return eqComp(ac, realToComplex(br.ToReal()))
		} else {
			return eqComp(ac, b.(CompNum))
		}
	}
	return false
}

func GTNum(a Number, b Number) bool {
	switch a.(type) {
	case IntNum:
		ai := a.(IntNum)
		if bi, ok := b.(IntNum); ok {
			return ai > bi
		} else if br, ok := b.(RealNum); ok {
			return RealNum(ai) > br
		} else if br, ok := b.(RatNum); ok {
			return RealNum(ai) > br.ToReal()
		}
	case RealNum:
		ar := a.(RealNum)
		if bi, ok := b.(IntNum); ok {
			return ar > RealNum(bi)
		} else if br, ok := b.(RealNum); ok {
			return ar > br
		} else if br, ok := b.(RatNum); ok {
			return ar > br.ToReal()
		}
	case RatNum:
		ar := a.(RatNum)
		if bi, ok := b.(IntNum); ok {
			return ar.ToReal() > RealNum(bi)
		} else if br, ok := b.(RealNum); ok {
			return ar.ToReal() > br
		} else if br, ok := b.(RatNum); ok {
			return ar.ToReal() > br.ToReal()
		}
		return a.(RatNum).ToReal() > b.(RatNum).ToReal()
	}
	return false
}

func GTENum(a Number, b Number) bool {
	return GTNum(a, b) || EqNum(a, b)
}

func LTNum(a Number, b Number) bool {
	switch a.(type) {
	case IntNum:
		ai := a.(IntNum)
		if bi, ok := b.(IntNum); ok {
			return ai < bi
		} else if br, ok := b.(RealNum); ok {
			return RealNum(ai) < br
		} else if br, ok := b.(RatNum); ok {
			return RealNum(ai) < br.ToReal()
		}
	case RealNum:
		ar := a.(RealNum)
		if bi, ok := b.(IntNum); ok {
			return ar < RealNum(bi)
		} else if br, ok := b.(RealNum); ok {
			return ar < br
		} else if br, ok := b.(RatNum); ok {
			return ar < br.ToReal()
		}
	case RatNum:
		ar := a.(RatNum)
		if bi, ok := b.(IntNum); ok {
			return ar.ToReal() < RealNum(bi)
		} else if br, ok := b.(RealNum); ok {
			return ar.ToReal() < br
		} else if br, ok := b.(RatNum); ok {
			return ar.ToReal() < br.ToReal()
		}
		return a.(RatNum).ToReal() < b.(RatNum).ToReal()
	}
	return false
}

func LTENum(a Number, b Number) bool {
	return LTNum(a, b) || EqNum(a, b)
}
