package runtime

import (
	"strconv"
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
			CallFunc: func(interpreter *Interpreter, args []any) any {
				if len(args) != 0 {
					interpreter.Runtime.ReportRuntimeError(nil, "String.length expects 0 arguments.")
					return nil
				}
				return float64(utf8.RuneCountInString(s.Value))
			},
		}
	case "upper":
		return &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(interpreter *Interpreter, args []any) any {
				if len(args) != 0 {
					interpreter.Runtime.ReportRuntimeError(nil, "String.upper expects 0 arguments.")
					return nil
				}
				return strings.ToUpper(s.Value)
			},
		}
	case "lower":
		return &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(interpreter *Interpreter, args []any) any {
				if len(args) != 0 {
					interpreter.Runtime.ReportRuntimeError(nil, "String.lower expects 0 arguments.")
					return nil
				}
				return strings.ToLower(s.Value)
			},
		}
	case "split":
		return &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(interpreter *Interpreter, args []any) any {
				if len(args) != 1 {
					interpreter.Runtime.ReportRuntimeError(nil, "String.split expects 1 argument.")
					return nil
				}
				sep, ok := args[0].(string)
				if !ok {
					interpreter.Runtime.ReportRuntimeError(nil, "String.split expects a string as argument.")
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
			CallFunc: func(interpreter *Interpreter, args []any) any {
				if len(args) != 2 {
					interpreter.Runtime.ReportRuntimeError(nil, "String.replace expects 2 arguments.")
					return nil
				}
				old, ok1 := args[0].(string)
				if !ok1 {
					interpreter.Runtime.ReportRuntimeError(nil, "String.replace expects the first argument to be a string.")
					return nil
				}
				newStr, ok2 := args[1].(string)
				if !ok2 {
					interpreter.Runtime.ReportRuntimeError(nil, "String.replace expects the second argument to be a string.")
					return nil
				}
				return strings.ReplaceAll(s.Value, old, newStr)
			},
		}
	case "contains":
		return &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(interpreter *Interpreter, args []any) any {
				if len(args) != 1 {
					interpreter.Runtime.ReportRuntimeError(nil, "String.contains expects 1 argument.")
					return nil
				}
				substr, ok := args[0].(string)
				if !ok {
					interpreter.Runtime.ReportRuntimeError(nil, "String.contains expects a string as argument.")
					return nil
				}
				return strings.Contains(s.Value, substr)
			},
		}
	case "index_of":
		return &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(interpreter *Interpreter, args []any) any {
				if len(args) != 1 {
					interpreter.Runtime.ReportRuntimeError(nil, "String.index_of expects 1 argument.")
					return nil
				}
				substr, ok := args[0].(string)
				if !ok {
					interpreter.Runtime.ReportRuntimeError(nil, "String.index_of expects a string as argument.")
					return nil
				}
				index := strings.Index(s.Value, substr)
				if index == -1 {
					return nil
				}
				return float64(index)
			},
		}
	case "last_index_of":
		return &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(interpreter *Interpreter, args []any) any {
				if len(args) != 1 {
					interpreter.Runtime.ReportRuntimeError(nil, "String.last_index_of expects 1 argument.")
					return nil
				}
				substr, ok := args[0].(string)
				if !ok {
					interpreter.Runtime.ReportRuntimeError(nil, "String.last_index_of expects a string as argument.")
					return nil
				}
				index := strings.LastIndex(s.Value, substr)
				if index == -1 {
					return nil
				}
				return float64(index)
			},
		}
	case "trim":
		return &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(interpreter *Interpreter, args []any) any {
				if len(args) != 0 {
					interpreter.Runtime.ReportRuntimeError(nil, "String.trim expects 0 arguments.")
					return nil
				}
				return strings.TrimSpace(s.Value)
			},
		}
	case "to_number":
		return &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(interpreter *Interpreter, args []any) any {
				if len(args) != 0 {
					interpreter.Runtime.ReportRuntimeError(nil, "String.to_number expects 0 arguments.")
					return nil
				}
				num, err := strconv.ParseFloat(s.Value, 64)
				if err != nil {
					interpreter.Runtime.ReportRuntimeError(nil, "String.to_number: "+err.Error())
					return nil
				}
				return num
			},
		}
	default:
		return nil
	}
}
