package lexer_test

import (
	"os"
	"strings"
	"testing"

	"github.com/itsubaki/q/pkg/openqasm/lexer"
)

func TestLexer(t *testing.T) {
	type item struct {
		token lexer.Token
		str   string
	}

	var cases = []struct {
		in   string
		want []item
	}{
		{
			in: "../_testdata/bell.qasm",
			want: []item{
				{lexer.OPENQASM, "OPENQASM"},
				{lexer.FLOAT, "3.0"},
				{lexer.SEMICOLON, ";"},

				{lexer.INCLUDE, "include"},
				{lexer.STRING, "\"stdgates.qasm\""},
				{lexer.SEMICOLON, ";"},

				{lexer.QUBIT, "qubit"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "2"},
				{lexer.RBRACKET, "]"},
				{lexer.IDENTIFIER, "q"},
				{lexer.SEMICOLON, ";"},

				{lexer.BIT, "bit"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "2"},
				{lexer.RBRACKET, "]"},
				{lexer.IDENTIFIER, "c"},
				{lexer.SEMICOLON, ";"},

				{lexer.RESET, "reset"},
				{lexer.IDENTIFIER, "q"},
				{lexer.SEMICOLON, ";"},

				{lexer.H, "h"},
				{lexer.IDENTIFIER, "q"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "0"},
				{lexer.RBRACKET, "]"},
				{lexer.SEMICOLON, ";"},

				{lexer.CX, "cx"},
				{lexer.IDENTIFIER, "q"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "0"},
				{lexer.RBRACKET, "]"},
				{lexer.COMMA, ","},
				{lexer.IDENTIFIER, "q"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "1"},
				{lexer.RBRACKET, "]"},
				{lexer.SEMICOLON, ";"},

				{lexer.IDENTIFIER, "c"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "0"},
				{lexer.RBRACKET, "]"},
				{lexer.EQUALS, "="},
				{lexer.MEASURE, "measure"},
				{lexer.IDENTIFIER, "q"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "0"},
				{lexer.RBRACKET, "]"},
				{lexer.SEMICOLON, ";"},

				{lexer.IDENTIFIER, "c"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "1"},
				{lexer.RBRACKET, "]"},
				{lexer.EQUALS, "="},
				{lexer.MEASURE, "measure"},
				{lexer.IDENTIFIER, "q"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "1"},
				{lexer.RBRACKET, "]"},
				{lexer.SEMICOLON, ";"},
				{lexer.EOF, ""},
			},
		},
	}

	for _, c := range cases {
		f, err := os.ReadFile(c.in)
		if err != nil {
			t.Fatalf("read file: %v", err)
		}

		lex := lexer.New(strings.NewReader(string(f)))
		for _, w := range c.want {
			token, str := lex.Tokenize()
			if token != w.token || str != w.str {
				t.Errorf("got=%v:%v, want=%v:%v", token, str, w.token, w.str)
			}
		}
	}
}

func TestLexerTokenize(t *testing.T) {
	var cases = []struct {
		in   string
		want []lexer.Token
	}{
		{"1", []lexer.Token{lexer.INT}},
		{"-1", []lexer.Token{lexer.MINUS, lexer.INT}},
		{"100", []lexer.Token{lexer.INT, lexer.EOF}},
		{"10.0", []lexer.Token{lexer.FLOAT, lexer.EOF}},
		{"\"abc\"", []lexer.Token{lexer.STRING, lexer.EOF}},
		{"'abc'", []lexer.Token{lexer.STRING, lexer.EOF}},
		{"abc", []lexer.Token{lexer.IDENTIFIER, lexer.EOF}},
		{" \t\n", []lexer.Token{lexer.WHITESPACE}},
		{"\\", []lexer.Token{lexer.ILLEGAL}},
		{"\"a", []lexer.Token{lexer.STRING, lexer.EOF}},
	}

	for _, c := range cases {
		lex := lexer.New(strings.NewReader(c.in))
		for _, w := range c.want {
			if got, _ := lex.TokenizeIgnore(); got != w {
				t.Fail()
			}
		}
	}
}
