package main

type Environment struct {
	runtime   *Nox
	Values    map[string]any
	Enclosing *Environment
}

func NewEnvironment(r *Nox, scope *Environment) *Environment {
	return &Environment{
		runtime:   r,
		Values:    map[string]any{},
		Enclosing: scope,
	}
}

func (e *Environment) Define(name string, value any) {
	if _, exists := e.Values[name]; exists {
		//panic("Variable already defined: " + name)
		e.runtime.ReportRuntimeError(&Token{Lexeme: name}, "Variable already defined: "+name)
	}
	e.Values[name] = value
}

func (e *Environment) Get(t *Token) any {
	if value, exists := e.Values[t.Lexeme]; exists {
		return value
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(t)
	}

	e.runtime.ReportRuntimeError(&Token{Lexeme: t.Lexeme}, "Undefined variable: "+t.Lexeme)
	return nil
}

func (e *Environment) GetAt(distance int, name string) any {
	ancestor := e.Ancestor(distance)
	return ancestor.Values[name]
}

func (e *Environment) Assign(name *Token, value any) {
	if _, exists := e.Values[name.Lexeme]; exists {
		e.Values[name.Lexeme] = value
		return
	}

	if e.Enclosing != nil {
		e.Enclosing.Assign(name, value)
		return
	}

	e.runtime.ReportRuntimeError(name, "Undefined variable: "+name.Lexeme)
}

func (e *Environment) AssignAt(d int, name *Token, value any) {
	e.Ancestor(d).Values[name.Lexeme] = value
}

func (e *Environment) Ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.Enclosing
	}
	return env
}
