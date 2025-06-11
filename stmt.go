package main

type StmtVisitor interface {
	VisitBlockStmt(stmt *BlockStmt) any
	//VisitClassStmt(stmt *ClassStmt) any
	VisitExpressionStmt(stmt *ExpressionStmt) any
	//VisitFunctionStmt(stmt *FunctionStmt) any
	//VisitIfStmt(stmt *IfStmt) any
	VisitPrintStmt(stmt *PrintStmt) any
	//VisitReturnStmt(stmt *ReturnStmt) any
	VisitVarStmt(stmt *VarStmt) any
	//VisitWhileStmt(stmt *WhileStmt) any
}

type Stmt interface {
	String() string
	Accept(visitor StmtVisitor) any
}

// VarStmt represents a variable declaration statement in the language.
type ExpressionStmt struct {
	Expression Expr
}

func (e *ExpressionStmt) String() string {
	return e.Expression.String() + ";"
}

func (e *ExpressionStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitExpressionStmt(e)
}

// PrintStmt represents a print statement in the language.
type PrintStmt struct {
	Expression Expr
}

func (p *PrintStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitPrintStmt(p)
}

func (p *PrintStmt) String() string {
	return "print " + p.Expression.String() + ";"
}

// VarStmt represents a variable declaration statement in the language.
type VarStmt struct {
	Name  *Token
	Value Expr
}

func (v *VarStmt) String() string {
	return "let " + v.Name.Lexeme + " = " + v.Value.String() + ";"
}

func (v *VarStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitVarStmt(v)
}

// BlockStmt represents a block of statements in the language.
type BlockStmt struct {
	Statements []Stmt
}

func (b *BlockStmt) String() string {
	var result string
	for _, stmt := range b.Statements {
		result += stmt.String() + "\n"
	}
	return "{\n" + result + "}"
}

func (b *BlockStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitBlockStmt(b)
}
