package parser

import (
	"fmt"
	"strconv"

	"github.com/potix2/goscheme/scm"
)

const (
	EOF = -1
)

type Error struct {
	Message  string
	Pos      Position
	Filename string
	Fatal    bool
}

func (e *Error) Error() string {
	return e.Message
}

type Token struct {
	Tok int
	Lit string
	Pos Position
}

func (t *Token) SetPosition(pos Position) {
	t.Pos = pos
}

type Position struct {
	Line   int
	Column int
}

type Pos interface {
	Position() Position
	SetPosition(Position)
}

type PosImpl struct {
	pos Position
}

func (x *PosImpl) Position() Position {
	return x.pos
}

func (x *PosImpl) SetPosition(pos Position) {
	x.pos = pos
}

type Scanner struct {
	src      []rune
	offset   int
	lineHead int
	line     int
}

func (s *Scanner) Init(src string) {
	s.src = []rune(src)
	fmt.Printf("INIT: %s\n", src)
}

func (s *Scanner) Scan() (tok int, lit string, pos Position, err error) {
	s.skipBrank()
	pos = s.pos()
	switch ch1, ch2 := s.peek2(); {
	case isBoolean(ch1, ch2):
		lit, err = s.scanBoolean()
		if err != nil {
			return
		}
		tok = BOOLEAN
	case ch1 == '#':
		lit, err = s.scanNumber()
		if err != nil {
			return
		}
		tok = NUMBER
	case ch1 == '.':
		if isPeculiarIdent(ch1, ch2) {
			lit, err = s.scanIdent()
			if err != nil {
				return
			}
			tok = IDENT
		} else if isCharDelimiter(ch2) {
			tok = int(ch1)
			lit = string(ch1)
			s.next()
		} else {
			lit, err = s.scanNumber()
			if err != nil {
				return
			}
			tok = NUMBER
		}
	case ch1 == '"':
		lit, err = s.scanString()
		if err != nil {
			return
		}
		tok = STRING
	case isSign(ch1):
		if isPeculiarIdent(ch1, ch2) {
			lit, err = s.scanIdent()
			if err != nil {
				return
			}
			tok = IDENT
		} else if isCharDelimiter(ch2) {
			lit, err = s.scanIdent()
			if err != nil {
				return
			}
			tok = IDENT
		} else {
			lit, err = s.scanNumber()
			if err != nil {
				return
			}
			tok = NUMBER
		}
	case isDigit(ch1):
		lit, err = s.scanNumber()
		if err != nil {
			return
		}
		tok = NUMBER
	case isInitial(ch1):
		lit, err = s.scanIdent()
		if err != nil {
			return
		}
		tok = IDENT
	default:
		switch ch1 {
		case -1:
			tok = EOF
		case '(', ')', '\'', 'e', '+', '-', '/', '.':
			tok = int(ch1)
			lit = string(ch1)
		}
		s.next()
	}
	return
}

func isBrank(ch rune) bool             { return ch == ' ' || ch == '\t' || ch == '\n' }
func isLetter(ch rune) bool            { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }
func isDigit(ch rune) bool             { return (ch >= '0' && ch <= '9') }
func isSign(ch rune) bool              { return ch == '+' || ch == '-' }
func isSignSubsequent(ch rune) bool    { return isInitial(ch) || isSign(ch) || ch == '@' }
func isSubsequent(ch rune) bool        { return isInitial(ch) || isDigit(ch) || isSpecialSubsequent(ch) }
func isSpecialSubsequent(ch rune) bool { return isSign(ch) || ch == '.' || ch == '@' }
func isDotSubsequent(ch rune) bool     { return isSignSubsequent(ch) || ch == '.' }
func isIdent(ch rune) bool             { return isInitial(ch) || isDigit(ch) }
func isInitial(ch rune) bool           { return isLetter(ch) || isSpecialSuffix(ch) }

func isCharDelimiter(ch rune) bool {
	return ch == '(' || ch == ')' || ch == '[' || ch == ']' || ch == '"' || ch == ';' || ch == '#' || isBrank(ch)
}

func isSpecialSuffix(ch rune) bool {
	return ch == '!' ||
		ch == '$' ||
		ch == '%' ||
		ch == '&' ||
		ch == '*' ||
		ch == '/' ||
		ch == ':' ||
		ch == '<' ||
		ch == '=' ||
		ch == '>' ||
		ch == '?' ||
		ch == '^' ||
		ch == '_' ||
		ch == '~'
}

func isPeculiarIdent(ch1 rune, ch2 rune) bool {
	return (isSign(ch1) && isBrank(ch2)) ||
		(isSign(ch1) && isSignSubsequent(ch2)) ||
		(isSign(ch1) && ch2 == '.') ||
		(ch1 == '.' && isDotSubsequent(ch2))
}

func isNumber(ch1 rune, ch2 rune, ch3 rune) bool {
	if ch1 == '#' {
		return ch2 == 'i' || ch2 == 'e' || ch2 == 'b' || ch2 == 'o' || ch2 == 'd' || ch2 == 'x'
	} else if isSign(ch1) {
		return isDigit(ch2) || ch2 == 'i' || ch2 == 'n' || (ch2 == '.' && isDigit(ch3))
	} else {
		return isDigit(ch1) || ch1 == '.'
	}
}

func isBoolean(ch1 rune, ch2 rune) bool { return ch1 == '#' && (ch2 == 't' || ch2 == 'f') }
func isLineComment(ch rune) bool        { return ch == ';' }

//func isSpecialChar(ch rune) bool { return ch == '#' }

func (s *Scanner) peek() rune {
	if !s.reachEOF() {
		return s.src[s.offset]
	} else {
		return -1
	}
}

func (s *Scanner) peek2() (rune, rune) {
	var ch1, ch2 rune
	ch1 = -1
	ch2 = -1
	ch1 = s.peek()
	if ch1 == -1 {
		return ch1, ch2
	}

	s.next()

	ch2 = s.peek()
	s.back()
	return ch1, ch2
}

func (s *Scanner) reachEOF() bool {
	return len(s.src) <= s.offset
}

func (s *Scanner) pos() Position {
	return Position{Line: s.line + 1, Column: s.offset - s.lineHead + 1}
}

func (s *Scanner) skipBrank() {
	for isBrank(s.peek()) {
		s.next()
	}

	if isLineComment(s.peek()) {
		s.skipLineComment()
	}
}

func (s *Scanner) skipLineComment() {
	for s.peek() != '\n' {
		s.next()
	}
	s.next()
}

func (s *Scanner) back() {
	s.offset--
}

func (s *Scanner) next() {
	if !s.reachEOF() {
		if s.peek() == '\n' {
			s.lineHead = s.offset + 1
			s.line++
		}
		s.offset++
	}
}

func (s *Scanner) scanIdent() (lit string, err error) {
	var ret []rune
	if isInitial(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}

	for isSubsequent(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}
	return string(ret), nil
}

func (s *Scanner) scanPeculiarIdent() (lit string, err error) {
	var ret []rune
	if isSign(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
		if isSignSubsequent(s.peek()) {
			ret = append(ret, s.peek())
			s.next()
		} else if s.peek() == '.' {
			ret = append(ret, s.peek())
			s.next()
			if isDotSubsequent(s.peek()) {
				ret = append(ret, s.peek())
				s.next()
			} else {
				return "", &Error{Message: "syntax error", Pos: s.pos(), Fatal: false}
			}
		} else {
			return string(ret), nil
		}
	} else if s.peek() == ',' {
		ret = append(ret, s.peek())
		s.next()
		if isDotSubsequent(s.peek()) {
			ret = append(ret, s.peek())
			s.next()
		} else {
			return "", &Error{Message: "syntax error", Pos: s.pos(), Fatal: false}
		}
	}

	for isSubsequent(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}
	return string(ret), nil
}

func (s *Scanner) scanNumber() (lit string, err error) {
	var ret []rune
	for !s.reachEOF() && !isCharDelimiter(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}

	return string(ret), nil
}

var escapedChars = map[rune]rune{
	'"':  '"',
	'\\': '\\',
	'|':  '|',
	'a':  '\a',
	'b':  '\t',
	't':  '\t',
	'n':  '\n',
	'r':  '\r',
}

func (s *Scanner) scanString() (lit string, err error) {
	//assert(s.peek() == '"')
	s.next()

	var ret []rune
	for !s.reachEOF() && s.peek() != '"' {
		if s.peek() == '\\' {
			s.next()
			switch s.peek() {
			case '"', '\\', '|', 'a', 'b', 't', 'n', 'r':
				ret = append(ret, escapedChars[s.peek()])
			case 'x':
				s.next()
				var hexDigital []rune
				for s.peek() != ';' {
					hexDigital = append(hexDigital, s.peek())
					s.next()
				}
				ui, err := strconv.ParseUint(string(hexDigital), 16, 64)
				if err != nil {
					return "", err
				}
				ret = append(ret, rune(ui))
			}
		} else {
			ret = append(ret, s.peek())
		}
		s.next()
	}
	s.next()

	return string(ret), nil
}

func (s *Scanner) scanBoolean() (lit string, err error) {
	if s.peek() != '#' {
		return "", &Error{Message: "invalid state of lexer", Pos: s.pos(), Fatal: false}
	}
	s.next()

	var buf []rune
	for isLetter(s.peek()) {
		buf = append(buf, s.peek())
		s.next()
	}

	ret := string(buf)
	switch ret {
	case "t", "true":
		return "true", nil
	case "f", "false":
		return "false", nil
	default:
		return "", &Error{Message: "syntax error", Pos: s.pos(), Fatal: false}
	}
}

type Lexer struct {
	s     *Scanner
	lit   string
	pos   Position
	e     error
	expr  scm.Expr
	exprs []scm.Expr
}

func (l *Lexer) Lex(lval *yySymType) int {
	tok, lit, pos, err := l.s.Scan()
	if err != nil {
		l.e = &Error{Message: err.Error(), Pos: pos, Fatal: true}
	}
	lval.tok = Token{Tok: tok, Lit: lit}
	lval.tok.SetPosition(pos)
	l.lit = lit
	l.pos = pos
	return tok
}

func (l *Lexer) Error(e string) {
	l.e = &Error{Message: e, Pos: l.pos, Fatal: false}
}

func Parse(s *Scanner) ([]scm.Expr, error) {
	l := Lexer{s: s, exprs: []scm.Expr{}}
	if yyParse(&l) != 0 {
		return nil, l.e
	}
	return l.exprs, l.e
}

func Read(src string) ([]scm.Expr, error) {
	scanner := &Scanner{
		src: []rune(src),
	}
	return Parse(scanner)
}
