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

func (f *Function) Bind(instance *Instance) *Function {
	env := NewEnvironment(f.runtime, f.closure)
	env.Define("self", instance)
	bound := &Function{
		runtime:       f.runtime,
		Declaration:   f.Declaration,
		closure:       env,
		IsInitializer: f.IsInitializer,
	}
	return bound
}

func (f *Function) String() string {
	return fmt.Sprintf("<function %s>", f.Declaration.Name.Lexeme)
}
