package main

type Resolver struct {
	interpreter     *Interpreter
	scopes          ResolverStack
	currentFunction FunctionType
}

type FunctionType int

const (
	FunctionTypeNone FunctionType = iota
	FunctionTypeFunction
)

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          ResolverStack{},
		currentFunction: FunctionTypeNone,
	}
}

func (r *Resolver) ResolveStatements(statements []Stmt) {
	for _, s := range statements {
		s.Accept(r)
	}
}

func (r *Resolver) ResolveStatement(s Stmt) {
	s.Accept(r)
}

func (r *Resolver) ResolveExpr(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) BeginScope() {
	scope := map[string]bool{}
	r.scopes.Push(scope)
}

func (r *Resolver) EndScope() {
	r.scopes.Pop()
}

func (r *Resolver) Declare(name *Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope, _ := r.scopes.Peek()

	if _, exists := (*scope)[name.Lexeme]; exists {
		r.interpreter.runtime.ReportRuntimeError(name, "Variable already defined: "+name.Lexeme)
		return
	}

	(*scope)[name.Lexeme] = false
}

func (r *Resolver) Define(name *Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope, _ := r.scopes.Peek()
	(*scope)[name.Lexeme] = true
}

func (r *Resolver) ResolveLocalExpr(expr Expr, name *Token) {
	if r.scopes.IsEmpty() {
		return
	}

	for i := len(r.scopes) - 1; i >= 0; i-- {
		scope := r.scopes[i]
		if _, exists := scope[name.Lexeme]; exists {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
	// r.interpreter.runtime.ReportRuntimeError(name, "Undefined variable: "+name.Lexeme)
}

func (r *Resolver) ResolveFunction(stmt *FunctionStmt, functionType FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = functionType
	r.BeginScope()
	for _, param := range stmt.Parameters {
		r.Declare(param)
		r.Define(param)
	}
	r.ResolveStatements(stmt.Body)
	r.EndScope()
	r.currentFunction = enclosingFunction
}

// ==== Visitor methods ====

func (r *Resolver) VisitBlockStmt(stmt *BlockStmt) any {
	r.BeginScope()
	r.ResolveStatements(stmt.Statements)
	r.EndScope()
	return nil
}

func (r *Resolver) VisitVarStmt(stmt *VarStmt) any {
	r.Declare(stmt.Name)
	if stmt.Initializer != nil {
		r.ResolveExpr(stmt.Initializer)
	}
	r.Define(stmt.Name)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *VariableExpr) any {
	if !r.scopes.IsEmpty() {
		scope, _ := r.scopes.Peek()
		if declared, exists := (*scope)[expr.Name.Lexeme]; exists && !declared {
			r.interpreter.runtime.ReportRuntimeError(expr.Name, "Cannot read local variable in its own initializer.")
		}
	}

	r.ResolveLocalExpr(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr *AssignExpr) any {
	r.ResolveExpr(expr.Value)
	r.ResolveLocalExpr(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *FunctionStmt) any {
	r.Declare(stmt.Name)
	r.Define(stmt.Name)
	r.ResolveFunction(stmt, FunctionTypeFunction)
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ExpressionStmt) any {
	r.ResolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *IfStmt) any {
	r.ResolveExpr(stmt.Condition)
	r.ResolveStatement(stmt.Then)
	if stmt.Else != nil {
		r.ResolveStatement(stmt.Else)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *PrintStmt) any {
	r.ResolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *ReturnStmt) any {
	if r.currentFunction == FunctionTypeNone {
		r.interpreter.runtime.ReportRuntimeError(stmt.Keyword, "Cannot return from top-level code.")
		return nil
	}

	if stmt.Value != nil {
		r.ResolveExpr(stmt.Value)
	}
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *WhileStmt) any {
	r.ResolveExpr(stmt.Condition)
	r.ResolveStatement(stmt.Body)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *BinaryExpr) any {
	r.ResolveExpr(expr.Left)
	r.ResolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *CallExpr) any {
	r.ResolveExpr(expr.Callee)
	for _, arg := range expr.Arguments {
		r.ResolveExpr(arg)
	}
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *GroupingExpr) any {
	r.ResolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *LiteralExpr) any {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *LogicalExpr) any {
	r.ResolveExpr(expr.Left)
	r.ResolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *UnaryExpr) any {
	r.ResolveExpr(expr.Right)
	return nil
}

// VisitGetExpr implements ExprVisitor.
func (r *Resolver) VisitGetExpr(expr *GetExpr) any {
	panic("unimplemented")
}

// VisitSetExpr implements ExprVisitor.
func (r *Resolver) VisitSetExpr(expr *SetExpr) any {
	panic("unimplemented")
}

// VisitSuperExpr implements ExprVisitor.
func (r *Resolver) VisitSuperExpr(expr *SuperExpr) any {
	panic("unimplemented")
}

// VisitThisExpr implements ExprVisitor.
func (r *Resolver) VisitThisExpr(expr *ThisExpr) any {
	panic("unimplemented")
}
