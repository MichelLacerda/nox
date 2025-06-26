package ast

type ExprVisitor interface {
	VisitLiteralExpr(expr *LiteralExpr) any
	VisitGroupingExpr(expr *GroupingExpr) any
	VisitUnaryExpr(expr *UnaryExpr) any
	VisitBinaryExpr(expr *BinaryExpr) any
	VisitVariableExpr(expr *VariableExpr) any
	VisitAssignExpr(expr *AssignExpr) any
	VisitCallExpr(expr *CallExpr) any
	VisitGetExpr(expr *GetExpr) any
	VisitSetExpr(expr *SetExpr) any
	VisitLogicalExpr(expr *LogicalExpr) any
	VisitSuperExpr(expr *SuperExpr) any
	VisitSelfExpr(expr *SelfExpr) any
	VisitListExpr(expr *ListExpr) any
	VisitIndexExpr(expr *IndexExpr) any
	VisitSetIndexExpr(expr *SetIndexExpr) any
	VisitDictExpr(expr *DictExpr) any
	VisitSafeExpr(expr *SafeExpr) any
}

func (a *AssignExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitAssignExpr(a)
}

func (a *BinaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinaryExpr(a)
}

func (c *CallExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitCallExpr(c)
}

func (g *GetExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitGetExpr(g)
}

func (g *GroupingExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitGroupingExpr(g)
}

func (l *LiteralExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteralExpr(l)
}

func (l *LogicalExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLogicalExpr(l)
}

func (s *SetExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitSetExpr(s)
}

func (s *SetIndexExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitSetIndexExpr(s)
}

func (s *SuperExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitSuperExpr(s)
}

func (t *SelfExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitSelfExpr(t)
}

func (u *UnaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnaryExpr(u)
}

func (v *VariableExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitVariableExpr(v)
}

func (l *ListExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitListExpr(l)
}

func (i *IndexExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitIndexExpr(i)
}

func (d *DictExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitDictExpr(d)
}

func (s *SafeExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitSafeExpr(s)
}
