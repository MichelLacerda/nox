package runtime

import "github.com/MichelLacerda/nox/internal/ast"

func (r *Resolver) VisitDictExpr(expr *ast.DictExpr) any {
	for _, pair := range expr.Pairs {
		r.ResolveExpr(pair.Key)
		r.ResolveExpr(pair.Value)
	}
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *ast.VariableExpr) any {
	if !r.scopes.IsEmpty() {
		scope, _ := r.scopes.Peek()
		if declared, exists := scope[expr.Name.Lexeme]; exists && !declared {
			r.interpreter.Runtime.ReportRuntimeError(expr.Name, "Cannot read local variable in its own initializer.")
		}
	}

	r.ResolveLocalExpr(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr *ast.AssignExpr) any {
	r.ResolveExpr(expr.Value)
	r.ResolveLocalExpr(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitIndexExpr(expr *ast.IndexExpr) any {
	r.ResolveExpr(expr.List)
	r.ResolveExpr(expr.Index)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *ast.BinaryExpr) any {
	r.ResolveExpr(expr.Left)
	r.ResolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *ast.CallExpr) any {
	r.ResolveExpr(expr.Callee)
	for _, arg := range expr.Arguments {
		r.ResolveExpr(arg)
	}
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *ast.GroupingExpr) any {
	r.ResolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *ast.LiteralExpr) any {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *ast.LogicalExpr) any {
	r.ResolveExpr(expr.Left)
	r.ResolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *ast.UnaryExpr) any {
	r.ResolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitGetExpr(expr *ast.GetExpr) any {
	r.ResolveExpr(expr.Object)
	return nil
}

func (r *Resolver) VisitSetExpr(expr *ast.SetExpr) any {
	r.ResolveExpr(expr.Object)
	r.ResolveExpr(expr.Value)
	return nil
}

func (r *Resolver) VisitSuperExpr(expr *ast.SuperExpr) any {
	if r.currentClass == ClassTypeNone {
		r.interpreter.Runtime.ReportRuntimeError(expr.Keyword, "Cannot use 'super' outside of a class.")
	} else if r.currentClass != ClassTypeSubclass {
		r.interpreter.Runtime.ReportRuntimeError(expr.Keyword, "Cannot use 'super' in a class with no superclass.")
	}
	r.ResolveLocalExpr(expr, expr.Keyword)
	return nil
}

func (r *Resolver) VisitSelfExpr(expr *ast.SelfExpr) any {
	if r.currentClass == ClassTypeNone {
		r.interpreter.Runtime.ReportRuntimeError(expr.Keyword, "Cannot use 'self' outside of a class.")
		return nil
	}

	r.ResolveLocalExpr(expr, expr.Keyword)
	return nil
}

func (r *Resolver) VisitListExpr(expr *ast.ListExpr) any {
	for _, element := range expr.Elements {
		r.ResolveExpr(element)
	}
	return nil
}

func (r *Resolver) VisitSafeExpr(expr *ast.SafeExpr) any {
	r.ResolveExpr(expr.Expr)
	return nil
}
