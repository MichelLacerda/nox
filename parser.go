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

// func (p *Parser) Parse() (Expr, error) {
// 	return p.Expression()
// }

func (p *Parser) Parse() ([]Stmt, error) {
	statements := make([]Stmt, 0)

	for !p.IsAtEnd() {
		// stmt, err := p.Statement()
		// if err != nil {
		// 	return nil, err
		// }
		// statements = append(statements, stmt)
		d, err := p.declaration()

		if err != nil {
			return nil, err
		}

		statements = append(statements, d)
	}

	return statements, nil
}

func (p *Parser) declaration() (Stmt, error) {
	if p.Match(TokenType_LET) {
		stmt, err := p.VarDeclaration()
		if err != nil {
			p.Synchronize()
			return nil, err
		}
		return stmt, nil
	}

	stmt, err := p.Statement()
	if err != nil {
		p.Synchronize()
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) VarDeclaration() (Stmt, error) {
	name, err := p.Consume(TokenType_IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.Match(TokenType_EQUAL) {
		initializer, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.Consume(TokenType_SEMICOLON, "Expect ';' after variable declaration."); err != nil {
		return nil, err
	}

	return &VarStmt{name, initializer}, nil
}

func (p *Parser) Statement() (Stmt, error) {
	if p.Match(TokenType_FOR) {
		return p.ForStatement()
	}

	if p.Match(TokenType_IF) {
		return p.IfStatement()
	}

	if p.Match(TokenType_PRINT) {
		return p.PrintStatement()
	}

	if p.Match(TokenType_WHILE) {
		return p.WhileStatement()
	}

	if p.Match(TokenType_LEFT_BRACE) {
		return &BlockStmt{
			Statements: p.Block(),
		}, nil
	}

	return p.ExpressionStatement()
}

func (p *Parser) ForStatement() (Stmt, error) {
	if _, err := p.Consume(TokenType_LEFT_PAREN, "Expect '(' after 'for'."); err != nil {
		return nil, err
	}

	var initializer Stmt
	if p.Match(TokenType_SEMICOLON) {
		initializer = nil
	} else if p.Match(TokenType_LET) {
		varDecl, err := p.VarDeclaration()
		if err != nil {
			return nil, err
		}
		initializer = varDecl
	} else {
		stmt, err := p.ExpressionStatement()
		if err != nil {
			return nil, err
		}
		initializer = stmt
	}

	var condition Expr
	if !p.Check(TokenType_SEMICOLON) {
		var err error
		condExpr, err := p.Expression()
		if err != nil {
			return nil, err
		}
		condition = condExpr
	}

	if _, err := p.Consume(TokenType_SEMICOLON, "Expect ';' after loop condition."); err != nil {
		return nil, err
	}

	var increment Expr
	if !p.Check(TokenType_RIGHT_PAREN) {
		exprInc, err := p.Expression()
		if err != nil {
			return nil, err
		}
		increment = exprInc
	}
	p.Consume(TokenType_RIGHT_PAREN, "Expect ')' after for clauses.")

	body, err := p.Statement()
	if err != nil {
		return nil, err
	}

	// Sempre cria um novo bloco para o incremento, mesmo se body já for bloco
	if increment != nil {
		var stmts []Stmt
		if block, ok := body.(*BlockStmt); ok {
			stmts = append([]Stmt{}, block.Statements...)
			stmts = append(stmts, &ExpressionStmt{increment})
		} else {
			stmts = []Stmt{body, &ExpressionStmt{increment}}
		}
		body = &BlockStmt{Statements: stmts}
	}

	if condition == nil {
		condition = &LiteralExpr{Value: true}
	}

	body = &WhileStmt{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &BlockStmt{
			Statements: []Stmt{initializer, body},
		}
	}

	return body, nil
}

func (p *Parser) WhileStatement() (Stmt, error) {
	if _, err := p.Consume(TokenType_LEFT_PAREN, "Expect '(' after 'while'."); err != nil {
		return nil, err
	}

	condition, err := p.Expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.Consume(TokenType_RIGHT_PAREN, "Expect ')' after while condition."); err != nil {
		return nil, err
	}

	body, err := p.Statement()
	if err != nil {
		return nil, err
	}

	return &WhileStmt{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) IfStatement() (Stmt, error) {
	if _, err := p.Consume(TokenType_LEFT_PAREN, "Expect '(' after 'if'."); err != nil {
		return nil, err
	}

	condition, err := p.Expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.Consume(TokenType_RIGHT_PAREN, "Expect ')' after if condition."); err != nil {
		return nil, err
	}

	thenStmt, err := p.Statement()
	if err != nil {
		return nil, err
	}

	var elseStmt Stmt
	if p.Match(TokenType_ELSE) {
		elseStmt, err = p.Statement()
		if err != nil {
			return nil, err
		}
	}

	return &IfStmt{
		Condition: condition,
		Then:      thenStmt,
		Else:      elseStmt,
	}, nil
}

func (p *Parser) Block() []Stmt {
	statements := make([]Stmt, 0)

	for !p.IsAtEnd() && !p.Check(TokenType_RIGHT_BRACE) {
		stmt, err := p.declaration()
		if err != nil {
			p.Synchronize()
			continue
		}
		statements = append(statements, stmt)
	}

	_, err := p.Consume(TokenType_RIGHT_BRACE, "Expect '}' after block.")

	if err != nil {
		return nil
	}

	return statements
}

func (p *Parser) PrintStatement() (Stmt, error) {
	expr, err := p.Expression()
	if err != nil {
		return nil, err
	}

	_, err = p.Consume(TokenType_SEMICOLON, "Expect ';' after value.")

	if err != nil {
		return nil, err
	}

	return &PrintStmt{
		Expression: expr,
	}, nil
}

func (p *Parser) ExpressionStatement() (Stmt, error) {
	expr, err := p.Expression()
	if err != nil {
		return nil, err
	}

	p.Consume(TokenType_SEMICOLON, "Expect ';' after expression.")

	return &ExpressionStmt{
		Expression: expr,
	}, nil
}

func (p *Parser) Expression() (Expr, error) {
	return p.Assignment()
}

func (p *Parser) Assignment() (Expr, error) {
	expr, err := p.Or()
	if err != nil {
		return nil, err
	}

	if p.Match(TokenType_EQUAL) {
		equals := p.Previous()
		value, err := p.Assignment()
		if err != nil {
			return nil, err
		}

		if varExpr, ok := expr.(*VariableExpr); ok {
			return &AssignExpr{
				Name:  varExpr.Name,
				Value: value,
			}, nil
		}

		return nil, ParserError{
			Token:   equals,
			Message: "Invalid assignment target.",
		}
	}

	return expr, nil
}

func (p *Parser) Or() (Expr, error) {
	expr, err := p.And()
	if err != nil {
		return nil, err
	}

	for p.Match(TokenType_OR) {
		operator := p.Previous()
		right, err := p.And()
		if err != nil {
			return nil, err
		}
		expr = &LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) And() (Expr, error) {
	expr, err := p.Equality()
	if err != nil {
		return nil, err
	}

	for p.Match(TokenType_AND) {
		operator := p.Previous()
		right, err := p.Equality()
		if err != nil {
			return nil, err
		}
		expr = &LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
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

	if p.Match(TokenType_IDENTIFIER) {
		t := p.Previous()
		return &VariableExpr{t}, nil
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
