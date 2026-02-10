package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// identifer
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// operator
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	// Equal
	EQ = "=="
	// Not Equal
	NEQ = "!="

	// Less Than
	LT = "<"
	// Greater Than
	GT = ">"

	// Less Than Equal
	LTE = "<="
	// Greater Than Equal
	GTE = ">="

	// delimiter
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	LBRACKET = "["
	RBRACKET = "]"

	COLON = ":"

	// keyword
	FUNC   = "FUNCTION"
	LET    = "LET"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
	FOR    = "FOR"
	BREAK  = "BREAK"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func MakeToken(t TokenType, char byte) Token {
	return Token{
		Type:    t,
		Literal: string(char),
	}
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keyword[ident]; ok {
		return tok
	}

	return IDENT
}

var keyword = map[string]TokenType{
	"ft":     FUNC,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"for":    FOR,
	"break":  BREAK,
}
