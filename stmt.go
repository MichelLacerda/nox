package main

type StmtVisitor interface {
	VisitBlockStmt(stmt *BlockStmt) any
	VisitClassStmt(stmt *ClassStmt) any
	VisitExpressionStmt(stmt *ExpressionStmt) any
	VisitFunctionStmt(stmt *FunctionStmt) any
	VisitIfStmt(stmt *IfStmt) any
	VisitPrintStmt(stmt *PrintStmt) any
	VisitReturnStmt(stmt *ReturnStmt) any
	VisitVarStmt(stmt *VarStmt) any
	VisitWhileStmt(stmt *WhileStmt) any
	VisitForInStmt(stmt *ForInStmt) any
	VisitBreakStmt(stmt *BreakStmt) any
	VisitContinueStmt(stmt *ContinueStmt) any
	VisitWithStmt(stmt *WithStmt) any
}

type Stmt interface {
	String() string
	Accept(visitor StmtVisitor) any
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

// ClassStmt represents a class declaration statement in the language.
type ClassStmt struct {
	Name       *Token
	Superclass *VariableExpr
	Methods    []*FunctionStmt
}

func NewClassStmt(name *Token, superclass *VariableExpr, methods []*FunctionStmt) *ClassStmt {
	return &ClassStmt{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}
}

// Accept implements Stmt.
func (c *ClassStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitClassStmt(c)
}

// String implements Stmt.
func (c *ClassStmt) String() string {
	result := "class " + c.Name.Lexeme
	// if c.Superclass != nil {
	// 	result += " < " + c.Superclass.Name.Lexeme
	// }
	result += " {\n"
	for _, method := range c.Methods {
		result += method.String() + "\n"
	}
	result += "}"
	return result
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

// FunctionStmt represents a function declaration statement in the language.
type FunctionStmt struct {
	Name       *Token
	Parameters []*Token
	Body       []Stmt
}

func (f *FunctionStmt) String() string {
	// result := "func " + f.Name.Lexeme + "("
	// for i, param := range f.Parameters {
	// 	if i > 0 {
	// 		result += ", "
	// 	}
	// 	result += param.Lexeme
	// }
	// result += ") {\n"
	// for _, stmt := range f.Body {
	// 	result += stmt.String() + "\n"
	// }
	// result += "}"
	return "<function " + f.Name.Lexeme + ">"
}

func (f *FunctionStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitFunctionStmt(f)
}

// IfStmt represents an if statement in the language.
type IfStmt struct {
	Condition Expr
	Then      Stmt
	Else      Stmt // Optional else statement
}

func (i *IfStmt) String() string {
	result := "if " + i.Condition.String() + " {\n" + i.Then.String() + "\n}"
	if i.Else != nil {
		result += " else {\n" + i.Else.String() + "\n}"
	}
	return result
}

func (i *IfStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitIfStmt(i)
}

// PrintStmt represents a print statement in the language.
type PrintStmt struct {
	Expressions []Expr
}

func (p *PrintStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitPrintStmt(p)
}

func (p *PrintStmt) String() string {
	var result string
	for _, expr := range p.Expressions {
		result += expr.String() + " "
	}
	return "print " + result + ";"
}

// ReturnStmt represents a return statement in the language.
type ReturnStmt struct {
	Keyword *Token
	Value   Expr
}

func (r *ReturnStmt) String() string {
	if r.Value != nil {
		return r.Keyword.Lexeme + " " + r.Value.String() + ";"
	}
	return r.Keyword.Lexeme + ";"
}

func (r *ReturnStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitReturnStmt(r)
}

// VarStmt represents a variable declaration statement in the language.
type VarStmt struct {
	Name        *Token
	Initializer Expr
}

func (v *VarStmt) String() string {
	return "let " + v.Name.Lexeme + " = " + v.Initializer.String() + ";"
}

func (v *VarStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitVarStmt(v)
}

// WhileStmt represents a while loop statement in the language.
type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (w *WhileStmt) String() string {
	return "while " + w.Condition.String() + " {\n" + w.Body.String() + "\n}"
}

func (w *WhileStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitWhileStmt(w)
}

// ForInStmt represents a for-in loop statement in the language.
type ForInStmt struct {
	IndexVar *Token // pode ser nil, para o `_`
	ValueVar *Token
	Iterable Expr
	Body     Stmt
}

func (f *ForInStmt) String() string {
	result := "for "
	if f.IndexVar != nil {
		result += f.IndexVar.Lexeme + " in "
	}
	if f.ValueVar != nil {
		result += f.ValueVar.Lexeme + " in "
	}
	result += f.Iterable.String() + " {\n" + f.Body.String() + "\n}"
	return result
}

func (f *ForInStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitForInStmt(f)
}

// ListStmt represents a list of statements in the language.
type ListStmt struct {
	Statements []Stmt
}

func (l *ListStmt) String() string {
	var result string
	for _, stmt := range l.Statements {
		result += stmt.String() + "\n"
	}
	return result
}

// BreakStmt represents a break statement in the language.
type BreakStmt struct {
	Keyword *Token
}

func (b *BreakStmt) String() string {
	return b.Keyword.Lexeme + ";"
}

func (b *BreakStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitBreakStmt(b)
}

// ContinueStmt represents a continue statement in the language.
type ContinueStmt struct {
	Keyword *Token
}

func (c *ContinueStmt) String() string {
	return c.Keyword.Lexeme + ";"
}

func (c *ContinueStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitContinueStmt(c)
}

type WithStmt struct {
	Resource Expr   // ex: open("file.txt", "r")
	Alias    *Token // ex: f
	Body     Stmt   // ex: bloco { ... }
}

func (w *WithStmt) String() string {
	result := "with " + w.Resource.String()
	if w.Alias != nil {
		result += " as " + w.Alias.Lexeme
	}
	result += " {\n" + w.Body.String() + "\n}"
	return result
}

func (w *WithStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitWithStmt(w)
}
