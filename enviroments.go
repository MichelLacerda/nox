package main

type Environment struct {
	runtime *Nox
	Values  map[string]any
	Scope   *Environment
}

func NewEnvironment(r *Nox, scope *Environment) *Environment {
	return &Environment{
		runtime: r,
		Values:  make(map[string]any),
		Scope:   scope,
	}
}

func (e *Environment) Define(name string, value any) {
	if _, exists := e.Values[name]; exists {
		//panic("Variable already defined: " + name)
		e.runtime.ReportRuntimeError(&Token{Lexeme: name}, "Variable already defined: "+name)
	}
	e.Values[name] = value
}

func (e *Environment) Get(name string) any {
	if value, exists := e.Values[name]; exists {
		return value
	}

	if e.Scope != nil {
		return e.Scope.Get(name)
	}

	e.runtime.ReportRuntimeError(&Token{Lexeme: name}, "Undefined variable: "+name)
	return nil
}

func (e *Environment) Assign(name *Token, value any) {
	if _, exists := e.Values[name.Lexeme]; exists {
		e.Values[name.Lexeme] = value
		return
	}

	if e.Scope != nil {
		e.Scope.Assign(name, value)
		return
	}

	e.runtime.ReportRuntimeError(name, "Undefined variable: "+name.Lexeme)
}
