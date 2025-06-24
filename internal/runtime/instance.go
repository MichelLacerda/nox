package runtime

import (
	"strings"

	"github.com/MichelLacerda/nox/internal/token"
)

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
	var res strings.Builder
	res.WriteString("<instance of ")
	res.WriteString(i.Class.Name)
	if i.Class.Super != nil {
		res.WriteString("(")
		res.WriteString(i.Class.Super.Name)
		res.WriteString(")")
	}
	res.WriteString(">")
	return res.String()
}

func (i *Instance) Get(name *token.Token) any {
	if value, exists := i.Fields[name.Lexeme]; exists {
		return value
	}

	if method, exists := i.Class.FindMethod(name.Lexeme); exists {
		bound := method.Bind(i)
		return bound
	}

	panic(&RuntimeError{
		Token:   name,
		Message: "Undefined property '" + name.Lexeme + "' in instance of class '" + i.Class.Name + "'.",
	})
}

func (i *Instance) Set(name *token.Token, value any) {
	i.Fields[name.Lexeme] = value
}

func (i *Instance) IsInstanceOf(class *Class) any {
	if i.Class == class {
		return true
	}
	if i.Class.Super != nil {
		return i.Class.Super.IsInstanceOf(class)
	}
	return false
}
