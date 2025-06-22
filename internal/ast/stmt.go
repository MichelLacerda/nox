package ast

import "github.com/MichelLacerda/nox/internal/token"

type Stmt interface {
	String() string
	Accept(visitor StmtVisitor) any
}

type BlockStmt struct {
	Statements []Stmt
}

type ClassStmt struct {
	Name       *token.Token
	Superclass Expr
	Methods    []*FunctionStmt
}

func NewClassStmt(name *token.Token, superclass Expr, methods []*FunctionStmt) *ClassStmt {
	return &ClassStmt{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}
}

type ExpressionStmt struct {
	Expression Expr
}

type FunctionStmt struct {
	Name       *token.Token
	Parameters []*token.Token
	Body       []Stmt
}

type IfStmt struct {
	Condition Expr
	Then      Stmt
	Else      Stmt // Optional else statement
}

type PrintStmt struct {
	Expressions []Expr
}

type ReturnStmt struct {
	Keyword *token.Token
	Value   Expr
}

type VarStmt struct {
	Name        *token.Token
	Initializer Expr
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

type ForInStmt struct {
	IndexVar *token.Token // pode ser nil, para o `_`
	ValueVar *token.Token
	Iterable Expr
	Body     Stmt
}

type ListStmt struct {
	Statements []Stmt
}

type BreakStmt struct {
	Keyword *token.Token
}

type ContinueStmt struct {
	Keyword *token.Token
}

type WithStmt struct {
	Resource Expr         // ex: open("file.txt", "r")
	Alias    *token.Token // ex: f
	Body     Stmt         // ex: block { ... }
}

type ImportStmt struct {
	// import "<path>"" [as <alias>]
	Path  *token.Token // STRING Token, ex: "std/math"
	Alias *token.Token // IDENTIFIER Token, ex: "math"
}

type ExportStmt struct {
	Declaration Stmt // pode ser VarDecl, FuncDecl, ClassDecl
}
