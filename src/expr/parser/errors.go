package parser

import (
	"expr/lexer"
	"fmt"
)

type TokenError struct {
	Token *lexer.Token
}

func (t TokenError) Error() string {
	return "TokenError: " + t.Token.String()
}

type ParseError struct {
	GotToken        *lexer.Token
	ExpectTokenType lexer.TokenType
	Message         string
}

func (t ParseError) Error() string {
	return fmt.Sprintf("ParseError: expect %s, got %s, %s", t.ExpectTokenType.String(), t.GotToken.String(), t.Message)
}
