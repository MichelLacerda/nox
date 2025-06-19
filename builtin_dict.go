package main

import (
	"fmt"
)

type DictInstance struct {
	Entries map[string]any
}

func NewDictInstance(entries map[string]any) *DictInstance {
	return &DictInstance{Entries: entries}
}

func (d *DictInstance) Get(name *Token) any {
	switch name.Lexeme {
	case "get":
		return &BuiltinFunction{arity: 1, call: func(i *Interpreter, args []any) any {
			key, ok := args[0].(string)
			if !ok {
				i.runtime.ReportRuntimeError(name, "dict.get: key must be a string")
				return nil
			}
			val, exists := d.Entries[key]
			if !exists {
				return nil
			}
			return val
		}}

	case "set":
		return &BuiltinFunction{arity: 2, call: func(i *Interpreter, args []any) any {
			key, ok := args[0].(string)
			if !ok {
				i.runtime.ReportRuntimeError(name, "dict.set: key must be a string")
				return nil
			}
			d.Entries[key] = args[1]
			return nil
		}}

	case "remove":
		return &BuiltinFunction{arity: 1, call: func(i *Interpreter, args []any) any {
			key, ok := args[0].(string)
			if !ok {
				i.runtime.ReportRuntimeError(name, "dict.remove: key must be a string")
				return nil
			}
			_, existed := d.Entries[key]
			delete(d.Entries, key)
			return existed
		}}

	case "keys":
		return &BuiltinFunction{arity: 0, call: func(i *Interpreter, _ []any) any {
			keys := []any{}
			for k := range d.Entries {
				keys = append(keys, k)
			}
			return NewListInstance(keys)
		}}
	case "values":
		return &BuiltinFunction{arity: 0, call: func(i *Interpreter, _ []any) any {
			values := []any{}
			for _, v := range d.Entries {
				values = append(values, v)
			}
			return NewListInstance(values)
		}}
	case "clear":
		return &BuiltinFunction{arity: 0, call: func(_ *Interpreter, _ []any) any {
			d.Entries = map[string]any{}
			return nil
		}}

	case "contains":
		return &BuiltinFunction{arity: 1, call: func(_ *Interpreter, args []any) any {
			key, ok := args[0].(string)
			if !ok {
				return false
			}
			_, exists := d.Entries[key]
			return exists
		}}

	case "length":
		return &BuiltinFunction{arity: 0, call: func(_ *Interpreter, _ []any) any {
			return float64(len(d.Entries))
		}}

	default:
		return nil
	}
}

func (d *DictInstance) String() string {
	return fmt.Sprintf("<dict %v>", d.Entries)
}
