package main

type Instance struct {
	Class  *Class
	Fields map[string]any
}

func NewInstance(c *Class) *Instance {
	return &Instance{
		Class:  c,
		Fields: map[string]any{},
	}
}

func (i *Instance) String() string {
	return "<instance of " + i.Class.Name + ">"
}

func (i *Instance) Get(name *Token) any {
	if value, exists := i.Fields[name.Lexeme]; exists {
		return value
	}

	if method, exists := i.FindMethod(name.Lexeme); exists {
		return method.Bind(i)
	}

	panic(RuntimeError{
		Token:   name,
		Message: "Undefined property '" + name.Lexeme + "' in instance of class '" + i.Class.Name + "'.",
	})
}

func (i *Instance) FindMethod(name string) (*Function, bool) {
	if method, exists := i.Class.Methods[name]; exists {
		return method, true
	}
	return nil, false
}

// func (i *Instance) Set(name *Token, value any) {
// 	i.Fields[name.Lexeme] = value
// }

func (i *Instance) Set(name *Token, value any) {
	// Protege métodos
	if _, isMethod := i.Class.Methods[name.Lexeme]; isMethod {
		panic(RuntimeError{
			Token:   name,
			Message: "Cannot overwrite method '" + name.Lexeme + "' with a field.",
		})
	}

	i.Fields[name.Lexeme] = value
}
