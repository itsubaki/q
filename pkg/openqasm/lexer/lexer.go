package lexer

import (
	"bufio"
	"bytes"
	"io"
)

type Lexer struct {
	eof rune
	r   *bufio.Reader
}

func New(r io.Reader) *Lexer {
	return &Lexer{
		eof: rune(-1),
		r:   bufio.NewReader(r),
	}
}

func (l *Lexer) Tokenize() (Token, string) {
	return l.TokenizeIgnore(WHITESPACE)
}

func (l *Lexer) TokenizeIgnore(t ...Token) (Token, string) {
	ignore := make(map[Token]bool)
	for _, tt := range t {
		ignore[tt] = true
	}

	for {
		token, str := l.Scan()
		if _, ok := ignore[token]; ok {
			continue
		}

		return token, str
	}
}

func (l *Lexer) Scan() (Token, string) {
	ch := l.read()
	if ch == l.eof {
		return EOF, ""
	}

	if isWhitespace(ch) {
		l.unread()
		return l.whitespace()
	}

	if isLetter(ch) {
		l.unread()
		str := l.scan()

		if v, ok := keyword[str]; ok {
			return v, str
		}

		return IDENTIFIER, str
	}

	if isDigit(ch) {
		l.unread()
		return l.scanNumber()
	}

	if isString(ch) {
		l.unread()
		return STRING, l.scanString()
	}

	if v, ok := operator[string(ch)]; ok {
		return v, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (l *Lexer) scan() string {
	var buf bytes.Buffer
	_, _ = buf.WriteRune(l.read())

	for {
		ch := l.read()
		if ch == l.eof {
			break
		}

		if isLetter(ch) || isDigit(ch) {
			_, _ = buf.WriteRune(ch)
			continue
		}

		l.unread()
		break
	}

	return buf.String()
}

func (l *Lexer) scanString() string {
	var buf bytes.Buffer
	_, _ = buf.WriteRune(l.read())

	for {
		ch := l.read()
		if ch == l.eof {
			break
		}

		_, _ = buf.WriteRune(ch)

		if isString(ch) {
			break
		}
	}

	return buf.String()
}

func (l *Lexer) scanNumber() (Token, string) {
	var buf bytes.Buffer
	_, _ = buf.WriteRune(l.read())

	token := INT
	for {
		ch := l.read()
		if ch == l.eof {
			break
		}

		if ch == '.' {
			_, _ = buf.WriteRune(ch)
			token = FLOAT
			continue
		}

		if isDigit(ch) {
			_, _ = buf.WriteRune(ch)
			continue
		}

		l.unread()
		break
	}

	return token, buf.String()
}

func (l *Lexer) whitespace() (Token, string) {
	var buf bytes.Buffer
	_, _ = buf.WriteRune(l.read())

	for {
		ch := l.read()
		if ch == l.eof {
			break
		}

		if isWhitespace(ch) {
			buf.WriteRune(ch)
			continue
		}

		l.unread()
		break
	}

	return WHITESPACE, buf.String()
}

func (l *Lexer) read() rune {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		return l.eof
	}

	return ch
}

func (l *Lexer) unread() {
	_ = l.r.UnreadRune()
}

func isWhitespace(ch rune) bool {
	if ch == ' ' {
		return true
	}
	if ch == '\t' {
		return true
	}
	if ch == '\n' {
		return true
	}

	return false
}

func isLetter(ch rune) bool {
	if ch >= 'a' && ch <= 'z' {
		return true
	}
	if ch >= 'A' && ch <= 'Z' {
		return true
	}

	return false
}

func isDigit(ch rune) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}

	return false
}

func isString(ch rune) bool {
	if ch == '"' {
		return true
	}

	if ch == '\'' {
		return true
	}

	return false
}
