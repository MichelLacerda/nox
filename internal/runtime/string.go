package runtime

import (
	"strings"
	"unicode/utf8"
)

type StringInstance struct {
	Value string
}

var _ HasMethods = (*StringInstance)(nil)

func (s *StringInstance) GetMethod(name string) any {
	switch name {
	case "length":
		return &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(i *Interpreter, _ []any) any {
				return float64(utf8.RuneCountInString(s.Value))
			},
		}
	case "upper":
		return &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(i *Interpreter, _ []any) any {
				return strings.ToUpper(s.Value)
			},
		}
	case "lower":
		return &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(i *Interpreter, _ []any) any {
				return strings.ToLower(s.Value)
			},
		}
	case "split":
		return &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "String.split expects 1 argument.")
					return nil
				}
				sep, ok := args[0].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "String.split expects a string as argument.")
					return nil
				}
				parts := strings.Split(s.Value, sep)
				result := make([]any, len(parts))
				for i, part := range parts {
					result[i] = part
				}
				return ListInstance{
					Elements: result,
				}
			},
		}
	case "replace":
		return &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 2 {
					i.Runtime.ReportRuntimeError(nil, "String.replace expects 2 arguments.")
					return nil
				}
				old, ok1 := args[0].(string)
				new, ok2 := args[1].(string)
				if !ok1 || !ok2 {
					i.Runtime.ReportRuntimeError(nil, "String.replace expects two strings as arguments.")
					return nil
				}
				return strings.ReplaceAll(s.Value, old, new)
			},
		}
	case "contains":
		return &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "String.contains expects 1 argument.")
					return nil
				}
				substr, ok := args[0].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "String.contains expects a string as argument.")
					return nil
				}
				return strings.Contains(s.Value, substr)
			},
		}
	case "index_of":
		return &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "String.indexOf expects 1 argument.")
					return nil
				}
				substr, ok := args[0].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "String.indexOf expects a string as argument.")
					return nil
				}
				index := strings.Index(s.Value, substr)
				if index == -1 {
					return nil // Retorna nil se não encontrar
				}
				return float64(index)
			},
		}
	case "last_index_of":
		return &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "String.lastIndexOf expects 1 argument.")
					return nil
				}
				substr, ok := args[0].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "String.lastIndexOf expects a string as argument.")
					return nil
				}
				index := strings.LastIndex(s.Value, substr)
				if index == -1 {
					return nil // Retorna nil se não encontrar
				}
				return float64(index)
			},
		}
	case "trim":
		return &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(i *Interpreter, _ []any) any {
				return strings.TrimSpace(s.Value)
			},
		}
	default:
		return nil
	}
}
