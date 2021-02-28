package lexer

import (
	"bufio"
	"bytes"
	"io"
)

type Lexer struct {
	reader *bufio.Reader
	line   int
}

func New(reader io.Reader) *Lexer {
	return &Lexer{reader: bufio.NewReader(reader), line: 1}
}

func (l *Lexer) NextToken() *Token {
	l.skipBlank()
	c := l.peekChar()
	switch c {
	case 0:
		l.readChar()
		return l.newToken(T_EOF, "EOF")
	case ',':
		l.readChar()
		return l.newToken(T_COMMA, "")
	case ';':
		l.readChar()
		return l.newToken(T_SEMICOLON, "")
	case '+':
		l.readChar()
		return l.newToken(T_PLUS, "")
	case '-':
		l.readChar()
		return l.newToken(T_MINUS, "")
	case '*':
		l.readChar()
		return l.newToken(T_ASTERISK, "")
	case '/':
		l.readChar()
		return l.newToken(T_SLASH, "")
	case '(':
		l.readChar()
		return l.newToken(T_LPAREN, "")
	case ')':
		l.readChar()
		return l.newToken(T_RPAREN, "")
	case '[':
		l.readChar()
		return l.newToken(T_LBRACKET, "")
	case ']':
		l.readChar()
		return l.newToken(T_RBRACKET, "")
	case '{':
		l.readChar()
		return l.newToken(T_LBRACE, "")
	case '}':
		l.readChar()
		return l.newToken(T_RBRACE, "")
	case '"':
		return l.readString()
	case '.':
		l.readChar()
		if isNumber(l.peekChar()) {
			return l.readFloat("0")
		} else {
			return l.newToken(T_DOT, "")
		}
	case ':':
		l.readChar()
		if l.peekChar() == '=' {
			l.readChar()
			return l.newToken(T_DECLARATION, "")
		} else {
			return l.newToken(T_COLON, "")
		}
	case '!':
		l.readChar()
		if l.peekChar() == '=' {
			l.readChar()
			return l.newToken(T_NEQ, "")
		} else {
			return l.newToken(T_BANG, "")
		}
	case '<':
		l.readChar()
		if l.peekChar() == '=' {
			l.readChar()
			return l.newToken(T_LE, "")
		} else {
			return l.newToken(T_LT, "")
		}
	case '>':
		l.readChar()
		if l.peekChar() == '=' {
			l.readChar()
			return l.newToken(T_GE, "")
		} else {
			return l.newToken(T_GT, "")
		}
	case '=':
		l.readChar()
		if l.peekChar() == '=' {
			l.readChar()
			return l.newToken(T_EQ, "")
		} else {
			return l.newToken(T_ASSIGN, "")
		}
	default:
		if isNumber(c) {
			return l.readNumber()
		} else if isLetter(c) {
			return l.readWord()
		}
		l.readChar()
		return l.newToken(T_ILLEGAL, "")
	}
}

func (l *Lexer) newToken(t TokenType, message string) *Token {
	return &Token{Type: t, Line: l.line, Message: message}
}

func (l *Lexer) skipBlank() {
	for c := l.peekChar(); c != 0 && c <= 32; c = l.peekChar() {
		if l.readChar() == '\n' {
			l.line++
		}
	}
}

func (l *Lexer) peekChar() byte {
	b, err := l.reader.Peek(1)
	if err == io.EOF {
		return 0
	}
	return b[0]
}

func (l *Lexer) readChar() byte {
	b, err := l.reader.ReadByte()
	if err == io.EOF {
		return 0
	}
	return b
}

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
}

func (l *Lexer) readNumber() *Token {
	buf := bytes.Buffer{}
	for {
		c := l.peekChar()
		if c == '.' {
			l.readChar()
			return l.readFloat(buf.String())
		} else if isNumber(c) {
			buf.WriteByte(l.readChar())
		} else {
			break
		}
	}
	s := buf.String()
	if s == "" {
		s = "0"
	}
	return l.newToken(T_INT, s)
}

func (l *Lexer) readFloat(intPart string) *Token {
	decimal := l.readNumber()
	if decimal.Type == T_INT {
		return l.newToken(T_FLOAT, intPart+"."+decimal.Message)
	} else {
		return l.newToken(T_ILLEGAL, intPart+"."+decimal.Message)
	}
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (l *Lexer) readWord() *Token {
	buf := bytes.Buffer{}
	for c := l.peekChar(); isLetter(c) || isNumber(c); c = l.peekChar() {
		buf.WriteByte(l.readChar())
	}
	s := buf.String()
	t := keywordOrIdent(s)
	message := ""
	if t == T_IDENT {
		message = s
	}
	return l.newToken(t, message)
}

func (l *Lexer) readString() *Token {
	l.readChar()
	buf := bytes.Buffer{}
	for l.peekChar() != '"' {
		buf.WriteByte(l.readChar())
	}
	l.readChar()
	return l.newToken(T_STRING, buf.String())
}
