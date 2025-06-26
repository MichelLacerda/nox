package ast

import "github.com/MichelLacerda/nox/internal/token"

type Expr interface {
	String() string
	Accept(visitor ExprVisitor) any
}

type AssignExpr struct {
	Name  *token.Token
	Value Expr
}

type BinaryExpr struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

type CallExpr struct {
	Callee      Expr
	Parenthesis *token.Token // The opening parenthesis
	Arguments   []Expr
}

type GetExpr struct {
	Object Expr
	Name   *token.Token
}

type GroupingExpr struct {
	Expression Expr
}

type LiteralExpr struct {
	Value any
}

type LogicalExpr struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

type SetExpr struct {
	Object Expr
	Name   *token.Token
	Value  Expr
}

type SuperExpr struct {
	Keyword *token.Token // The 'super' keyword
	Method  *token.Token
}

type SelfExpr struct {
	Keyword *token.Token // The 'self' keyword
}

type UnaryExpr struct {
	Operator *token.Token
	Right    Expr
}

type VariableExpr struct {
	Name *token.Token
}

type ListExpr struct {
	Elements []Expr
	Bracket  *token.Token
}

type IndexExpr struct {
	Object Expr
	Index  Expr
}

type SetIndexExpr struct {
	Object Expr
	Index  Expr
	Value  Expr
}

type DictExpr struct {
	Pairs []DictPair
}

type DictPair struct {
	Key   Expr
	Value Expr
}

type SafeExpr struct {
	Expr Expr
	Name *token.Token // Optional name for the safe expression
}
