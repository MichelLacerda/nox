package parser

import (
	"fmt"

	"github.com/MichelLacerda/nox/internal/ast"
	"github.com/MichelLacerda/nox/internal/token"
)

type Parser struct {
	tokens  []*token.Token
	current int
}

type ParserError struct {
	Token   *token.Token
	Message string
}

func (e ParserError) Error() string {
	return fmt.Sprintf("Parser Error at %s: %s", e.Token.Lexeme, e.Message)
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ([]ast.Stmt, error) {
	statements := []ast.Stmt{}

	for !p.IsAtEnd() {
		d, err := p.declaration()

		if err != nil {
			return nil, err
		}

		statements = append(statements, d)
	}

	return statements, nil
}

func (p *Parser) declaration() (ast.Stmt, error) {
	if p.Match(token.TokenType_EXPORT) {
		decl, err := p.exportDeclaration()
		if err != nil {
			return nil, err
		}
		return &ast.ExportStmt{Declaration: decl}, nil
	}

	if p.Match(token.TokenType_IMPORT) {
		stmt, err := p.ImportStmt()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	}

	if p.Match(token.TokenType_CLASS) {
		stmt, err := p.ClassDeclaration()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	}

	if p.Match(token.TokenType_FUNC) {
		stmt, err := p.Function("function")

		if err != nil {
			return nil, err // Não sincroniza, apenas retorna o erro
		}

		return stmt, nil
	}

	if p.Match(token.TokenType_LET) {
		stmt, err := p.VarDeclaration()
		if err != nil {
			return nil, err // Não sincroniza, apenas retorna o erro
		}
		return stmt, nil
	}

	stmt, err := p.Statement()
	if err != nil {
		return nil, err // Não sincroniza, apenas retorna o erro
	}
	return stmt, nil
}

func (p *Parser) exportDeclaration() (ast.Stmt, error) {
	if p.Match(token.TokenType_FUNC) {
		return p.Function("function")
	}
	if p.Match(token.TokenType_CLASS) {
		return p.ClassDeclaration()
	}
	if p.Match(token.TokenType_LET) {
		return p.VarDeclaration()
	}

	return nil, ParserError{
		Token:   p.Peek(),
		Message: "Expect declaration after 'export'.",
	}
}

func (p *Parser) ImportStmt() (ast.Stmt, error) {
	pathToken, err := p.Consume(token.TokenType_STRING, "Expect module path.")
	if err != nil {
		return nil, err
	}

	var aliasToken *token.Token
	if p.Match(token.TokenType_AS) {
		tok := p.Previous()

		if !p.Match(token.TokenType_IDENTIFIER) {
			return nil, ParserError{
				Token:   tok,
				Message: "Expect identifier after 'as'",
			}
		}
		aliasToken = p.Previous()
	}

	// if _, err := p.Consume(token.TokenType_SEMICOLON, "Expect ';' after use statement."); err != nil {
	// 	return nil, err
	// }

	// Optional semicolon
	p.Match(token.TokenType_SEMICOLON)

	return &ast.ImportStmt{Path: pathToken, Alias: aliasToken}, nil
}

func (p *Parser) ClassDeclaration() (ast.Stmt, error) {
	name, err := p.Consume(token.TokenType_IDENTIFIER, "Expect class name.")
	if err != nil {
		return nil, err
	}

	var superclass ast.Expr
	if p.Match(token.TokenType_LESS) {
		tok, err := p.Consume(token.TokenType_IDENTIFIER, "Expect superclass name after '<'.")
		if err != nil {
			return nil, err
		}

		expr := ast.Expr(&ast.VariableExpr{Name: tok})
		for p.Match(token.TokenType_DOT) {
			name, _ := p.Consume(token.TokenType_IDENTIFIER, "Expect property name after '.'.")
			expr = &ast.GetExpr{
				Object: expr,
				Name:   name,
			}
		}
		superclass = expr
	}

	if _, err := p.Consume(token.TokenType_LEFT_BRACE, "Expect '{' before class body."); err != nil {
		return nil, err
	}

	methods := []*ast.FunctionStmt{}
	for !p.IsAtEnd() && !p.Check(token.TokenType_RIGHT_BRACE) {
		funcStmt, err := p.Function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, funcStmt)
	}

	if _, err := p.Consume(token.TokenType_RIGHT_BRACE, "Expect '}' after class body."); err != nil {
		return nil, err
	}

	return ast.NewClassStmt(name, superclass, methods), nil
}

func (p *Parser) Function(kind string) (*ast.FunctionStmt, error) {
	name, err := p.Consume(token.TokenType_IDENTIFIER, "Expect "+kind+" name.")
	if err != nil {
		return nil, err
	}

	if _, err := p.Consume(token.TokenType_LEFT_PAREN, "Expect '(' after "+kind+" name."); err != nil {
		return nil, err
	}

	parameters := []*token.Token{}
	for !p.Check(token.TokenType_RIGHT_PAREN) {
		if len(parameters) >= 255 {
			return nil, ParserError{
				Token:   p.Peek(),
				Message: "Cannot have more than 255 parameters.",
			}
		}
		param, err := p.Consume(token.TokenType_IDENTIFIER, "Expect parameter name.")
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, param)
		if !p.Match(token.TokenType_COMMA) {
			break
		}
	}

	if _, err := p.Consume(token.TokenType_RIGHT_PAREN, "Expect ')' after parameters."); err != nil {
		return nil, err
	}

	// Consome o '{' antes do corpo da função
	if _, err := p.Consume(token.TokenType_LEFT_BRACE, "Expect '{' before "+kind+" body."); err != nil {
		return nil, err
	}

	body, err := p.Block()
	if err != nil {
		return nil, err
	}

	return &ast.FunctionStmt{
		Name:       name,
		Parameters: parameters,
		Body:       body,
	}, nil
}

func (p *Parser) VarDeclaration() (ast.Stmt, error) {
	name, err := p.Consume(token.TokenType_IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr
	if p.Match(token.TokenType_EQUAL) {
		initializer, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}

	// Permite ;, fim de bloco ou EOF como final de declaração
	// if !(p.Match(token.TokenType_SEMICOLON) || p.Check(token.TokenType_RIGHT_BRACE) || p.Check(token.TokenType_EOF)) {
	// 	// Usa o token da variável para reportar a linha correta
	// 	return nil, ParserError{
	// 		Token:   name,
	// 		Message: "Expect ';' after variable declaration.",
	// 	}
	// }

	// Optional semicolon
	p.Match(token.TokenType_SEMICOLON)

	return &ast.VarStmt{Name: name, Initializer: initializer}, nil
}

func (p *Parser) Statement() (ast.Stmt, error) {
	if p.Match(token.TokenType_FOR) {
		// return p.ForStatement()
		return p.ForInStatement()
	}

	if p.Match(token.TokenType_BREAK) {
		keyword := p.Previous()
		// p.Consume(token.TokenType_SEMICOLON, "Expect ';' after 'break'.")
		// Optional semicolon
		p.Match(token.TokenType_SEMICOLON)
		return &ast.BreakStmt{Keyword: keyword}, nil
	}

	if p.Match(token.TokenType_CONTINUE) {
		keyword := p.Previous()
		// p.Consume(token.TokenType_SEMICOLON, "Expect ';' after 'continue'.")
		p.Match(token.TokenType_SEMICOLON)
		p.Match(token.TokenType_SEMICOLON)
		return &ast.ContinueStmt{Keyword: keyword}, nil
	}

	if p.Match(token.TokenType_WITH) {
		return p.WithStatement()
	}

	if p.Match(token.TokenType_IF) {
		return p.IfStatement()
	}

	if p.Match(token.TokenType_PRINT) {
		return p.PrintStatement()
	}

	if p.Match(token.TokenType_RETURN) {
		return p.ReturnStatement()
	}

	// if p.Match(token.TokenType_WHILE) {
	// 	return p.WhileStatement()
	// }

	if p.Match(token.TokenType_LEFT_BRACE) {
		stmts, err := p.Block()
		if err != nil {
			return nil, err
		}
		return &ast.BlockStmt{
			Statements: stmts,
		}, nil
	}

	return p.ExpressionStatement()
}

func (p *Parser) WithStatement() (ast.Stmt, error) {
	resourceExpr, err := p.Expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.Consume(token.TokenType_AS, "Expect 'as' after with resource."); err != nil {
		return nil, err
	}

	alias, err := p.Consume(token.TokenType_IDENTIFIER, "Expect variable name after 'as'.")
	if err != nil {
		return nil, err
	}

	body, err := p.Statement()
	if err != nil {
		return nil, err
	}

	return &ast.WithStmt{
		Resource: resourceExpr,
		Alias:    alias,
		Body:     body,
	}, nil
}

func (p *Parser) ReturnStatement() (ast.Stmt, error) {
	keyword := p.Previous()

	var value ast.Expr
	if !p.Check(token.TokenType_SEMICOLON) && !p.Check(token.TokenType_RIGHT_BRACE) && !p.IsAtEnd() {
		var err error
		value, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}

	// Optional semicolon — can be omitted before } or end of file
	p.Match(token.TokenType_SEMICOLON)

	return &ast.ReturnStmt{
		Keyword: keyword,
		Value:   value,
	}, nil
}

func (p *Parser) ForInStatement() (ast.Stmt, error) {
	// Suporta o estilo Go: for { ... }
	if p.Match(token.TokenType_LEFT_BRACE) {
		bodyStmts, err := p.Block()
		if err != nil {
			return nil, err
		}
		body := &ast.BlockStmt{Statements: bodyStmts}

		return &ast.ForInStmt{
			IndexVar: nil,
			ValueVar: nil,
			Iterable: &ast.LiteralExpr{Value: true}, // sinaliza loop infinito
			Body:     body,
		}, nil
	}

	// Parse do estilo: for index, value in iterable { ... }
	// index pode ser "_"
	var indexVar *token.Token = nil
	var valueVar *token.Token

	// Primeiro identificador
	firstIdent, err := p.Consume(token.TokenType_IDENTIFIER, "Expect loop variable.")
	if err != nil {
		return nil, err
	}

	if p.Match(token.TokenType_COMMA) {
		indexVar = firstIdent
		secondIdent, err := p.Consume(token.TokenType_IDENTIFIER, "Expect value variable after comma.")
		if err != nil {
			return nil, err
		}
		valueVar = secondIdent
	} else {
		valueVar = firstIdent
	}

	if _, err := p.Consume(token.TokenType_IN, "Expect 'in' after loop variables."); err != nil {
		return nil, err
	}

	iterable, err := p.Expression()
	if err != nil {
		return nil, err
	}

	body, err := p.Statement()
	if err != nil {
		return nil, err
	}

	return &ast.ForInStmt{
		IndexVar: indexVar,
		ValueVar: valueVar,
		Iterable: iterable,
		Body:     body,
	}, nil
}

// func (p *Parser) ForStatement() (ast.Stmt, error) {
// 	if _, err := p.Consume(token.TokenType_LEFT_PAREN, "Expect '(' after 'for'."); err != nil {
// 		return nil, err
// 	}

// 	var initializer ast.Stmt
// 	if p.Match(token.TokenType_SEMICOLON) {
// 		initializer = nil
// 	} else if p.Match(token.TokenType_LET) {
// 		varDecl, err := p.VarDeclaration()
// 		if err != nil {
// 			return nil, err
// 		}
// 		initializer = varDecl
// 	} else {
// 		stmt, err := p.ExpressionStatement()
// 		if err != nil {
// 			return nil, err
// 		}
// 		initializer = stmt
// 	}

// 	var condition ast.Expr
// 	if !p.Check(token.TokenType_SEMICOLON) {
// 		var err error
// 		condExpr, err := p.Expression()
// 		if err != nil {
// 			return nil, err
// 		}
// 		condition = condExpr
// 	}

// 	// if _, err := p.Consume(token.TokenType_SEMICOLON, "Expect ';' after loop condition."); err != nil {
// 	// 	return nil, err
// 	// }

// 	var increment ast.Expr
// 	if !p.Check(token.TokenType_RIGHT_PAREN) {
// 		exprInc, err := p.Expression()
// 		if err != nil {
// 			return nil, err
// 		}
// 		increment = exprInc
// 	}
// 	p.Consume(token.TokenType_RIGHT_PAREN, "Expect ')' after for clauses.")

// 	body, err := p.Statement()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Sempre cria um novo bloco para o incremento, mesmo se body já for bloco
// 	if increment != nil {
// 		var stmts []ast.Stmt
// 		if block, ok := body.(*ast.BlockStmt); ok {
// 			stmts = append([]ast.Stmt{}, block.Statements...)
// 			stmts = append(stmts, &ast.ExpressionStmt{Expression: increment})
// 		} else {
// 			stmts = []ast.Stmt{body, &ast.ExpressionStmt{Expression: increment}}
// 		}
// 		body = &ast.BlockStmt{Statements: stmts}
// 	}

// 	if condition == nil {
// 		condition = &ast.LiteralExpr{Value: true}
// 	}

// 	body = &ast.WhileStmt{
// 		Condition: condition,
// 		Body:      body,
// 	}

// 	if initializer != nil {
// 		body = &ast.BlockStmt{
// 			Statements: []ast.Stmt{initializer, body},
// 		}
// 	}

// 	return body, nil
// }

// func (p *Parser) WhileStatement() (ast.Stmt, error) {
// 	if _, err := p.Consume(token.TokenType_LEFT_PAREN, "Expect '(' after 'while'."); err != nil {
// 		return nil, err
// 	}

// 	condition, err := p.Expression()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if _, err := p.Consume(token.TokenType_RIGHT_PAREN, "Expect ')' after while condition."); err != nil {
// 		return nil, err
// 	}

// 	body, err := p.Statement()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &ast.WhileStmt{
// 		Condition: condition,
// 		Body:      body,
// 	}, nil
// }

func (p *Parser) IfStatement() (ast.Stmt, error) {
	// if _, err := p.Consume(token.TokenType_LEFT_PAREN, "Expect '(' after 'if'."); err != nil {
	// 	return nil, err
	// }
	p.Match(token.TokenType_LEFT_PAREN)

	condition, err := p.Expression()
	if err != nil {
		return nil, err
	}

	p.Match(token.TokenType_RIGHT_PAREN)

	// if _, err := p.Consume(token.TokenType_RIGHT_PAREN, "Expect ')' after if condition."); err != nil {
	// 	return nil, err
	// }

	thenStmt, err := p.Statement()
	if err != nil {
		return nil, err
	}

	var elseStmt ast.Stmt
	if p.Match(token.TokenType_ELSE) {
		elseStmt, err = p.Statement()
		if err != nil {
			return nil, err
		}
	}

	return &ast.IfStmt{
		Condition: condition,
		Then:      thenStmt,
		Else:      elseStmt,
	}, nil
}

func (p *Parser) Block() ([]ast.Stmt, error) {
	statements := []ast.Stmt{}

	for !p.IsAtEnd() && !p.Check(token.TokenType_RIGHT_BRACE) {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	_, err := p.Consume(token.TokenType_RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) PrintStatement() (ast.Stmt, error) {
	var expressions []ast.Expr

	for {
		expr, err := p.Expression()
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)

		if !p.Match(token.TokenType_COMMA) {
			break
		}
	}

	// if _, err := p.Consume(token.TokenType_SEMICOLON, "Expect ';' after value."); err != nil {
	// 	return nil, err
	// }

	// Optional semicolon
	p.Match(token.TokenType_SEMICOLON)

	return &ast.PrintStmt{Expressions: expressions}, nil
}

func (p *Parser) ExpressionStatement() (ast.Stmt, error) {
	expr, err := p.Expression()
	if err != nil {
		return nil, err
	}

	// p.Consume(token.TokenType_SEMICOLON, "Expect ';' after expression.")
	// Optional semicolon
	p.Match(token.TokenType_SEMICOLON)

	return &ast.ExpressionStmt{
		Expression: expr,
	}, nil
}

func (p *Parser) Expression() (ast.Expr, error) {
	return p.Assignment()
}

func (p *Parser) Assignment() (ast.Expr, error) {
	expr, err := p.Or()
	if err != nil {
		return nil, err
	}

	if p.Match(token.TokenType_EQUAL) {
		equals := p.Previous()
		value, err := p.Assignment()
		if err != nil {
			return nil, err
		}

		if varExpr, ok := expr.(*ast.VariableExpr); ok {
			return &ast.AssignExpr{
				Name:  varExpr.Name,
				Value: value,
			}, nil
		} else if getExpr, ok := expr.(*ast.GetExpr); ok {
			return &ast.SetExpr{
				Object: getExpr.Object,
				Name:   getExpr.Name,
				Value:  value,
			}, nil
		}

		return nil, ParserError{
			Token:   equals,
			Message: "Invalid assignment target.",
		}
	}

	return expr, nil
}

func (p *Parser) Or() (ast.Expr, error) {
	expr, err := p.And()
	if err != nil {
		return nil, err
	}

	for p.Match(token.TokenType_OR) {
		operator := p.Previous()
		right, err := p.And()
		if err != nil {
			return nil, err
		}
		expr = &ast.LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) And() (ast.Expr, error) {
	expr, err := p.Equality()
	if err != nil {
		return nil, err
	}

	for p.Match(token.TokenType_AND) {
		operator := p.Previous()
		right, err := p.Equality()
		if err != nil {
			return nil, err
		}
		expr = &ast.LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) Equality() (ast.Expr, error) {
	expr, err := p.Comparison()

	if err != nil {
		return nil, err
	}

	for p.Match(token.TokenType_BANG_EQUAL, token.TokenType_EQUAL_EQUAL) {
		operator := p.Previous()
		right, err := p.Comparison()

		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) Comparison() (ast.Expr, error) {
	expr, err := p.Term()

	if err != nil {
		return nil, err
	}

	for p.Match(token.TokenType_GREATER, token.TokenType_GREATER_EQUAL, token.TokenType_LESS, token.TokenType_LESS_EQUAL) {
		operator := p.Previous()
		right, err := p.Term()

		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) Term() (ast.Expr, error) {
	expr, err := p.Factor()

	if err != nil {
		return nil, err
	}

	for p.Match(token.TokenType_MINUS, token.TokenType_PLUS) {
		operator := p.Previous()
		right, err := p.Factor()

		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) Factor() (ast.Expr, error) {
	expr, err := p.Unary()

	if err != nil {
		return nil, err
	}

	for p.Match(token.TokenType_SLASH, token.TokenType_STAR, token.TokenType_PERCENT, token.TokenType_DOUBLE_STAR) {
		operator := p.Previous()
		right, err := p.Unary()

		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) Unary() (ast.Expr, error) {
	if p.Match(token.TokenType_BANG, token.TokenType_MINUS, token.TokenType_NOT) {
		operator := p.Previous()
		right, err := p.Unary()

		if err != nil {
			return nil, err
		}

		return &ast.UnaryExpr{
			Operator: operator,
			Right:    right,
		}, nil
	}

	// try safe: ?expr
	if p.Match(token.TokenType_QUESTION) {
		question := p.Previous()
		expr, err := p.Unary() // pega função, get, index, etc
		if err != nil {
			return nil, err
		}
		return &ast.SafeExpr{Name: question, Expr: expr}, nil
	}

	return p.Call()
}

func (p *Parser) Call() (ast.Expr, error) {
	expr, err := p.Primary()

	if err != nil {
		return nil, err
	}

	for {
		if p.Match(token.TokenType_LEFT_PAREN) {
			exprCall, err := p.FinishCall(expr)
			if err != nil {
				return nil, err
			}
			expr = exprCall
		} else if p.Match(token.TokenType_DOT) {
			token, err := p.Consume(token.TokenType_IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = &ast.GetExpr{
				Object: expr,
				Name:   token,
			}
		} else if p.Match(token.TokenType_LEFT_BRACKET) {
			index, err := p.Expression()
			if err != nil {
				return nil, err
			}
			p.Consume(token.TokenType_RIGHT_BRACKET, "Expect ']' after index.")
			expr = &ast.IndexExpr{List: expr, Index: index}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *Parser) FinishCall(callee ast.Expr) (ast.Expr, error) {
	var arguments []ast.Expr
	if !p.Check(token.TokenType_RIGHT_PAREN) {
		for {
			arg, err := p.Expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)

			if !p.Match(token.TokenType_COMMA) {
				break
			}
		}
	}

	paren, err := p.Consume(token.TokenType_RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &ast.CallExpr{
		Callee:      callee,
		Parenthesis: paren,
		Arguments:   arguments,
	}, nil
}

func (p *Parser) Primary() (ast.Expr, error) {
	if p.Match(token.TokenType_FALSE) {
		return &ast.LiteralExpr{Value: false}, nil
	}

	if p.Match(token.TokenType_TRUE) {
		return &ast.LiteralExpr{Value: true}, nil
	}

	if p.Match(token.TokenType_NIL) {
		return &ast.LiteralExpr{Value: nil}, nil
	}

	if p.Match(token.TokenType_NUMBER, token.TokenType_STRING) {
		return &ast.LiteralExpr{Value: p.Previous().Literal}, nil
	}

	if p.Match(token.TokenType_LEFT_BRACE) {
		var pairs []ast.DictPair
		for !p.Check(token.TokenType_RIGHT_BRACE) && !p.IsAtEnd() {
			key, err := p.Expression()
			if err != nil {
				return nil, err
			}
			p.Consume(token.TokenType_COLON, "Expect ':' after key.")
			value, err := p.Expression()
			pairs = append(pairs, ast.DictPair{Key: key, Value: value})

			if !p.Match(token.TokenType_COMMA) {
				break
			}
		}
		p.Consume(token.TokenType_RIGHT_BRACE, "Expect '}' after dictionary.")
		return &ast.DictExpr{Pairs: pairs}, nil
	}

	if p.Match(token.TokenType_LEFT_BRACKET) {
		var elements []ast.Expr

		if !p.Check(token.TokenType_RIGHT_BRACKET) {
			for {
				expr, err := p.Expression()

				if err != nil {
					return nil, err
				}

				elements = append(elements, expr)

				if !p.Match(token.TokenType_COMMA) {
					break
				}
			}
		}

		closing, _ := p.Consume(token.TokenType_RIGHT_BRACKET, "Expect ']' after list elements.")

		return &ast.ListExpr{
			Elements: elements,
			Bracket:  closing,
		}, nil
	}

	if p.Match(token.TokenType_SUPER) {
		keyword := p.Previous()
		p.Consume(token.TokenType_DOT, "Expect '.' after 'super'.")
		if method, err := p.Consume(token.TokenType_IDENTIFIER, "Expect superclass method name."); err == nil {
			return &ast.SuperExpr{
				Keyword: keyword,
				Method:  method,
			}, nil
		} else {
			return nil, err
		}
	}

	if p.Match(token.TokenType_SELF) {
		return &ast.SelfExpr{Keyword: p.Previous()}, nil
	}

	if p.Match(token.TokenType_IDENTIFIER) {
		t := p.Previous()
		return &ast.VariableExpr{Name: t}, nil
	}

	if p.Match(token.TokenType_LEFT_PAREN) {
		expr, err := p.Expression()

		if err != nil {
			return nil, err
		}

		p.Consume(token.TokenType_RIGHT_PAREN, "Expect ')' after expression.")
		return &ast.GroupingExpr{Expression: expr}, nil
	}

	return nil, ParserError{
		Token:   p.Peek(),
		Message: "Expect expression.",
	}
}

func (p *Parser) Consume(tt token.TokenType, msg string) (*token.Token, error) {
	if p.Check(tt) {
		next := p.Advance()
		return next, nil
	}

	return nil, ParserError{
		Token:   p.Peek(),
		Message: msg,
	}
}

func (p *Parser) Match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.Check(t) {
			p.Advance()
			return true
		}
	}
	return false
}

func (p *Parser) Advance() *token.Token {
	if !p.IsAtEnd() {
		p.current++
	}
	return p.Previous()
}

func (p *Parser) Check(t token.TokenType) bool {
	if p.IsAtEnd() {
		return false
	}
	return p.Peek().Type == t
}

func (p *Parser) IsAtEnd() bool {
	return p.Peek().Type == token.TokenType_EOF
}

// Retorna token atual que ainda não foi consumido
func (p *Parser) Peek() *token.Token {
	return p.tokens[p.current]
}

// Retorna o token consumido mais recentemente
func (p *Parser) Previous() *token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) Synchronize() {
	p.Advance()
	for !p.IsAtEnd() {
		if p.Previous().Type == token.TokenType_SEMICOLON {
			return
		}

		switch p.Peek().Type {
		case
			token.TokenType_CLASS,
			token.TokenType_FUNC,
			token.TokenType_LET,
			token.TokenType_IF,
			token.TokenType_FOR,
			token.TokenType_WHILE,
			token.TokenType_PRINT,
			token.TokenType_RETURN:
			return
		}

		p.Advance()
	}
}
