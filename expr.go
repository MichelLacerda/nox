package main

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
}
type Expr interface {
	String() string
	Accept(visitor ExprVisitor) any
}

type AssignExpr struct {
	Name  *Token
	Value Expr
}

func (a *AssignExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitAssignExpr(a)
}

type BinaryExpr struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func (a *BinaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinaryExpr(a)
}

type CallExpr struct {
	Callee      Expr
	Parenthesis *Token // The opening parenthesis
	Arguments   []Expr
}

func (c *CallExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitCallExpr(c)
}

type GetExpr struct {
	Object Expr
	Name   *Token
}

func (g *GetExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitGetExpr(g)
}

type GroupingExpr struct {
	Expression Expr
}

func (g *GroupingExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitGroupingExpr(g)
}

type LiteralExpr struct {
	Value any
}

func (l *LiteralExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteralExpr(l)
}

type LogicalExpr struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func (l *LogicalExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLogicalExpr(l)
}

type SetExpr struct {
	Object Expr
	Name   *Token
	Value  Expr
}

func (s *SetExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitSetExpr(s)
}

type SuperExpr struct {
	Keyword *Token // The 'super' keyword
	Method  *Token
}

func (s *SuperExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitSuperExpr(s)
}

type SelfExpr struct {
	Keyword *Token // The 'self' keyword
}

func (t *SelfExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitSelfExpr(t)
}

type UnaryExpr struct {
	Operator *Token
	Right    Expr
}

func (u *UnaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnaryExpr(u)
}

type VariableExpr struct {
	Name *Token
}

func (v *VariableExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitVariableExpr(v)
}

type ListExpr struct {
	Elements []Expr
	Bracket  *Token
}

func (l *ListExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitListExpr(l)
}

type IndexExpr struct {
	List  Expr
	Index Expr
}

func (i *IndexExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitIndexExpr(i)
}
