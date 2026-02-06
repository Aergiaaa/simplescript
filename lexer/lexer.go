package lexer

import "github.com/Aergiaaa/idiotic_interpreter/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
}

func InitLexer(s string) *Lexer {
	l := &Lexer{input: s}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.skipWhiteSpace()

	makeToken := token.MakeToken
	switch l.char {
	case '"':
		t.Type = token.STRING
		t.Literal = l.readString()
	case '=':
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			t = token.Token{
				Type:    token.EQ,
				Literal: string(char) + string(l.char),
			}
		} else {
			t = makeToken(token.ASSIGN, l.char)
		}
	case '!':
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			t = token.Token{
				Type:    token.NEQ,
				Literal: string(char) + string(l.char),
			}
		} else {
			t = makeToken(token.BANG, l.char)
		}
	case '<':
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			t = token.Token{
				Type:    token.LTE,
				Literal: string(char) + string(l.char),
			}
		} else {
			t = makeToken(token.LT, l.char)
		}
	case '>':
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			t = token.Token{
				Type:    token.GTE,
				Literal: string(char) + string(l.char),
			}
		} else {
			t = makeToken(token.GT, l.char)
		}
	case '+':
		t = makeToken(token.PLUS, l.char)
	case '-':
		t = makeToken(token.MINUS, l.char)
	case '*':
		t = makeToken(token.ASTERISK, l.char)
	case '/':
		t = makeToken(token.SLASH, l.char)
	case '(':
		t = makeToken(token.LPAREN, l.char)
	case ')':
		t = makeToken(token.RPAREN, l.char)
	case '{':
		t = makeToken(token.LBRACE, l.char)
	case '}':
		t = makeToken(token.RBRACE, l.char)
	case '[':
		t = makeToken(token.LBRACKET, l.char)
	case ']':
		t = makeToken(token.RBRACKET, l.char)
	case ':':
		t = makeToken(token.COLON, l.char)
	case ',':
		t = makeToken(token.COMMA, l.char)
	case ';':
		t = makeToken(token.SEMICOLON, l.char)
	case 0:
		t.Literal = ""
		t.Type = token.EOF
	default:
		if isLetter(l.char) {
			t.Literal = l.readIdentifier()
			t.Type = token.LookupIdent(t.Literal)
			return t
		}

		if isDigit(l.char) {
			t.Type = token.INT
			t.Literal = l.readNum()
			return t
		}
		t = makeToken(token.ILLEGAL, l.char)
	}

	l.readChar()
	return t
}

func (l *Lexer) readString() string {
	pos := l.position + 1
	for {
		l.readChar()

		if l.char == '"' || l.char == 0 {
			break
		}
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readNum() string {
	pos := l.position
	for isDigit(l.char) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.char) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) skipWhiteSpace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func isDigit(char byte) bool {
	return '0' <= char && '9' >= char
}

func isLetter(char byte) bool {
	return 'a' <= char && 'z' >= char || 'A' <= char && 'Z' >= char || char == '_'
}
