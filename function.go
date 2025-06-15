package main

import (
	"fmt"
)

type Function struct {
	runtime       *Nox
	Declaration   *FunctionStmt
	closure       *Environment
	IsInitializer bool
}

func NewFunction(r *Nox, declaration *FunctionStmt, closure *Environment, isInitializer bool) *Function {
	return &Function{
		runtime:       r,
		Declaration:   declaration,
		closure:       closure,
		IsInitializer: isInitializer,
	}
}

func (f *Function) Call(i *Interpreter, args []any) (result any) {
	// CORRIGIDO: cria novo escopo para execução
	environment := NewEnvironment(f.runtime, f.closure)

	if len(args) != f.Arity() {
		f.runtime.ReportRuntimeError(f.Declaration.Name, fmt.Sprintf(
			"Expected %d arguments but got %d.", f.Arity(), len(args)))
		return nil
	}

	for idx, param := range f.Declaration.Parameters {
		environment.Define(param.Lexeme, args[idx])
	}

	defer func() {
		if r := recover(); r != nil {
			if ret, ok := r.(Return); ok {
				if f.IsInitializer {
					result = f.closure.GetAt(0, "self")
					return
				}
				result = ret.Value
			} else {
				panic(r)
			}
		}
	}()

	i.executeBlock(f.Declaration.Body, environment)

	if f.IsInitializer {
		result = f.closure.GetAt(0, "self")
		return result
	}

	return nil
}

func (f *Function) Arity() int {
	return len(f.Declaration.Parameters)
}

func (n *Function) Bind(i *Instance) *Function {
	env := NewEnvironment(n.runtime, n.closure)
	env.Define("self", i)

	return &Function{
		runtime:       n.runtime,
		Declaration:   n.Declaration,
		closure:       env,
		IsInitializer: n.IsInitializer,
	}
}

func (f *Function) String() string {
	return fmt.Sprintf("<function %s>", f.Declaration.Name.Lexeme)
}
