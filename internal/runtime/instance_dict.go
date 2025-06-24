package runtime

import (
	"fmt"

	"github.com/MichelLacerda/nox/internal/token"
)

type DictInstance struct {
	Entries map[string]any
}

func NewDictInstance(entries map[string]any) *DictInstance {
	return &DictInstance{Entries: entries}
}

func (d *DictInstance) Get(name *token.Token) any {
	switch name.Lexeme {
	case "get":
		return &BuiltinFunction{ArityValue: 1, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 1 {
				interpreter.Runtime.ReportRuntimeError(name, fmt.Sprintf("dict.get: expected 1 argument, got %d", len(args)))
				return nil
			}
			key, ok := args[0].(string)
			if !ok {
				interpreter.Runtime.ReportRuntimeError(name, "dict.get: key must be a string")
				return nil
			}
			val, exists := d.Entries[key]
			if !exists {
				return nil
			}
			return val
		}}

	case "set":
		return &BuiltinFunction{ArityValue: 2, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 2 {
				interpreter.Runtime.ReportRuntimeError(name, fmt.Sprintf("dict.set: expected 2 arguments, got %d", len(args)))
				return nil
			}
			key, ok := args[0].(string)
			if !ok {
				interpreter.Runtime.ReportRuntimeError(name, "dict.set: key must be a string")
				return nil
			}
			d.Entries[key] = args[1]
			return nil
		}}

	case "remove":
		return &BuiltinFunction{ArityValue: 1, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 1 {
				interpreter.Runtime.ReportRuntimeError(name, fmt.Sprintf("dict.remove: expected 1 argument, got %d", len(args)))
				return nil
			}
			key, ok := args[0].(string)
			if !ok {
				interpreter.Runtime.ReportRuntimeError(name, "dict.remove: key must be a string")
				return nil
			}
			_, existed := d.Entries[key]
			delete(d.Entries, key)
			return existed
		}}

	case "keys":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 0 {
				interpreter.Runtime.ReportRuntimeError(name, fmt.Sprintf("dict.keys: expected 0 arguments, got %d", len(args)))
				return nil
			}
			keys := []any{}
			for k := range d.Entries {
				keys = append(keys, k)
			}
			return NewListInstance(keys)
		}}
	case "values":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 0 {
				interpreter.Runtime.ReportRuntimeError(name, fmt.Sprintf("dict.values: expected 0 arguments, got %d", len(args)))
				return nil
			}
			values := []any{}
			for _, v := range d.Entries {
				values = append(values, v)
			}
			return NewListInstance(values)
		}}
	case "clear":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 0 {
				interpreter.Runtime.ReportRuntimeError(name, fmt.Sprintf("dict.clear: expected 0 arguments, got %d", len(args)))
				return nil
			}
			d.Entries = map[string]any{}
			return nil
		}}

	case "contains":
		return &BuiltinFunction{ArityValue: 1, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 1 {
				interpreter.Runtime.ReportRuntimeError(name, fmt.Sprintf("dict.contains: expected 1 argument, got %d", len(args)))
				return false
			}
			key, ok := args[0].(string)
			if !ok {
				interpreter.Runtime.ReportRuntimeError(name, "dict.contains: key must be a string")
				return false
			}
			_, exists := d.Entries[key]
			return exists
		}}

	case "length":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(interpreter *Interpreter, args []any) any {
			if len(args) != 0 {
				interpreter.Runtime.ReportRuntimeError(name, fmt.Sprintf("dict.length: expected 0 arguments, got %d", len(args)))
				return float64(0)
			}
			return float64(len(d.Entries))
		}}

	default:
		return nil
	}
}
