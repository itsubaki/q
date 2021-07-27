package lexer

type Token int

const (
	// Specials
	ILLEGAL Token = iota
	EOF
	WHITESPACE

	// Literals
	literal_begin
	IDENTIFIER
	STRING
	INT
	FLOAT
	literal_end

	// Operators
	operator_begin
	LBRACKET  // '['
	RBRACKET  // ']'
	LBRACE    // '{'
	RBRACE    // '}'
	LPAREN    // '('
	RPAREN    // ')'
	COLON     // ':'
	SEMICOLON // ';'
	DOT       // '.'
	COMMA     // ','
	EQUALS    // '='
	PLUS      // '+'
	MINUS     // '-'
	MUL       // '*'
	DIV       // '/'
	MOD       // '%'
	ARROW     // "->"
	operator_end

	// Keywords
	keyword_begin
	OPENQASM // OPENQASM
	INCLUDE  // include
	QUBIT    // qubit
	BIT      // bit
	RESET    // reset
	U        // u
	H        // h
	CX       // cx
	MEASURE  // measure
	keyword_end
)

var tokens = [...]string{
	// Specials
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	WHITESPACE: "WHITESPACE",

	// Literals
	IDENTIFIER: "IDENTIFIER",
	STRING:     "STRING",
	INT:        "INT",
	FLOAT:      "FLOAT",

	// Operators
	LBRACKET:  "[",
	RBRACKET:  "]",
	LBRACE:    "{",
	RBRACE:    "}",
	LPAREN:    "(",
	RPAREN:    ")",
	COLON:     ":",
	SEMICOLON: ";",
	DOT:       ".",
	COMMA:     ",",
	EQUALS:    "=",
	PLUS:      "+",
	MINUS:     "-",
	MUL:       "*",
	DIV:       "/",
	MOD:       "%",
	ARROW:     "->",

	// Keywords
	OPENQASM: "OPENQASM",
	INCLUDE:  "include",
	QUBIT:    "qubit",
	BIT:      "bit",
	RESET:    "reset",
	U:        "u",
	H:        "h",
	CX:       "cx",
	MEASURE:  "measure",
}

var (
	operator map[string]Token = make(map[string]Token)
	keyword  map[string]Token = make(map[string]Token)
)

func init() {
	for i := operator_begin + 1; i < operator_end; i++ {
		operator[tokens[i]] = i
	}

	for i := keyword_begin + 1; i < keyword_end; i++ {
		keyword[tokens[i]] = i
	}
}
