package main

type Parser struct {
	tokens  []*Token
	current int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() (Expr, error) {
	// expr, err := p.Expression()
	// if err != nil {
	// 	return nil, err
	// }

	// if !p.IsAtEnd() {
	// 	return nil, ParseError{
	// 		Token:   p.Peek(),
	// 		Message: "Unexpected token after expression.",
	// 	}
	// }

	// return expr, nil

	return p.Expression()
}

func (p *Parser) Expression() (Expr, error) {
	return p.Equality()
}

func (p *Parser) Equality() (Expr, error) {
	expr, err := p.Comparison()

	if err != nil {
		return nil, err
	}

	for p.Match(TokenType_BANG_EQUAL, TokenType_EQUAL_EQUAL) {
		operator := p.Previous()
		right, err := p.Comparison()

		if err != nil {
			return nil, err
		}

		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) Comparison() (Expr, error) {
	expr, err := p.Term()

	if err != nil {
		return nil, err
	}

	for p.Match(TokenType_GREATER, TokenType_GREATER_EQUAL, TokenType_LESS, TokenType_LESS_EQUAL) {
		operator := p.Previous()
		right, err := p.Term()

		if err != nil {
			return nil, err
		}
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) Term() (Expr, error) {
	expr, err := p.Factor()

	if err != nil {
		return nil, err
	}

	for p.Match(TokenType_MINUS, TokenType_PLUS) {
		operator := p.Previous()
		right, err := p.Factor()

		if err != nil {
			return nil, err
		}
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) Factor() (Expr, error) {
	expr, err := p.Unary()

	if err != nil {
		return nil, err
	}

	for p.Match(TokenType_SLASH, TokenType_STAR) {
		operator := p.Previous()
		right, err := p.Unary()

		if err != nil {
			return nil, err
		}

		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) Unary() (Expr, error) {
	if p.Match(TokenType_BANG, TokenType_MINUS) {
		operator := p.Previous()
		right, err := p.Unary()

		if err != nil {
			return nil, err
		}

		return &UnaryExpr{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.Primary()
}

func (p *Parser) Primary() (Expr, error) {
	if p.Match(TokenType_FALSE) {
		return &LiteralExpr{Value: false}, nil
	}
	if p.Match(TokenType_TRUE) {
		return &LiteralExpr{Value: true}, nil
	}
	if p.Match(TokenType_NIL) {
		return &LiteralExpr{Value: nil}, nil
	}

	if p.Match(TokenType_NUMBER, TokenType_STRING) {
		return &LiteralExpr{Value: p.Previous().Literal}, nil
	}

	if p.Match(TokenType_LEFT_PAREN) {
		expr, err := p.Expression()

		if err != nil {
			return nil, err
		}

		p.Consume(TokenType_RIGHT_PAREN, "Expect ')' after expression.")
		return &GroupingExpr{Expression: expr}, nil
	}

	return nil, ParserError{
		Token:   p.Peek(),
		Message: "Expect expression.",
	}
}

func (p *Parser) Consume(tt TokenType, msg string) (*Token, error) {
	if p.Check(tt) {
		next := p.Advance()
		return next, nil
	}

	return nil, ParserError{
		Token:   p.Peek(),
		Message: msg,
	}
}

func (p *Parser) Match(types ...TokenType) bool {
	for _, t := range types {
		if p.Check(t) {
			p.Advance()
			return true
		}
	}
	return false
}

func (p *Parser) Advance() *Token {
	if !p.IsAtEnd() {
		p.current++
	}
	return p.Previous()
}

func (p *Parser) Check(t TokenType) bool {
	if p.IsAtEnd() {
		return false
	}
	return p.Peek().Type == t
}

func (p *Parser) IsAtEnd() bool {
	return p.Peek().Type == TokenType_EOF
}

func (p *Parser) Peek() *Token {
	// Retorna token atual que ainda não foi consumido
	return p.tokens[p.current]
}

func (p *Parser) Previous() *Token {
	// Retorna o token consumido mais recentemente
	return p.tokens[p.current-1]
}

func (p *Parser) Synchronize() {
	p.Advance()
	for !p.IsAtEnd() {
		if p.Previous().Type == TokenType_SEMICOLON {
			return
		}

		switch p.Peek().Type {
		case
			TokenType_CLASS,
			TokenType_FUN,
			TokenType_LET,
			TokenType_IF,
			TokenType_FOR,
			TokenType_WHILE,
			TokenType_PRINT,
			TokenType_RETURN:
			return
		}

		p.Advance()
	}
	return
}
