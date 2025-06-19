package main

import (
	"fmt"
	"strings"
)

type ListInstance struct {
	Elements []any
}

func NewListInstance(elements []any) *ListInstance {
	return &ListInstance{Elements: elements}
}

func (l *ListInstance) Get(name *Token) any {
	switch name.Lexeme {
	case "append":
		return &BuiltinFunction{arity: 1, call: func(_ *Interpreter, args []any) any {
			l.Elements = append(l.Elements, args[0])
			return nil
		}}
	case "pop":
		return &BuiltinFunction{arity: 0, call: func(_ *Interpreter, _ []any) any {
			if len(l.Elements) == 0 {
				return nil
			}
			val := l.Elements[len(l.Elements)-1]
			l.Elements = l.Elements[:len(l.Elements)-1]
			return val
		}}
	case "insert":
		return &BuiltinFunction{arity: 2, call: func(_ *Interpreter, args []any) any {
			index, ok1 := args[0].(float64)
			if !ok1 || int(index) < 0 || int(index) > len(l.Elements) {
				return nil
			}
			i := int(index)
			l.Elements = append(l.Elements[:i], append([]any{args[1]}, l.Elements[i:]...)...)
			return nil
		}}
	case "remove":
		return &BuiltinFunction{arity: 1, call: func(_ *Interpreter, args []any) any {
			index, ok := args[0].(float64)
			if !ok || int(index) < 0 || int(index) >= len(l.Elements) {
				return nil
			}
			i := int(index)
			l.Elements = append(l.Elements[:i], l.Elements[i+1:]...)
			return nil
		}}
	case "clear":
		return &BuiltinFunction{arity: 0, call: func(_ *Interpreter, _ []any) any {
			l.Elements = []any{}
			return nil
		}}
	case "length":
		return &BuiltinFunction{arity: 0, call: func(_ *Interpreter, _ []any) any {
			return float64(len(l.Elements))
		}}
	case "contains":
		return &BuiltinFunction{arity: 1, call: func(_ *Interpreter, args []any) any {
			for _, el := range l.Elements {
				if el == args[0] {
					return true
				}
			}
			return false
		}}
	case "index_of":
		return &BuiltinFunction{arity: 1, call: func(_ *Interpreter, args []any) any {
			for i, el := range l.Elements {
				if el == args[0] {
					return float64(i)
				}
			}
			return nil
		}}
	case "reverse":
		return &BuiltinFunction{arity: 0, call: func(_ *Interpreter, _ []any) any {
			for i, j := 0, len(l.Elements)-1; i < j; i, j = i+1, j-1 {
				l.Elements[i], l.Elements[j] = l.Elements[j], l.Elements[i]
			}
			return nil
		}}
	case "join":
		return &BuiltinFunction{arity: 1, call: func(_ *Interpreter, args []any) any {
			sep, ok := args[0].(string)
			if !ok {
				return nil
			}
			strs := make([]string, len(l.Elements))
			for i, el := range l.Elements {
				strs[i] = fmt.Sprintf("%v", el)
			}
			return strings.Join(strs, sep)
		}}
	default:
		return nil
	}
}
