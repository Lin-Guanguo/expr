package parser

import (
	"expr/lexer"
	"strconv"
)

type prefixParseFn func() (Expression, *ParseError)
type infixParseFn func(leftExpr Expression) (Expression, *ParseError)

type Parser struct {
	lexer     *lexer.Lexer
	peekToken *lexer.Token

	Errors []error

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:          l,
		peekToken:      &lexer.Token{Type: lexer.T_LBRACE, Line: 0, Message: "begin of file"},
		Errors:         nil,
		prefixParseFns: make(map[lexer.TokenType]prefixParseFn),
		infixParseFns:  make(map[lexer.TokenType]infixParseFn),
	}
	p.registerParseFn()
	return p
}

func (p *Parser) nextToken() *lexer.Token {
	token := p.peekToken
	p.peekToken = p.lexer.NextToken()
	if p.peekToken.Type == lexer.T_ILLEGAL {
		p.Errors = append(p.Errors, &TokenError{Token: p.nextToken()})
	}
	return token
}

func (p *Parser) checkPeekToken(tt lexer.TokenType) *ParseError {
	if p.peekToken.Type != tt {
		return &ParseError{
			GotToken:        p.peekToken,
			ExpectTokenType: tt,
			Message:         "",
		}
	}
	return nil
}

func (p *Parser) ParseProgram() *BlockExpr {
	block := &BlockExpr{
		Token: p.peekToken,
		Exprs: nil,
	}
	p.nextToken()
	for p.peekToken.Type != lexer.T_EOF {
		expr, err := p.parseEntireExpr()
		if err != nil {
			p.Errors = append(p.Errors, err)
		} else {
			block.Exprs = append(block.Exprs, expr)
		}
		for p.peekToken.Type == lexer.T_SEMICOLON {
			p.nextToken()
		}
	}
	p.nextToken()
	return block
}

func (p *Parser) parseBlockExpr() (Expression, *ParseError) {
	token := p.nextToken()
	block := &BlockExpr{
		Token: token,
		Exprs: nil,
	}
	for p.peekToken.Type != lexer.T_RBRACE {
		expr, err := p.parseEntireExpr()
		if err != nil {
			p.Errors = append(p.Errors, err)
		} else {
			block.Exprs = append(block.Exprs, expr)
		}
		for p.peekToken.Type == lexer.T_SEMICOLON {
			p.nextToken()
		}
	}
	if token := p.nextToken(); token.Type == lexer.T_EOF {
		return nil, &ParseError{
			GotToken:        token,
			ExpectTokenType: lexer.T_RBRACE,
			Message:         "{}block brace mismatch",
		}
	}
	return block, nil
}

type Precedence int

const (
	LOWEST  Precedence = iota
	ASSIGN             // = :=
	OR                 // or
	AND                // and
	EQUALS             // == >= <=
	COMPARE            // > <
	SUM                // +
	PRODUCT            // *
	PREFIX             // -X or !X
	CALL               // myFunction(X)
)

var tokenPrecedences = map[lexer.TokenType]Precedence{
	lexer.T_SEMICOLON:   LOWEST,
	lexer.T_RPAREN:      LOWEST,
	lexer.T_DECLARATION: ASSIGN,
	lexer.T_ASSIGN:      ASSIGN,
	lexer.T_OR:          OR,
	lexer.T_AND:         AND,
	lexer.T_EQ:          EQUALS,
	lexer.T_NEQ:         EQUALS,
	lexer.T_GE:          EQUALS,
	lexer.T_LE:          EQUALS,
	lexer.T_GT:          COMPARE,
	lexer.T_LT:          COMPARE,
	lexer.T_PLUS:        SUM,
	lexer.T_MINUS:       SUM,
	lexer.T_ASTERISK:    PRODUCT,
	lexer.T_SLASH:       PRODUCT,
	lexer.T_LPAREN:      CALL,
	lexer.T_DOT:         CALL,
	lexer.T_LBRACKET:    CALL,
	lexer.T_COLON:       CALL,
}

func (p *Parser) parseEntireExpr() (Expression, *ParseError) {
	expr, err := p.parseExpr(LOWEST)
	return expr, err
}

func (p *Parser) parseExpr(precedence Precedence) (Expression, *ParseError) {
	prefixFn := p.prefixParseFns[p.peekToken.Type]
	if prefixFn == nil {
		token := p.nextToken()
		return nil, &ParseError{
			GotToken:        token,
			ExpectTokenType: 0,
			Message:         "no prefix parse function",
		}
	}

	leftExpr, err := prefixFn()
	if err != nil {
		return nil, err
	}

	for precedence < tokenPrecedences[p.peekToken.Type] {
		infixFn := p.infixParseFns[p.peekToken.Type]
		if infixFn == nil {
			return leftExpr, nil
		}
		leftExpr, err = infixFn(leftExpr)
		if err != nil {
			return nil, err
		}
	}
	return leftExpr, nil
}

func (p *Parser) registerParseFn() {
	p.prefixParseFns[lexer.T_INT] = p.parseInteger
	p.prefixParseFns[lexer.T_FLOAT] = p.parseFloat
	p.prefixParseFns[lexer.T_TRUE] = p.parseBoolean
	p.prefixParseFns[lexer.T_FALSE] = p.parseBoolean
	p.prefixParseFns[lexer.T_STRING] = p.parseString
	p.prefixParseFns[lexer.T_TABLE] = p.parseTableExpr
	p.prefixParseFns[lexer.T_LBRACKET] = p.parsePackExpr
	p.prefixParseFns[lexer.T_NIL] = p.parseNil
	p.prefixParseFns[lexer.T_IDENT] = p.parserIdentifier
	p.prefixParseFns[lexer.T_RETURN] = p.parserReturnExpr
	p.prefixParseFns[lexer.T_BREAK] = p.parserBreakExpr
	p.prefixParseFns[lexer.T_LPAREN] = p.parseGroupedExpr
	p.prefixParseFns[lexer.T_LBRACE] = p.parseBlockExpr
	p.prefixParseFns[lexer.T_IF] = p.parseIfExpr
	p.prefixParseFns[lexer.T_FUNCTION] = p.parseFuncExpr
	p.prefixParseFns[lexer.T_FOR] = p.parseForExpr
	p.prefixParseFns[lexer.T_MINUS] = p.parseArithPrefixExpr
	p.prefixParseFns[lexer.T_BANG] = p.parseArithPrefixExpr

	p.infixParseFns[lexer.T_DECLARATION] = p.parserDeclarationExpr
	p.infixParseFns[lexer.T_ASSIGN] = p.parserAssignExpr
	p.infixParseFns[lexer.T_LPAREN] = p.parseCallExpr
	p.infixParseFns[lexer.T_DOT] = p.parseDotExpr
	p.infixParseFns[lexer.T_PLUS] = p.parseArithInfixExpr
	p.infixParseFns[lexer.T_MINUS] = p.parseArithInfixExpr
	p.infixParseFns[lexer.T_ASTERISK] = p.parseArithInfixExpr
	p.infixParseFns[lexer.T_SLASH] = p.parseArithInfixExpr
	p.infixParseFns[lexer.T_EQ] = p.parseArithInfixExpr
	p.infixParseFns[lexer.T_NEQ] = p.parseArithInfixExpr
	p.infixParseFns[lexer.T_LE] = p.parseArithInfixExpr
	p.infixParseFns[lexer.T_GE] = p.parseArithInfixExpr
	p.infixParseFns[lexer.T_LT] = p.parseArithInfixExpr
	p.infixParseFns[lexer.T_GT] = p.parseArithInfixExpr
	p.infixParseFns[lexer.T_AND] = p.parseArithInfixExpr
	p.infixParseFns[lexer.T_OR] = p.parseArithInfixExpr
}

func (p *Parser) parseInteger() (Expression, *ParseError) {
	token := p.nextToken()
	val, err := strconv.ParseInt(token.Message, 0, 64)
	if err != nil {
		return nil, &ParseError{
			GotToken:        token,
			ExpectTokenType: lexer.T_INT,
			Message:         "parseInt error",
		}
	}
	return &IntegerExpr{
		Token: token,
		Value: val,
	}, nil
}

func (p *Parser) parseFloat() (Expression, *ParseError) {
	token := p.nextToken()
	val, err := strconv.ParseFloat(token.Message, 64)
	if err != nil {
		return nil, &ParseError{
			GotToken:        token,
			ExpectTokenType: lexer.T_FLOAT,
			Message:         "parseFloat error",
		}
	}
	return &FloatExpr{
		Token: token,
		Value: val,
	}, nil
}

func (p *Parser) parseBoolean() (Expression, *ParseError) {
	token := p.nextToken()
	if token.Type == lexer.T_TRUE {
		return &BooleanExpr{
			Token: token,
			Value: true,
		}, nil
	} else if token.Type == lexer.T_FALSE {
		return &BooleanExpr{
			Token: token,
			Value: false,
		}, nil
	} else {
		return nil, &ParseError{
			GotToken:        token,
			ExpectTokenType: lexer.T_TRUE,
			Message:         "parseBoolean token error,",
		}
	}
}

func (p *Parser) parseString() (Expression, *ParseError) {
	token := p.nextToken()
	if token.Type == lexer.T_STRING {
		return &StringExpr{
			Token: token,
			Value: token.Message,
		}, nil
	} else {
		return nil, &ParseError{
			GotToken:        token,
			ExpectTokenType: lexer.T_TRUE,
			Message:         "parseString token error,",
		}
	}
}

func (p *Parser) parseNil() (Expression, *ParseError) {
	token := p.nextToken()
	if token.Type == lexer.T_NIL {
		return &NilExpr{
			Token: token,
		}, nil
	} else {
		return nil, &ParseError{
			GotToken:        token,
			ExpectTokenType: lexer.T_TRUE,
			Message:         "parseNil token error,",
		}
	}
}

func (p *Parser) parserIdentifier() (Expression, *ParseError) {
	token := p.nextToken()
	return &Identifier{
		Token: token,
		Ident: token.Message,
	}, nil
}

func (p *Parser) parserReturnExpr() (Expression, *ParseError) {
	token := p.nextToken()
	val, err := p.parseEntireExpr()
	if err != nil {
		return nil, err
	}
	return &ReturnExpr{
		Token:       token,
		ReturnValue: val,
	}, nil
}

func (p *Parser) parserBreakExpr() (Expression, *ParseError) {
	token := p.nextToken()
	val, err := p.parseEntireExpr()
	if err != nil {
		return nil, err
	}
	return &BreakExpr{
		Token:      token,
		BreakValue: val,
	}, nil
}

func (p *Parser) parserDeclarationExpr(leftExpr Expression) (Expression, *ParseError) {
	token := p.nextToken()
	left, ok := leftExpr.(AssignableExpr)
	if !ok || !left.IsAssignable() {
		return nil, &ParseError{
			GotToken:        token,
			ExpectTokenType: 0,
			Message:         "Declaration’s left expression can't assign",
		}
	}
	right, err := p.parseExpr(ASSIGN - 1) // 降低一级优先级，形成右结合
	if err != nil {
		return nil, err
	}
	return &DeclarationExpr{
		Token: token,
		Left:  left,
		Value: right,
	}, nil
}

//值使用了 p.parseEntireExpr() 连续赋值会右结合
func (p *Parser) parserAssignExpr(leftExpr Expression) (Expression, *ParseError) {
	token := p.nextToken()
	left, ok := leftExpr.(AssignableExpr)
	if !ok || !left.IsAssignable() {
		return nil, &ParseError{
			GotToken:        token,
			ExpectTokenType: 0,
			Message:         "Assign’s left expression can't assign",
		}
	}
	right, err := p.parseExpr(ASSIGN - 1) // 降低一级优先级，形成右结合
	if err != nil {
		return nil, err
	}
	return &AssignExpr{
		Token: token,
		Left:  left,
		Value: right,
	}, nil
}

func (p *Parser) parseArithInfixExpr(leftExpr Expression) (Expression, *ParseError) {
	token := p.nextToken()
	switch token.Type {
	case lexer.T_PLUS, lexer.T_MINUS, lexer.T_ASTERISK, lexer.T_SLASH,
		lexer.T_EQ, lexer.T_NEQ, lexer.T_LE, lexer.T_GE, lexer.T_LT, lexer.T_GT,
		lexer.T_AND, lexer.T_OR:
		break
	default:
		return nil, &ParseError{
			GotToken:        token,
			ExpectTokenType: 0,
			Message:         "unknown infix arith operator",
		}
	}

	prePrecedence := tokenPrecedences[token.Type]
	right, err := p.parseExpr(prePrecedence)
	if err != nil {
		return nil, err
	}
	return &ArithInfixExpr{
		Token: token,
		Op:    token.Type,
		Left:  leftExpr,
		Right: right,
	}, nil
}

func (p *Parser) parseArithPrefixExpr() (Expression, *ParseError) {
	token := p.nextToken()
	switch token.Type {
	case lexer.T_MINUS, lexer.T_BANG:
		break
	default:
		return nil, &ParseError{
			GotToken:        token,
			ExpectTokenType: 0,
			Message:         "unknown prefix arith operator",
		}
	}

	right, err := p.parseExpr(PREFIX)
	if err != nil {
		return nil, err
	}

	return &ArithPrefixExpr{
		Token: token,
		Op:    token.Type,
		Right: right,
	}, nil
}

func (p *Parser) parseGroupedExpr() (Expression, *ParseError) {
	_ = p.nextToken()
	expr, err := p.parseExpr(LOWEST)
	if err != nil {
		return nil, err
	}
	if err := p.checkPeekToken(lexer.T_RPAREN); err != nil {
		err.Message = "()Parentheses mismatch"
		return nil, err
	}
	_ = p.nextToken()
	return expr, nil
}

func (p *Parser) parseCommaExprs(begin lexer.TokenType, end lexer.TokenType) ([]Expression, *ParseError) {
	if err := p.checkPeekToken(begin); err != nil {
		return nil, err
	}
	p.nextToken()
	exprs := make([]Expression, 0)
	for p.peekToken.Type != end {
		expr, err := p.parseEntireExpr()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
		if p.peekToken.Type == lexer.T_COMMA {
			p.nextToken()
		} else {
			break
		}
	}
	if err := p.checkPeekToken(end); err != nil {
		return nil, err
	}
	p.nextToken()
	return exprs, nil
}

func (p *Parser) parseImplicitBlockExpr() (*BlockExpr, *ParseError) {
	if p.peekToken.Type == lexer.T_LBRACE {
		expr, err := p.parseBlockExpr()
		if err != nil {
			return nil, err
		}
		return expr.(*BlockExpr), nil
	} else {
		token := p.peekToken
		expr, err := p.parseEntireExpr()
		if err != nil {
			return nil, err
		}
		token.Message = "Implicit BlockExpr"
		if p.peekToken.Type == lexer.T_SEMICOLON {
			p.nextToken()
		}
		return &BlockExpr{
			Token: token,
			Exprs: []Expression{expr},
		}, nil
	}
}

func (p *Parser) parseIfExpr() (Expression, *ParseError) {
	token := p.nextToken()

	condition, err := p.parseEntireExpr()
	if err != nil {
		return nil, err
	}

	consequence, err := p.parseImplicitBlockExpr()
	if err != nil {
		return nil, err
	}
	if p.peekToken.Type == lexer.T_ELSE {
		p.nextToken()
		alternative, err := p.parseImplicitBlockExpr()
		if err != nil {
			return nil, err
		}
		return &IfExpr{
			Token:       token,
			Condition:   condition,
			Consequence: consequence,
			Alternative: alternative,
		}, nil
	} else {
		return &IfExpr{
			Token:       token,
			Condition:   condition,
			Consequence: consequence,
			Alternative: nil,
		}, nil
	}
}

func (p *Parser) parseFuncExpr() (Expression, *ParseError) {
	token := p.nextToken()
	funcExpr := &FuncExpr{
		Token:      token,
		Parameters: nil,
		Capture:    nil,
		Body:       nil,
	}

	//read parameters
	{
		parameters, err := p.parseCommaExprs(lexer.T_LPAREN, lexer.T_RPAREN)
		if err != nil {
			return nil, err
		}
		for _, param := range parameters {
			if Ident, ok := param.(*Identifier); ok {
				funcExpr.Parameters = append(funcExpr.Parameters, Ident)
			} else {
				return nil, &ParseError{
					GotToken:        nil,
					ExpectTokenType: 0,
					Message:         "function parameter type error, could only be identifier",
				}
			}
		}
	}

	//read capture
	if p.peekToken.Type == lexer.T_LBRACKET {
		captures, err := p.parseCommaExprs(lexer.T_LBRACKET, lexer.T_RBRACKET)
		if err != nil {
			return nil, err
		}
		for _, param := range captures {
			if c, ok := param.(FuncCaptureExpr); ok {
				funcExpr.Capture = append(funcExpr.Capture, c)
			} else {
				return nil, &ParseError{
					GotToken:        nil,
					ExpectTokenType: 0,
					Message:         "function capture type error, could only be identifier or declaration",
				}
			}
		}
	}

	//read body
	body, err := p.parseImplicitBlockExpr()
	if err != nil {
		return nil, err
	}
	funcExpr.Body = body

	return funcExpr, nil
}

func (p *Parser) parseIndexExpr() (Expression, *ParseError) {
	if p.peekToken.Type == lexer.T_IDENT {
		stringIndexToken := p.nextToken()
		return &StringExpr{
			Token: stringIndexToken,
			Value: stringIndexToken.Message,
		}, nil
	} else {
		if err := p.checkPeekToken(lexer.T_LBRACKET); err != nil {
			return nil, err
		} else {
			p.nextToken()
		}
		index, err := p.parseEntireExpr()
		if err != nil {
			return nil, err
		}
		if err := p.checkPeekToken(lexer.T_RBRACKET); err != nil {
			return nil, err
		} else {
			p.nextToken()
		}
		return index, nil
	}
}

func (p *Parser) parseTableExpr() (Expression, *ParseError) {
	token := p.nextToken()
	table := &TableExpr{
		Token:     token,
		InitValue: nil,
	}
	if err := p.checkPeekToken(lexer.T_LBRACE); err != nil {
		return nil, err
	}
	p.nextToken()
	for p.peekToken.Type != lexer.T_RBRACE {
		key, err := p.parseIndexExpr()
		if err != nil {
			return nil, err
		}
		if err := p.checkPeekToken(lexer.T_ASSIGN); err != nil {
			return nil, err
		} else {
			p.nextToken()
		}
		value, err := p.parseEntireExpr()
		if err != nil {
			return nil, err
		}
		table.InitValue = append(table.InitValue, KeyValuePair{key, value})
		if p.peekToken.Type == lexer.T_COMMA {
			p.nextToken()
		} else {
			break
		}
	}
	if err := p.checkPeekToken(lexer.T_RBRACE); err != nil {
		return nil, err
	}
	p.nextToken()
	return table, nil
}

func (p *Parser) parseDotExpr(leftExpr Expression) (Expression, *ParseError) {
	token := p.nextToken()
	index, err := p.parseIndexExpr()
	if err != nil {
		return nil, err
	}
	return &IndexExpr{
		Token: token,
		Table: leftExpr,
		Index: index,
	}, nil
}
func (p *Parser) parseCallExpr(leftExpr Expression) (Expression, *ParseError) {
	token := p.peekToken
	parameters, err := p.parseCommaExprs(lexer.T_LPAREN, lexer.T_RPAREN)
	if err != nil {
		return nil, err
	}
	return &CallExpr{
		Token:      token,
		Function:   leftExpr,
		Parameters: parameters,
	}, nil
}

func (p *Parser) parseForExpr() (Expression, *ParseError) {
	token := p.nextToken()
	initExpr, err := p.parseEntireExpr()
	if err != nil {
		return nil, err
	}
	if err := p.checkPeekToken(lexer.T_SEMICOLON); err != nil {
		return nil, err
	} else {
		p.nextToken()
	}

	edgeExpr, err := p.parseEntireExpr()
	if err != nil {
		return nil, err
	}
	if err := p.checkPeekToken(lexer.T_SEMICOLON); err != nil {
		return nil, err
	} else {
		p.nextToken()
	}

	stepExpe, err := p.parseEntireExpr()
	if err != nil {
		return nil, err
	}
	if p.peekToken.Type == lexer.T_SEMICOLON {
		p.nextToken()
	}

	body, err := p.parseImplicitBlockExpr()
	if err != nil {
		return nil, err
	}
	return &ForExpr{
		Token:    token,
		InitExpr: initExpr,
		EdgeExpr: edgeExpr,
		StepExpr: stepExpe,
		Body:     body,
	}, nil
}

func (p *Parser) parsePackExpr() (Expression, *ParseError) {
	pack := &PackExpr{
		Token: p.peekToken,
		Exprs: nil,
	}
	exprs, err := p.parseCommaExprs(lexer.T_LBRACKET, lexer.T_RBRACKET)
	if err != nil {
		return nil, err
	}
	for _, expr := range exprs {
		pack.Exprs = append(pack.Exprs, expr)
	}
	return pack, nil
}
