package main

import (
	"fmt"
	"strings"
)

func (a *AssignExpr) String() string {
	return fmt.Sprintf("assign %s = %s", a.Name, a.Value.String())
}

func (b *BinaryExpr) String() string {
	return parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (c *CallExpr) String() string {
	var args []string
	for _, arg := range c.Arguments {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("call %s(%s)", c.Callee.String(), strings.Join(args, ", "))
}

func (g *GetExpr) String() string {
	return fmt.Sprintf("get %s from %s", g.Name, g.Object.String())
}

func (g *GroupingExpr) String() string {
	return parenthesize("group", g.Expression)
}

func (l *LiteralExpr) String() string {
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

func (l *LogicalExpr) String() string {
	return parenthesize(l.Operator.Lexeme, l.Left, l.Right)
}

func (s *SetExpr) String() string {
	return parenthesize("set "+s.Name.Lexeme, s.Object, s.Value)
}

func (s *SuperExpr) String() string {
	return fmt.Sprintf("super %s", s.Method)
}

func (t *SelfExpr) String() string {
	return "self"
}

func (u *UnaryExpr) String() string {
	return parenthesize(u.Operator.Lexeme, u.Right)
}

func (v *VariableExpr) String() string {
	return v.Name.Lexeme
}

func (l *ListExpr) String() string {
	var elements []string
	for _, elem := range l.Elements {
		elements = append(elements, elem.String())
	}
	return fmt.Sprintf("list[%s]", strings.Join(elements, ", "))
}

func (i *IndexExpr) String() string {
	return fmt.Sprintf("index %s[%s]", i.List.String(), i.Index.String())
}

func (i *DictExpr) String() string {
	var pairs []string
	for _, value := range i.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", value.Key, value.Value))
	}
	return fmt.Sprintf("dict{%s}", strings.Join(pairs, ", "))
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
