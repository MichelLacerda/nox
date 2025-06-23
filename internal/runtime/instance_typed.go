package runtime

import "github.com/MichelLacerda/nox/internal/token"

type TypedInstance struct {
	Value any
}

func NewTypedInstance(value any) *TypedInstance {
	return &TypedInstance{Value: value}
}

func (c *TypedInstance) Get(name *token.Token) any {
	switch name.Lexeme {
	case "value":
		return c.Value
	case "of":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(_ *Interpreter, _ []any) any {
			return TypeOf(c.Value)
		}}
	case "is_null":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(_ *Interpreter, _ []any) any {
			return c.Value == nil
		}}
	case "is_bool":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(_ *Interpreter, _ []any) any {
			return TypeOf(c.Value) == "bool"
		}}
	case "is_number":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(_ *Interpreter, _ []any) any {
			return TypeOf(c.Value) == "number"
		}}
	case "is_string":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(_ *Interpreter, _ []any) any {
			return TypeOf(c.Value) == "string"
		}}
	case "is_list":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(_ *Interpreter, _ []any) any {
			return TypeOf(c.Value) == "list"
		}}
	case "is_map":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(_ *Interpreter, _ []any) any {
			return TypeOf(c.Value) == "map"
		}}
	case "is_function":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(_ *Interpreter, _ []any) any {
			return TypeOf(c.Value) == "function"
		}}
	case "is_class":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(_ *Interpreter, _ []any) any {
			return TypeOf(c.Value) == "class"
		}}
	case "is_instance":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(_ *Interpreter, _ []any) any {
			return TypeOf(c.Value) == "instance"
		}}
	case "is_iterable":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(_ *Interpreter, _ []any) any {
			t := TypeOf(c.Value)
			return t == "list" || t == "map" || t == "string"
		}}
	case "is_callable":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(_ *Interpreter, _ []any) any {
			t := TypeOf(c.Value)
			return t == "function" || t == "class" || t == "instance"
		}}
	case "is_truthy":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(_ *Interpreter, _ []any) any {
			if c.Value == nil {
				return false
			}
			if b, ok := c.Value.(bool); ok {
				return b
			}
			if n, ok := c.Value.(float64); ok {
				return n != 0
			}
			if s, ok := c.Value.(string); ok {
				return s != ""
			}
			if l, ok := c.Value.([]any); ok {
				return len(l) > 0
			}
			if m, ok := c.Value.(map[string]any); ok {
				return len(m) > 0
			}
			return true
		}}
	case "is_falsey":
		return &BuiltinFunction{ArityValue: 0, CallFunc: func(_ *Interpreter, _ []any) any {
			if c.Value == nil {
				return true
			}
			if b, ok := c.Value.(bool); ok {
				return !b
			}
			if n, ok := c.Value.(float64); ok {
				return n == 0
			}
			if s, ok := c.Value.(string); ok {
				return s == ""
			}
			if l, ok := c.Value.([]any); ok {
				return len(l) == 0
			}
			if m, ok := c.Value.(map[string]any); ok {
				return len(m) == 0
			}
			return false
		}}
	default:
		return nil
	}
}
