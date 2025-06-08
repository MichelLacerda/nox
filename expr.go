package main

type Expr interface {
	String() string
}

type Assign struct {
	Name  string
	Value Expr
}

type Binary struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

type Call struct {
	Callee    Expr
	Arguments []Expr
}

type Get struct {
	Object Expr
	Name   string
}

type Grouping struct {
	Expression Expr
}

type Literal struct {
	Value any
}

type Logical struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

type Set struct {
	Object Expr
	Name   string
	Value  Expr
}

type Super struct {
	Method string
}

type This struct{}

type Unary struct {
	Operator *Token
	Right    Expr
}

type Variable struct {
	Name string
}
