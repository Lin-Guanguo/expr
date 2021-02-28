package lexer

import "fmt"

type TokenType int

type Token struct {
	Type    TokenType
	Line    int
	Message string
}

func (t *Token) String() string {
	return fmt.Sprintf("%s line: %d; %s;", t.Type.String(), t.Line, t.Message)
}

func (t *Token) TypeIs(tt TokenType) bool {
	return t.Type == tt
}

const (
	T_NONE TokenType = iota
	T_ILLEGAL
	T_EOF

	T_IDENT
	T_INT
	T_FLOAT
	T_STRING
	T_TABLE

	T_COMMA
	T_SEMICOLON
	T_DOT
	T_COLON

	T_ASSIGN
	T_PLUS
	T_MINUS
	T_ASTERISK
	T_SLASH
	T_BANG

	T_LT
	T_LE
	T_GT
	T_GE
	T_EQ
	T_NEQ

	T_LPAREN
	T_RPAREN
	T_LBRACKET
	T_RBRACKET
	T_LBRACE
	T_RBRACE

	T_DECLARATION

	T_FUNCTION
	T_TRUE
	T_FALSE
	T_NIL
	T_IF
	T_ELSE
	T_FOR
	T_BREAK
	T_CONTINUE
	T_RETURN
	T_AND
	T_OR
)

func (t TokenType) String() string {
	switch t {
	case T_NONE:
		return "T_NONE"
	case T_ILLEGAL:
		return "T_ILLEGAL"
	case T_EOF:
		return "T_EOF"
	case T_IDENT:
		return "T_IDENT"
	case T_INT:
		return "T_INT"
	case T_FLOAT:
		return "T_FLOAT"
	case T_STRING:
		return "T_STRING"
	case T_TABLE:
		return "T_TABLE"
	case T_COMMA:
		return "T_COMMA"
	case T_SEMICOLON:
		return "T_SEMICOLON"
	case T_DOT:
		return "T_DOT"
	case T_COLON:
		return "T_COLON"
	case T_ASSIGN:
		return "T_ASSIGN"
	case T_PLUS:
		return "T_PLUS"
	case T_MINUS:
		return "T_MINUS"
	case T_ASTERISK:
		return "T_ASTERISK"
	case T_SLASH:
		return "T_SLASH"
	case T_BANG:
		return "T_BANG"
	case T_LT:
		return "T_LT"
	case T_LE:
		return "T_LE"
	case T_GT:
		return "T_GT"
	case T_GE:
		return "T_GE"
	case T_EQ:
		return "T_EQ"
	case T_NEQ:
		return "T_NEQ"
	case T_LPAREN:
		return "T_LPAREN"
	case T_RPAREN:
		return "T_RPAREN"
	case T_LBRACKET:
		return "T_LBRACKET"
	case T_RBRACKET:
		return "T_RBRACKET"
	case T_LBRACE:
		return "T_LBRACE"
	case T_RBRACE:
		return "T_RBRACE"
	case T_DECLARATION:
		return "T_DECLARATION"
	case T_FUNCTION:
		return "T_FUNCTION"
	case T_TRUE:
		return "T_TRUE"
	case T_FALSE:
		return "T_FALSE"
	case T_NIL:
		return "T_NIL"
	case T_IF:
		return "T_IF"
	case T_ELSE:
		return "T_ELSE"
	case T_FOR:
		return "T_FOR"
	case T_BREAK:
		return "T_BREAK"
	case T_CONTINUE:
		return "T_CONTINUE"
	case T_RETURN:
		return "T_RETURN"
	case T_AND:
		return "T_AND"
	case T_OR:
		return "T_OR"
	default:
		panic("unknown token type to string")
	}
}

func keywordOrIdent(word string) TokenType {
	switch word {
	case "func":
		return T_FUNCTION
	case "true":
		return T_TRUE
	case "false":
		return T_FALSE
	case "if":
		return T_IF
	case "else":
		return T_ELSE
	case "for":
		return T_FOR
	case "return":
		return T_RETURN
	case "break":
		return T_BREAK
	case "continue":
		return T_CONTINUE
	case "nil":
		return T_NIL
	case "and":
		return T_AND
	case "or":
		return T_OR
	case "table":
		return T_TABLE
	default:
		return T_IDENT
	}
}
