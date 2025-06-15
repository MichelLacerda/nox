package main

type MethodType map[string]*Function

type Class struct {
	Name    string
	Methods MethodType
}

func NewClass(name string, methods MethodType) *Class {
	return &Class{
		Name:    name,
		Methods: methods,
	}
}

func (c *Class) Call(i *Interpreter, args []any) any {
	instance := NewInstance(c)
	if initializer, exists := c.FindMethod("init"); exists {
		bound := initializer.Bind(instance)
		// Execute e capture o retorno corretamente
		bound.Call(i, args)
	}
	return instance
}

func (c *Class) Arity() int {
	if initializer, exists := c.FindMethod("init"); exists {
		return initializer.Arity()
	}

	return 0
}

func (c *Class) String() string {
	return "<class " + c.Name + ">"
}

func (c *Class) Bind(instance *Instance) Callable {
	return c
}

func (c *Class) FindMethod(name string) (*Function, bool) {
	method, exists := c.Methods[name]
	return method, exists
}
