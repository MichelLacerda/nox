package main

// Adicionado para logs

type Environment struct {
	runtime   *Nox
	Values    map[string]any
	Enclosing *Environment
}

func (e *Environment) String() string {
	return "<environment>"
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

// Busca uma variável por nome, ignorando o token (usado para self em inicializadores)
func (e *Environment) GetByName(name string) any {
	if value, exists := e.Values[name]; exists {
		return value
	}
	if e.Enclosing != nil {
		return e.Enclosing.GetByName(name)
	}
	e.runtime.ReportRuntimeError(&Token{Lexeme: name, line: 0}, "Undefined variable: "+name)
	return nil
}

func (e *Environment) GetAt(distance int, name string) any {
	ancestor := e.Ancestor(distance)
	if ancestor == nil {
		return nil
	}
	value, exists := ancestor.Values[name]
	if exists {
		return value
	}
	return nil
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
		if env.Enclosing == nil {
			return nil
		}
		env = env.Enclosing
	}
	return env
}
