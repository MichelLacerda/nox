package main

type FunctionType int

const (
	FunctionTypeNone FunctionType = iota
	FunctionTypeFunction
	FunctionTypeInitializer
	FunctionTypeMethod
)

type ClassType int

const (
	ClassTypeNone ClassType = iota
	ClassTypeClass
	ClassTypeSubclass
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          ResolverStack
	currentFunction FunctionType
	currentClass    ClassType
	insideLoop      bool
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          ResolverStack{},
		currentFunction: FunctionTypeNone,
		currentClass:    ClassTypeNone,
		insideLoop:      false,
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

	if _, exists := scope[name.Lexeme]; exists {
		r.interpreter.runtime.ReportRuntimeError(name, "Variable already defined: "+name.Lexeme)
		return
	}

	scope[name.Lexeme] = false
}

func (r *Resolver) Define(name *Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope, _ := r.scopes.Peek()
	scope[name.Lexeme] = true
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

func (r *Resolver) VisitDictExpr(expr *DictExpr) any {
	for _, pair := range expr.Pairs {
		r.ResolveExpr(pair.Key)
		r.ResolveExpr(pair.Value)
	}
	return nil
}

// ==== Visitor methods ====

func (r *Resolver) VisitBlockStmt(stmt *BlockStmt) any {
	r.BeginScope()
	r.ResolveStatements(stmt.Statements)
	r.EndScope()
	return nil
}

func (r *Resolver) VisitClassStmt(stmt *ClassStmt) any {
	enclosingClass := r.currentClass
	r.currentClass = ClassTypeClass

	r.Declare(stmt.Name)
	r.Define(stmt.Name)

	if stmt.Superclass != nil {
		if stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme {
			r.interpreter.runtime.ReportRuntimeError(stmt.Superclass.Name, "A class cannot inherit from itself.")
		}

		r.currentClass = ClassTypeSubclass
		r.ResolveExpr(stmt.Superclass)

		r.BeginScope() // Create a new scope for the superclass
		if s, ok := r.scopes.Peek(); ok {
			s["super"] = true
		}
	}

	if s, ok := r.scopes.Peek(); ok {
		s["self"] = true
	}

	for _, method := range stmt.Methods {
		functionType := FunctionTypeMethod
		if method.Name.Lexeme == "init" {
			functionType = FunctionTypeInitializer
		}

		r.BeginScope()
		if s, ok := r.scopes.Peek(); ok {
			s["self"] = true
		}
		r.ResolveFunction(method, functionType)
		r.EndScope()
	}

	r.BeginScope()
	if s, ok := r.scopes.Peek(); ok {
		s[stmt.Name.Lexeme] = true
	}
	r.EndScope()

	if stmt.Superclass != nil {
		r.EndScope() // End the scope created for the superclass
	}

	r.currentClass = enclosingClass
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
		if declared, exists := scope[expr.Name.Lexeme]; exists && !declared {
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
	for _, expr := range stmt.Expressions {
		r.ResolveExpr(expr)
	}
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *ReturnStmt) any {
	if r.currentFunction == FunctionTypeNone {
		r.interpreter.runtime.ReportRuntimeError(stmt.Keyword, "Cannot return from top-level code.")
		return nil
	}

	if stmt.Value != nil {
		if r.currentFunction == FunctionTypeInitializer {
			r.interpreter.runtime.ReportRuntimeError(stmt.Keyword, "Cannot return a value from an initializer.")
		}
		r.ResolveExpr(stmt.Value)
	}
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *WhileStmt) any {
	wasInside := r.insideLoop
	r.insideLoop = true

	r.ResolveExpr(stmt.Condition)
	r.ResolveStatement(stmt.Body)

	r.insideLoop = wasInside
	return nil
}

func (r *Resolver) VisitForInStmt(stmt *ForInStmt) any {
	wasInside := r.insideLoop
	r.insideLoop = true

	r.ResolveExpr(stmt.Iterable)

	r.BeginScope()
	if stmt.IndexVar != nil {
		r.Declare(stmt.IndexVar)
		r.Define(stmt.IndexVar)
	}
	if stmt.ValueVar != nil {
		r.Declare(stmt.ValueVar)
		r.Define(stmt.ValueVar)
	}
	r.ResolveStatement(stmt.Body)
	r.EndScope()

	r.insideLoop = wasInside
	return nil
}

func (r *Resolver) VisitIndexExpr(expr *IndexExpr) any {
	r.ResolveExpr(expr.List)
	r.ResolveExpr(expr.Index)
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

func (r *Resolver) VisitGetExpr(expr *GetExpr) any {
	r.ResolveExpr(expr.Object)
	return nil
}

func (r *Resolver) VisitSetExpr(expr *SetExpr) any {
	r.ResolveExpr(expr.Object)
	r.ResolveExpr(expr.Value)
	return nil
}

func (r *Resolver) VisitSuperExpr(expr *SuperExpr) any {
	if r.currentClass == ClassTypeNone {
		r.interpreter.runtime.ReportRuntimeError(expr.Keyword, "Cannot use 'super' outside of a class.")
	} else if r.currentClass != ClassTypeSubclass {
		r.interpreter.runtime.ReportRuntimeError(expr.Keyword, "Cannot use 'super' in a class with no superclass.")
	}
	r.ResolveLocalExpr(expr, expr.Keyword)
	return nil
}

func (r *Resolver) VisitSelfExpr(expr *SelfExpr) any {
	if r.currentClass == ClassTypeNone {
		r.interpreter.runtime.ReportRuntimeError(expr.Keyword, "Cannot use 'self' outside of a class.")
		return nil
	}

	r.ResolveLocalExpr(expr, expr.Keyword)
	return nil
}

func (r *Resolver) VisitListExpr(expr *ListExpr) any {
	for _, element := range expr.Elements {
		r.ResolveExpr(element)
	}
	return nil
}

func (r *Resolver) VisitSafeExpr(expr *SafeExpr) any {
	r.ResolveExpr(expr.Expr)
	return nil
}

func (r *Resolver) VisitBreakStmt(stmt *BreakStmt) any {
	if !r.insideLoop {
		r.errorToken(stmt.Keyword, "Can't use 'break' outside of a loop.")
	}
	return nil
}

func (r *Resolver) VisitContinueStmt(stmt *ContinueStmt) any {
	if !r.insideLoop {
		r.errorToken(stmt.Keyword, "Can't use 'continue' outside of a loop.")
	}
	return nil
}

func (r *Resolver) VisitWithStmt(stmt *WithStmt) any {
	r.ResolveExpr(stmt.Resource)

	r.BeginScope()
	r.Declare(stmt.Alias)
	r.Define(stmt.Alias)
	r.ResolveStatement(stmt.Body)
	r.EndScope()

	return nil
}

func (r *Resolver) VisitImportStmt(stmt *ImportStmt) any {
	return nil
}

func (r *Resolver) VisitExportStmt(stmt *ExportStmt) any {
	r.ResolveStatement(stmt.Declaration)
	return nil
}

func (r *Resolver) errorToken(token *Token, message string) {
	panic(ParserError{Token: token, Message: message})
}
