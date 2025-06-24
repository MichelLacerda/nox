package runtime

import (
	"fmt"
	"strings"

	"github.com/MichelLacerda/nox/internal/token"
)

type ListInstance struct {
	Elements []any
}

func NewListInstance(elements []any) *ListInstance {
	return &ListInstance{Elements: elements}
}

func (l *ListInstance) Get(name *token.Token) any {
	switch name.Lexeme {
	case "append":
		return &BuiltinFunction{ArityValue: 1, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 1 {
				interpreter.Runtime.ReportRuntimeError(nil, "append(arg) expects 1 argument.")
				return nil
			}
			if l == nil {
				interpreter.Runtime.ReportRuntimeError(nil, "ListInstance is nil.")
				return nil
			}
			l.Elements = append(l.Elements, args[0])
			return nil
		}}
	case "pop":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 0 {
				interpreter.Runtime.ReportRuntimeError(nil, "pop() expects 0 arguments.")
				return nil
			}
			if l == nil {
				interpreter.Runtime.ReportRuntimeError(nil, "ListInstance is nil.")
				return nil
			}
			if len(l.Elements) == 0 {
				interpreter.Runtime.ReportRuntimeError(nil, "pop() called on empty list.")
				return nil
			}
			val := l.Elements[len(l.Elements)-1]
			l.Elements = l.Elements[:len(l.Elements)-1]
			return val
		}}
	case "insert":
		return &BuiltinFunction{ArityValue: 2, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 2 {
				interpreter.Runtime.ReportRuntimeError(nil, "insert(index, value) expects 2 arguments.")
				return nil
			}
			if l == nil {
				interpreter.Runtime.ReportRuntimeError(nil, "ListInstance is nil.")
				return nil
			}
			index, ok1 := args[0].(float64)
			if !ok1 {
				interpreter.Runtime.ReportRuntimeError(nil, "insert(index, value) expects index as number.")
				return nil
			}
			i := int(index)
			if i < 0 || i > len(l.Elements) {
				interpreter.Runtime.ReportRuntimeError(nil, "insert(index, value) index out of range.")
				return nil
			}
			l.Elements = append(l.Elements[:i], append([]any{args[1]}, l.Elements[i:]...)...)
			return nil
		}}
	case "remove":
		return &BuiltinFunction{ArityValue: 1, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 1 {
				interpreter.Runtime.ReportRuntimeError(nil, "remove(index) expects 1 argument.")
				return nil
			}
			if l == nil {
				interpreter.Runtime.ReportRuntimeError(nil, "ListInstance is nil.")
				return nil
			}
			index, ok := args[0].(float64)
			if !ok {
				interpreter.Runtime.ReportRuntimeError(nil, "remove(index) expects index as number.")
				return nil
			}
			i := int(index)
			if i < 0 || i >= len(l.Elements) {
				interpreter.Runtime.ReportRuntimeError(nil, "remove(index) index out of range.")
				return nil
			}
			l.Elements = append(l.Elements[:i], l.Elements[i+1:]...)
			return nil
		}}
	case "clear":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 0 {
				interpreter.Runtime.ReportRuntimeError(nil, "clear() expects 0 arguments.")
				return nil
			}
			if l == nil {
				interpreter.Runtime.ReportRuntimeError(nil, "ListInstance is nil.")
				return nil
			}
			l.Elements = []any{}
			return nil
		}}
	case "length":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 0 {
				interpreter.Runtime.ReportRuntimeError(nil, "length() expects 0 arguments.")
				return nil
			}
			if l == nil {
				interpreter.Runtime.ReportRuntimeError(nil, "ListInstance is nil.")
				return nil
			}
			return float64(len(l.Elements))
		}}
	case "contains":
		return &BuiltinFunction{ArityValue: 1, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 1 {
				interpreter.Runtime.ReportRuntimeError(nil, "contains(value) expects 1 argument.")
				return nil
			}
			if l == nil {
				interpreter.Runtime.ReportRuntimeError(nil, "ListInstance is nil.")
				return nil
			}
			for _, el := range l.Elements {
				if el == args[0] {
					return true
				}
			}
			return false
		}}
	case "index_of":
		return &BuiltinFunction{ArityValue: 1, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 1 {
				interpreter.Runtime.ReportRuntimeError(nil, "index_of(value) expects 1 argument.")
				return nil
			}
			if l == nil {
				interpreter.Runtime.ReportRuntimeError(nil, "ListInstance is nil.")
				return nil
			}
			for i, el := range l.Elements {
				if el == args[0] {
					return float64(i)
				}
			}
			return nil
		}}
	case "reverse":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 0 {
				interpreter.Runtime.ReportRuntimeError(nil, "reverse() expects 0 arguments.")
				return nil
			}
			if l == nil {
				interpreter.Runtime.ReportRuntimeError(nil, "ListInstance is nil.")
				return nil
			}
			for i, j := 0, len(l.Elements)-1; i < j; i, j = i+1, j-1 {
				l.Elements[i], l.Elements[j] = l.Elements[j], l.Elements[i]
			}
			return nil
		}}
	case "join":
		return &BuiltinFunction{ArityValue: 1, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 1 {
				interpreter.Runtime.ReportRuntimeError(nil, "join(sep) expects 1 argument.")
				return nil
			}

			if l == nil {
				interpreter.Runtime.ReportRuntimeError(nil, "ListInstance is nil.")
				return nil
			}

			sep, ok := args[0].(string)
			if !ok {
				interpreter.Runtime.ReportRuntimeError(nil, "join(sep) expects a string as the separator.")
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
