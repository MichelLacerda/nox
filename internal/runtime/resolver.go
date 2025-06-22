package runtime

import (
	"github.com/MichelLacerda/nox/internal/ast"
	"github.com/MichelLacerda/nox/internal/parser"
	"github.com/MichelLacerda/nox/internal/token"
)

type ResolverStack []map[string]bool

func (s ResolverStack) IsEmpty() bool {
	if len(s) == 0 {
		return true
	}
	return false
}

func (s *ResolverStack) Push(m map[string]bool) {
	*s = append(*s, m)
}

func (s *ResolverStack) Pop() (map[string]bool, bool) {
	if len(*s) == 0 {
		return nil, false
	}
	index := len(*s) - 1
	elem := (*s)[index]
	*s = (*s)[:index]
	return elem, true
}

func (s *ResolverStack) Peek() (map[string]bool, bool) {
	if len(*s) == 0 {
		return nil, false
	}
	return (*s)[len(*s)-1], true
}

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

func (r *Resolver) ResolveStatements(statements []ast.Stmt) {
	for _, s := range statements {
		s.Accept(r)
	}
}

func (r *Resolver) ResolveStatement(s ast.Stmt) {
	s.Accept(r)
}

func (r *Resolver) ResolveExpr(expr ast.Expr) {
	expr.Accept(r)
}

func (r *Resolver) BeginScope() {
	scope := map[string]bool{}
	r.scopes.Push(scope)
}

func (r *Resolver) EndScope() {
	r.scopes.Pop()
}

func (r *Resolver) Declare(name *token.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope, _ := r.scopes.Peek()

	if _, exists := scope[name.Lexeme]; exists {
		r.interpreter.Runtime.ReportRuntimeError(name, "Variable already defined: "+name.Lexeme)
		return
	}

	scope[name.Lexeme] = false
}

func (r *Resolver) Define(name *token.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope, _ := r.scopes.Peek()
	scope[name.Lexeme] = true
}

func (r *Resolver) ResolveLocalExpr(expr ast.Expr, name *token.Token) {
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

func (r *Resolver) ResolveFunction(stmt *ast.FunctionStmt, functionType FunctionType) {
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

func (r *Resolver) errorToken(token *token.Token, message string) {
	panic(parser.ParserError{Token: token, Message: message})
}
