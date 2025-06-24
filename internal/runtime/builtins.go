package runtime

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/MichelLacerda/nox/internal/token"
	"github.com/MichelLacerda/nox/internal/util"
)

type BuiltinFunction struct {
	ArityValue int
	CallFunc   func(interpreter *Interpreter, args []any) any
}

func (b *BuiltinFunction) Arity() int {
	return b.ArityValue
}

func (b *BuiltinFunction) Call(interpreter *Interpreter, args []any) any {
	return b.CallFunc(interpreter, args)
}

func (b *BuiltinFunction) String() string {
	return "<builtin fn(" + fmt.Sprint(b.ArityValue) + ")>"
}

func RegisterClockBuiltin(i *Interpreter) *BuiltinFunction {
	return &BuiltinFunction{
		ArityValue: 0,
		CallFunc: func(i *Interpreter, args []any) any {
			if len(args) != 0 {
				i.Runtime.ReportRuntimeError(nil, "clock() expects no arguments.")
				return nil
			}
			return float64(time.Now().UnixNano()) / 1e9 // Retorna o tempo em segundos
		},
	}
}

func RegisterLenBuiltin(i *Interpreter) *BuiltinFunction {
	return &BuiltinFunction{
		ArityValue: 1,
		CallFunc: func(i *Interpreter, args []any) any {
			if len(args) != 1 {
				i.Runtime.ReportRuntimeError(nil, "len() expects 1 argument.")
				return nil
			}
			arg := args[0]
			switch v := arg.(type) {
			case string:
				return float64(utf8.RuneCountInString(v))
			case []any:
				return float64(len(v))
			case map[string]any:
				return float64(len(v))
			case *ListInstance:
				return float64(len(v.Elements))
			case *DictInstance:
				return float64(len(v.Entries))
			case *FileObject:
				if v.File == nil {
					i.Runtime.ReportRuntimeError(nil, "File is not open.")
					return nil
				}
				info, err := v.File.Stat()
				if err != nil {
					i.Runtime.ReportRuntimeError(nil, "Failed to get file info: "+err.Error())
					return nil
				}
				return float64(info.Size()) // Retorna o tamanho do arquivo em bytes
			default:
				i.Runtime.ReportRuntimeError(nil, fmt.Sprintf("len() expects a string, list, dict, or file, but got %T.", arg))
				return nil
			}
		},
	}
}

func RegisterRangeBuiltin(i *Interpreter) *BuiltinFunction {
	return &BuiltinFunction{
		ArityValue: -1, // Aceita 1, 2 ou 3 argumentos
		CallFunc: func(i *Interpreter, args []any) any {
			if len(args) < 1 || len(args) > 3 {
				i.Runtime.ReportRuntimeError(nil, "range() expects 1 to 3 arguments.")
				return nil
			}
			var start, end, step float64
			switch len(args) {
			case 1:
				start = 0
				end = util.ToFloat(args[0])
				step = 1

			case 2:
				start = util.ToFloat(args[0])
				end = util.ToFloat(args[1])
				step = 1
			case 3:
				start = util.ToFloat(args[0])
				end = util.ToFloat(args[1])
				step = util.ToFloat(args[2])
			default:
				i.Runtime.ReportRuntimeError(nil, "range() expects 1 to 3 arguments.")
				return nil
			}
			if step == 0 {
				i.Runtime.ReportRuntimeError(nil, "range() step must not be zero.")
				return nil
			}
			var result []any
			if step > 0 {
				for i := start; i < end; i += step {
					result = append(result, i)
				}
			} else {
				for i := start; i > end; i += step {
					result = append(result, i)
				}
			}
			return NewListInstance(result)
		},
	}
}

func RegisterAssertBuiltin(i *Interpreter) *BuiltinFunction {
	return &BuiltinFunction{
		ArityValue: 2,
		CallFunc: func(i *Interpreter, args []any) any {
			if len(args) != 2 {
				i.Runtime.ReportRuntimeError(nil, "assert(condition, message) expects 2 arguments.")
				return nil
			}

			condition := i.isTruthy(args[0])
			message := args[1]

			if condition {
				return nil // tudo certo
			}

			// modo debug → apenas imprime
			if i.debug {
				fmt.Printf("Assertion failed: %v\n", i.stringify(message))
				return nil
			}

			// modo normal → erro fatal
			i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "assert"}, fmt.Sprintf("Assertion failed: %v", message))
			return nil
		},
	}
}

func RegisterIoBuiltins(i *Interpreter) *BuiltinFunction {
	return &BuiltinFunction{
		ArityValue: 2,
		CallFunc: func(i *Interpreter, args []any) any {
			if len(args) != 2 {
				i.Runtime.ReportRuntimeError(nil, "open() expects 2 arguments.")
				return nil
			}

			path, ok1 := args[0].(string)
			mode, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "open"}, "open(path, mode) expects strings")
				return nil
			}

			flags, err := util.ParseFileMode(mode)
			if err != nil {
				i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "open"}, err.Error())
				return nil
			}

			f, err := os.OpenFile(path, flags, 0666)
			if err != nil {
				i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "open"}, "failed to open file: "+err.Error())
				return nil
			}
			return &FileObject{
				File:   f,
				Reader: bufio.NewReader(f),
			}
		},
	}
}

func RegisterFmtBuiltin(i *Interpreter) *BuiltinFunction {
	return &BuiltinFunction{
		ArityValue: -1,
		CallFunc: func(inter *Interpreter, args []any) any {
			if len(args) == 0 {
				return ""
			}
			format, ok := args[0].(string)
			if !ok {
				// Se o primeiro argumento não for string, apenas concatena todos
				var sb strings.Builder
				for i, arg := range args {
					if i > 0 {
						sb.WriteString(" ")
					}
					sb.WriteString(inter.stringify(arg))
				}
				return sb.String()
			}
			// Substitui cada '{}' pelo argumento correspondente
			result := ""
			parts := strings.Split(format, "{}")
			for i, part := range parts {
				result += part
				if i+1 < len(parts) && i+1 < len(args) {
					result += inter.stringify(args[i+1])
				}
			}
			// Se houver mais argumentos do que '{}', adiciona-os ao final
			if len(args) > len(parts) {
				for j := len(parts); j < len(args); j++ {
					result += " " + inter.stringify(args[j])
				}
			}
			return result
		},
	}
}

func RegisterMathBuiltin(i *Interpreter) *MapInstance {
	return NewMapInstance(map[string]any{
		"abs": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "math.abs(value) expects 1 argument.")
					return nil
				}

				n, ok := args[0].(float64)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "math.abs(value) expects a number argument.")
					return nil
				}
				return math.Abs(n)
			},
		},
		"sqrt": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "math.sqrt(value) expects 1 argument.")
					return nil
				}

				n, ok := args[0].(float64)
				if !ok || n < 0 {
					i.Runtime.ReportRuntimeError(nil, "Argument must be a non-negative number.")
					return nil
				}
				return math.Sqrt(n)
			},
		},
		"sin":   mathUnary("math.sin", math.Sin),
		"cos":   mathUnary("math.cos", math.Cos),
		"tan":   mathUnary("math.tan", math.Tan),
		"floor": mathUnary("math.floor", math.Floor),
		"ceil":  mathUnary("math.ceil", math.Ceil),
		"round": mathUnary("math.round", math.Round),
		"exp":   mathUnary("math.exp", math.Exp),
		"log": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "math.log(value) expects 1 argument.")
					return nil
				}
				n, ok := args[0].(float64)
				if !ok || n <= 0 {
					i.Runtime.ReportRuntimeError(nil, "Argument must be a positive number.")
					return nil
				}
				return math.Log(n)
			},
		},
		"pow": &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 2 {
					i.Runtime.ReportRuntimeError(nil, "math.pow(base, exponent) expects 2 arguments.")
					return nil
				}
				a, ok1 := args[0].(float64)
				b, ok2 := args[1].(float64)
				if !ok1 || !ok2 {
					i.Runtime.ReportRuntimeError(nil, "Arguments must be numbers.")
					return nil
				}
				return math.Pow(a, b)
			},
		},
		"max": &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 2 {
					i.Runtime.ReportRuntimeError(nil, "math.max(a, b) expects 2 arguments.")
					return nil
				}
				a, ok1 := args[0].(float64)
				b, ok2 := args[1].(float64)
				if !ok1 || !ok2 {
					i.Runtime.ReportRuntimeError(nil, "Arguments must be numbers.")
					return nil
				}
				return math.Max(a, b)
			},
		},
		"min": &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 2 {
					i.Runtime.ReportRuntimeError(nil, "math.min(a, b) expects 2 arguments.")
					return nil
				}

				a, ok1 := args[0].(float64)
				b, ok2 := args[1].(float64)
				if !ok1 || !ok2 {
					i.Runtime.ReportRuntimeError(nil, "Arguments must be numbers.")
					return nil
				}
				return math.Min(a, b)
			},
		},
	})
}

func mathUnary(name string, fn func(float64) float64) *BuiltinFunction {
	return &BuiltinFunction{
		ArityValue: 1,
		CallFunc: func(i *Interpreter, args []any) any {
			if len(args) != 1 {
				i.Runtime.ReportRuntimeError(nil, name+"() expects 1 argument.")
				return nil
			}
			n, ok := args[0].(float64)
			if !ok {
				i.Runtime.ReportRuntimeError(nil, "Argument must be a number.")
				return nil
			}
			return fn(n)
		},
	}
}

func RegisterTypeBuiltins(i *Interpreter) *MapInstance {
	return NewMapInstance(map[string]any{
		"of": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "type.of() expects 1 argument.")
					return nil
				}
				return TypeOf(args[0])
			},
		},
		"is": &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 2 {
					i.Runtime.ReportRuntimeError(nil, "type.is() expects 2 arguments.")
					return nil
				}
				value := args[0]
				expectedType, ok := args[1].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "type.is() expects the second argument to be a string representing the type.")
					return nil
				}
				actualType := TypeOf(value)
				return actualType == expectedType
			},
		},
		"is_nil": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "type.is_nil() expects 1 argument.")
					return nil
				}
				value := args[0]
				if value == nil {
					return true
				}
				return false
			},
		},
		"is_bool": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "type.is_bool() expects 1 argument.")
					return nil
				}
				value := args[0]
				return TypeOf(value) == "bool"
			},
		},
		"is_number": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "type.is_number() expects 1 argument.")
					return nil
				}
				value := args[0]
				return TypeOf(value) == "number"
			},
		},
		"is_string": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "type.is_string() expects 1 argument.")
					return nil
				}
				value := args[0]
				return TypeOf(value) == "string"
			},
		},
		"is_list": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "type.is_list() expects 1 argument.")
					return nil
				}
				value := args[0]
				return TypeOf(value) == "list"
			},
		},
		"is_dict": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "type.is_dict() expects 1 argument.")
					return nil
				}
				value := args[0]
				return TypeOf(value) == "dict"
			},
		},
		"is_function": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "type.is_function() expects 1 argument.")
					return nil
				}
				value := args[0]
				return TypeOf(value) == "function"
			},
		},
		"is_class": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "type.is_class() expects 1 argument.")
					return nil
				}
				value := args[0]
				return TypeOf(value) == "class"
			},
		},
		"is_instance": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "type.is_instance() expects 1 argument.")
					return nil
				}
				value := args[0]
				return TypeOf(value) == "instance"
			},
		},
		"is_iterable": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "type.is_iterable() expects 1 argument.")
					return nil
				}
				value := args[0]
				t := TypeOf(value)
				return t == "list" || t == "dict" || t == "string"
			},
		},
		"is_callable": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "type.is_callable() expects 1 argument.")
					return nil
				}
				value := args[0]
				t := TypeOf(value)
				return t == "function" || t == "class"
			},
		},
		"is_truthy": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "type.is_truthy() expects 1 argument.")
					return nil
				}
				value := args[0]
				if value == nil {
					return false
				}
				if b, ok := value.(bool); ok {
					return b
				}
				if n, ok := value.(float64); ok {
					return n != 0
				}
				if s, ok := value.(string); ok {
					return s != ""
				}
				if l, ok := value.([]any); ok {
					return len(l) > 0
				}
				if m, ok := value.(map[string]any); ok {
					return len(m) > 0
				}
				return true
			},
		},
		"is_falsey": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "type.is_falsey() expects 1 argument.")
					return nil
				}
				value := args[0]
				if value == nil {
					return true
				}
				if b, ok := value.(bool); ok {
					return !b
				}
				if n, ok := value.(float64); ok {
					return n == 0
				}
				if s, ok := value.(string); ok {
					return s == ""
				}
				if l, ok := value.([]any); ok {
					return len(l) == 0
				}
				if m, ok := value.(map[string]any); ok {
					return len(m) == 0
				}
				return false
			},
		},
		"instance_of": &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 2 {
					i.Runtime.ReportRuntimeError(nil, "type.instance(instance, class) expects 2 arguments.")
					return nil
				}
				instance, ok1 := args[0].(*Instance)
				class, ok2 := args[1].(*Class)
				if !ok1 || !ok2 {
					i.Runtime.ReportRuntimeError(nil, "type.instance(instance, class) expects an instance and a class.")
					return nil
				}
				return instance.IsInstanceOf(class)
			},
		},
	})
}

func RegisterMathConstants(i *Interpreter) {
	i.globals.Define("PI", 3.141592)
	i.globals.Define("E", 2.718281)
	i.globals.Define("PHI", 1.618033)
	i.globals.Define("TAU", 6.283185) // tau = 2 * PI
	i.globals.Define("sqrt2", 1.414213)
	i.globals.Define("sqrtE", 1.648721)
	i.globals.Define("sqrtPi", 1.772453)
	i.globals.Define("sqrtPhi", 1.272019)
	i.globals.Define("ln2", 0.693147)
	i.globals.Define("log2E", 1/0.693147)
	i.globals.Define("ln10", 2.302585)
	i.globals.Define("log10E", 1/2.302585)
}

func RegisterBuiltins(i *Interpreter) {
	i.globals.Define("clock", RegisterClockBuiltin(i))
	i.globals.Define("len", RegisterLenBuiltin(i))
	i.globals.Define("range", RegisterRangeBuiltin(i))
	i.globals.Define("assert", RegisterAssertBuiltin(i))
	i.globals.Define("open", RegisterIoBuiltins(i))
	i.globals.Define("math", RegisterMathBuiltin(i))
	i.globals.Define("fmt", RegisterFmtBuiltin(i))
	i.globals.Define("type", RegisterTypeBuiltins(i))
	RegisterMathConstants(i)
}

func TypeOf(v any) any {
	switch v.(type) {
	case nil:
		return "nil"
	case bool:
		return "bool"
	case float64:
		return "number"
	case string:
		return "string"
	case []any:
		return "list"
	case map[string]any:
		return "dict"
	case *BuiltinFunction, *Function:
		return "function"
	case *Class:
		return "class"
	case *Instance:
		return "instance"
	case *DictInstance:
		return "dict"
	case *ListInstance:
		return "list"
	default:
		return "unknown"
	}
}
