package runtime

type MethodType map[string]*Function

type Class struct {
	Name    string
	Methods MethodType
	Super   *Class
}

func NewClass(name string, super *Class, methods MethodType) *Class {
	return &Class{
		Name:    name,
		Super:   super,
		Methods: methods,
	}
}

func (c *Class) Call(i *Interpreter, args []any) any {
	instance := NewInstance(c)
	if initializer, exists := c.FindMethod("init"); exists {
		bound := initializer.Bind(instance)
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
	if method, exists := c.Methods[name]; exists {
		return method, true
	}

	if c.Super != nil {
		return c.Super.FindMethod(name)
	}

	return nil, false
}

func (c *Class) IsInstanceOf(class *Class) any {
	if c == class {
		return true
	}
	if c.Super != nil {
		return c.Super.IsInstanceOf(class)
	}
	return false
}
