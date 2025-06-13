package main

import (
	"fmt"
	"time"
)

type NoxCallable interface {
	Call(interpreter *Interpreter, args []any) any
	Arity() int
	String() string
}

type NoxFunction struct {
	Declaration *FunctionStmt
	closure     *Environment
}

func NewNoxFunction(declaration *FunctionStmt, closure *Environment) NoxCallable {
	return &NoxFunction{
		Declaration: declaration,
		closure:     closure,
	}
}

func (f *NoxFunction) Call(interpreter *Interpreter, args []any) (result any) {
	environment := NewEnvironment(interpreter.runtime, f.closure)
	for i, token := range f.Declaration.Parameters {
		environment.Define(token.Lexeme, args[i])
	}

	defer func() {
		if r := recover(); r != nil {
			if returnValue, ok := r.(Return); ok {
				result = returnValue.Value
			} else {
				panic(r) // repropaga se não for um Return
			}
		}
	}()

	// Executa o corpo da função
	interpreter.executeBlock(f.Declaration.Body, environment)
	return nil
}

func (f *NoxFunction) Arity() int {
	return len(f.Declaration.Parameters)
}

func (f *NoxFunction) String() string {
	return fmt.Sprintf("<function %s>", f.Declaration.Name.Lexeme)
}

// ===== Native functions =====

type ClockCallable struct{}

func (c ClockCallable) Arity() int {
	return 0
}

func (c ClockCallable) Call(interpreter *Interpreter, arguments []any) any {
	return float64(time.Now().UnixNano()) / 1e9 // segundos com fração
}

func (c ClockCallable) String() string {
	return "<native fn>"
}
