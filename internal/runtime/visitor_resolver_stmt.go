package runtime

import "github.com/MichelLacerda/nox/internal/ast"

func (r *Resolver) VisitBlockStmt(stmt *ast.BlockStmt) any {
	r.BeginScope()
	r.ResolveStatements(stmt.Statements)
	r.EndScope()
	return nil
}

func (r *Resolver) VisitClassStmt(stmt *ast.ClassStmt) any {
	enclosingClass := r.currentClass
	r.currentClass = ClassTypeClass

	r.Declare(stmt.Name)
	r.Define(stmt.Name)

	if stmt.Superclass != nil {
		if variable, ok := stmt.Superclass.(*ast.VariableExpr); ok {
			if stmt.Name.Lexeme == variable.Name.Lexeme {
				r.interpreter.Runtime.ReportRuntimeError(variable.Name, "A class cannot inherit from itself.")
			}
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

func (r *Resolver) VisitVarStmt(stmt *ast.VarStmt) any {
	r.Declare(stmt.Name)
	if stmt.Initializer != nil {
		r.ResolveExpr(stmt.Initializer)
	}
	r.Define(stmt.Name)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *ast.FunctionStmt) any {
	r.Declare(stmt.Name)
	r.Define(stmt.Name)
	r.ResolveFunction(stmt, FunctionTypeFunction)
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ast.ExpressionStmt) any {
	r.ResolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *ast.IfStmt) any {
	r.ResolveExpr(stmt.Condition)
	r.ResolveStatement(stmt.Then)
	if stmt.Else != nil {
		r.ResolveStatement(stmt.Else)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *ast.PrintStmt) any {
	for _, expr := range stmt.Expressions {
		r.ResolveExpr(expr)
	}
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *ast.ReturnStmt) any {
	if r.currentFunction == FunctionTypeNone {
		r.interpreter.Runtime.ReportRuntimeError(stmt.Keyword, "Cannot return from top-level code.")
		return nil
	}

	if stmt.Value != nil {
		if r.currentFunction == FunctionTypeInitializer {
			r.interpreter.Runtime.ReportRuntimeError(stmt.Keyword, "Cannot return a value from an initializer.")
		}
		r.ResolveExpr(stmt.Value)
	}
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *ast.WhileStmt) any {
	wasInside := r.insideLoop
	r.insideLoop = true

	r.ResolveExpr(stmt.Condition)
	r.ResolveStatement(stmt.Body)

	r.insideLoop = wasInside
	return nil
}

func (r *Resolver) VisitForInStmt(stmt *ast.ForInStmt) any {
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

func (r *Resolver) VisitBreakStmt(stmt *ast.BreakStmt) any {
	if !r.insideLoop {
		r.errorToken(stmt.Keyword, "Can't use 'break' outside of a loop.")
	}
	return nil
}

func (r *Resolver) VisitContinueStmt(stmt *ast.ContinueStmt) any {
	if !r.insideLoop {
		r.errorToken(stmt.Keyword, "Can't use 'continue' outside of a loop.")
	}
	return nil
}

func (r *Resolver) VisitWithStmt(stmt *ast.WithStmt) any {
	r.ResolveExpr(stmt.Resource)

	r.BeginScope()
	r.Declare(stmt.Alias)
	r.Define(stmt.Alias)
	r.ResolveStatement(stmt.Body)
	r.EndScope()

	return nil
}

func (r *Resolver) VisitImportStmt(stmt *ast.ImportStmt) any {
	return nil
}

func (r *Resolver) VisitExportStmt(stmt *ast.ExportStmt) any {
	r.ResolveStatement(stmt.Declaration)
	return nil
}
