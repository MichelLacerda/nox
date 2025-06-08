package main

import (
	"fmt"
	"strings"
)

func (a *Assign) String() string {
	return fmt.Sprintf("assign %s = %s", a.Name, a.Value.String())
}

func (b *Binary) String() string {
	return parenthesize(b.Operator.lexeme, b.Left, b.Right)
}

func (c *Call) String() string {
	var args []string
	for _, arg := range c.Arguments {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("call %s(%s)", c.Callee.String(), strings.Join(args, ", "))
}

func (g *Get) String() string {
	return fmt.Sprintf("get %s from %s", g.Name, g.Object.String())
}

func (g *Grouping) String() string {
	return parenthesize("group", g.Expression)
}

func (l *Literal) String() string {
	if l.Value == nil {
		return "nil"
	}

	switch l.Value.(type) {
	case string:
		return fmt.Sprintf("%q", l.Value)
	default:
		return fmt.Sprintf("%v", l.Value)
	}
}

func (l *Logical) String() string {
	return parenthesize(l.Operator.lexeme, l.Left, l.Right)
}

func (s *Set) String() string {
	return parenthesize("set "+s.Name, s.Object, s.Value)
}

func (s *Super) String() string {
	return fmt.Sprintf("super %s", s.Method)
}

func (t *This) String() string {
	return "this"
}

func (u *Unary) String() string {
	return parenthesize(u.Operator.lexeme, u.Right)
}

func (v *Variable) String() string {
	return v.Name
}

func parenthesize(name string, parts ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, part := range parts {
		builder.WriteString(" ")
		builder.WriteString(part.String())
	}

	builder.WriteString(")")
	return builder.String()
}
