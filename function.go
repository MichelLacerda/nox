package main

import (
	"fmt"
)

type Function struct {
	runtime       *Nox
	Declaration   *FunctionStmt
	closure       *Environment
	IsInitializer bool   // Indica se é um inicializador de classe
	SelfToken     *Token // Token de self do parser
}

func NewFunction(r *Nox, declaration *FunctionStmt, closure *Environment, isInitializer bool, selfToken *Token) *Function {
	return &Function{
		runtime:       r,
		Declaration:   declaration,
		closure:       closure,
		IsInitializer: isInitializer,
		SelfToken:     selfToken,
	}
}

func (f *Function) Call(i *Interpreter, args []any) (result any) {
	environment := NewEnvironment(f.runtime, f.closure)
	for i, token := range f.Declaration.Parameters {
		environment.Define(token.Lexeme, args[i])
	}

	defer func() {
		if r := recover(); r != nil {
			if ret, ok := r.(Return); ok {
				if f.IsInitializer {
					// Ignora qualquer valor retornado explicitamente em init
					// result = environment.Get(f.SelfToken)
					result = f.closure.GetAt(0, "self") // Retorna o próprio objeto
					return
				}
				result = ret.Value
			} else {
				panic(r)
			}
		}
	}()
	// Executa o corpo da função
	i.executeBlock(f.Declaration.Body, environment)

	// if f.IsInitializer {
	// 	// Sempre retorna self, mesmo sem return explícito
	// 	return environment.Get(f.SelfToken)
	// }

	if f.IsInitializer {
		// Se for um inicializador, retorna o próprio objeto
		return f.closure.GetAt(0, "self")
	}
	return nil
}

func (f *Function) Arity() int {
	return len(f.Declaration.Parameters)
}

func (n *Function) Bind(i *Instance) *Function {
	env := NewEnvironment(n.runtime, n.closure)
	env.Define("self", i)
	return NewFunction(n.runtime, n.Declaration, env, n.IsInitializer, n.SelfToken)
}

func (f *Function) String() string {
	return fmt.Sprintf("<function %s>", f.Declaration.Name.Lexeme)
}
