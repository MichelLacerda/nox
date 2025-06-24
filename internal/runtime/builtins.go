package runtime

import (
	"bufio"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"os/exec"
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

func RegisterRandomBuiltins(i *Interpreter) *MapInstance {
	return NewMapInstance(map[string]any{
		"float": &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 0 {
					i.Runtime.ReportRuntimeError(nil, "math.random() expects no arguments.")
					return nil
				}
				return rand.Float64()
			},
		},
		"int": &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 2 {
					i.Runtime.ReportRuntimeError(nil, "math.irand(min, max) expects 2 arguments.")
					return nil
				}
				min, ok1 := args[0].(float64)
				max, ok2 := args[1].(float64)
				if !ok1 || !ok2 {
					i.Runtime.ReportRuntimeError(nil, "math.irand(min, max) expects two numbers.")
					return nil
				}
				if min >= max {
					i.Runtime.ReportRuntimeError(nil, "math.irand(min, max) expects min < max.")
					return nil
				}
				return float64(rand.IntN(int(max-min))) + min
			},
		},
	})
}

func RegisterOsBuiltins(i *Interpreter) *MapInstance {
	return NewMapInstance(map[string]any{
		"exit": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "os.exit(code) expects 1 argument.")
					return nil
				}
				code, ok := args[0].(float64)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "os.exit(code) expects a number argument.")
					return nil
				}
				if code < 0 || code > 255 {
					i.Runtime.ReportRuntimeError(nil, "os.exit(code) expects a code between 0 and 255.")
					return nil
				}
				os.Exit(int(code))
				return nil // nunca alcançado, mas necessário para satisfazer a assinatura
			},
		},
		"exec": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "os.exec(command) expects 1 argument.")
					return nil
				}
				command, ok := args[0].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "os.exec(command) expects a string argument.")
					return nil
				}
				var cmd *exec.Cmd
				if os.PathSeparator == '\\' {
					// Provavelmente Windows
					cmd = exec.Command("cmd", "/C", command)
				} else {
					// Unix-like
					cmd = exec.Command("sh", "-c", command)
				}
				output, err := cmd.CombinedOutput()
				if err != nil {
					i.Runtime.ReportRuntimeError(nil, fmt.Sprintf("os.exec failed: %v", err))
					return nil
				}
				return string(output)
			},
		},
		"getenv": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "os.getenv(name) expects 1 argument.")
					return nil
				}
				name, ok := args[0].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "os.getenv(name) expects a string argument.")
					return nil
				}
				value, exists := os.LookupEnv(name)
				if !exists {
					i.Runtime.ReportRuntimeError(nil, fmt.Sprintf("Environment variable '%s' not found.", name))
					return nil
				}
				return value
			},
		},
		"setenv": &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 2 {
					i.Runtime.ReportRuntimeError(nil, "os.setenv(name, value) expects 2 arguments.")
					return nil
				}
				name, ok1 := args[0].(string)
				value, ok2 := args[1].(string)
				if !ok1 || !ok2 {
					i.Runtime.ReportRuntimeError(nil, "os.setenv(name, value) expects both arguments to be strings.")
					return nil
				}
				err := os.Setenv(name, value)
				if err != nil {
					i.Runtime.ReportRuntimeError(nil, fmt.Sprintf("Failed to set environment variable '%s': %v", name, err))
					return nil
				}
				return nil
			},
		},
		"cwd": &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 0 {
					i.Runtime.ReportRuntimeError(nil, "os.cwd() expects no arguments.")
					return nil
				}
				cwd, err := os.Getwd()
				if err != nil {
					i.Runtime.ReportRuntimeError(nil, fmt.Sprintf("os.cwd() failed: %v", err))
					return nil
				}
				return cwd
			},
		},
		"listdir": &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 2 {
					i.Runtime.ReportRuntimeError(nil, "os.listdir(path, only_dirs) expects 2 arguments.")
					return nil
				}
				path, ok := args[0].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "os.listdir(path) expects a string argument.")
					return nil
				}
				onlyDirs, ok := args[1].(bool)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "os.listdir(path, only_dirs) expects a boolean argument.")
					return nil
				}
				entries, err := os.ReadDir(path)
				if err != nil {
					i.Runtime.ReportRuntimeError(nil, fmt.Sprintf("os.listdir failed: %v", err))
					return nil
				}
				var names []any
				for _, entry := range entries {
					if onlyDirs {
						if entry.IsDir() {
							names = append(names, entry.Name()+"/") // adiciona barra para diretórios
						}
					} else {
						names = append(names, entry.Name())
					}
				}
				return NewListInstance(names)
			},
		},
		"chmod": &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 2 {
					i.Runtime.ReportRuntimeError(nil, "os.chmod(path, mode) expects 2 arguments.")
					return nil
				}
				path, ok1 := args[0].(string)
				mode, ok2 := args[1].(float64)
				if !ok1 || !ok2 {
					i.Runtime.ReportRuntimeError(nil, "os.chmod(path, mode) expects a string and a number.")
					return nil
				}
				err := os.Chmod(path, os.FileMode(mode))
				if err != nil {
					i.Runtime.ReportRuntimeError(nil, fmt.Sprintf("os.chmod failed: %v", err))
					return nil
				}
				return nil
			},
		},
		"chdir": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "os.chdir(path) expects 1 argument.")
					return nil
				}
				path, ok := args[0].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "os.chdir(path) expects a string argument.")
					return nil
				}
				err := os.Chdir(path)
				if err != nil {
					i.Runtime.ReportRuntimeError(nil, fmt.Sprintf("os.chdir failed: %v", err))
					return nil
				}
				return nil
			},
		},
		"mkdir": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "os.mkdir(path) expects 1 argument.")
					return nil
				}
				path, ok := args[0].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "os.mkdir(path) expects a string argument.")
					return nil
				}
				err := os.Mkdir(path, 0755) // Permissões padrão 0755
				if err != nil {
					i.Runtime.ReportRuntimeError(nil, fmt.Sprintf("os.mkdir failed: %v", err))
					return nil
				}
				return nil
			},
		},
		"rmdir": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "os.rmdir(path) expects 1 argument.")
					return nil
				}
				path, ok := args[0].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "os.rmdir(path) expects a string argument.")
					return nil
				}
				err := os.Remove(path)
				if err != nil {
					i.Runtime.ReportRuntimeError(nil, fmt.Sprintf("os.rmdir failed: %v", err))
					return nil
				}
				return nil
			},
		},
		"walk": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "os.walk(path) expects 1 argument.")
					return nil
				}
				root, ok := args[0].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "os.walk(path) expects a string argument.")
					return nil
				}

				var walkDir func(string) []any
				walkDir = func(path string) []any {
					entries, err := os.ReadDir(path)
					if err != nil {
						// Retorna erro como string para não interromper toda a recursão
						return []any{fmt.Sprintf("os.walk failed at %s: %v", path, err)}
					}
					var files []any
					for _, entry := range entries {
						fullPath := path + string(os.PathSeparator) + entry.Name()
						if entry.IsDir() {
							files = append(files, NewDictInstance(map[string]any{
								"name":     entry.Name() + "/",
								"type":     "directory",
								"children": NewListInstance(walkDir(fullPath)),
							}))
						} else {
							files = append(files, NewDictInstance(map[string]any{
								"name": entry.Name(),
								"type": "file",
							}))
						}
					}
					return files
				}

				return NewListInstance(walkDir(root))
			},
		},
		"info": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				if len(args) != 1 {
					i.Runtime.ReportRuntimeError(nil, "os.path(path) expects 1 argument.")
					return nil
				}
				path, ok := args[0].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "os.path(path) expects a string argument.")
					return nil
				}
				// Retorna um dicionário com informações sobre o caminho
				info, err := os.Stat(path)
				if err != nil {
					i.Runtime.ReportRuntimeError(nil, fmt.Sprintf("os.path failed: %v", err))
					return nil
				}
				return NewDictInstance(map[string]any{
					"name":        info.Name(),
					"size":        float64(info.Size()),
					"mode":        float64(info.Mode()),
					"mod_time":    info.ModTime().Format(time.RFC3339),
					"is_dir":      info.IsDir(),
					"is_file":     !info.IsDir(),
					"is_symlink":  info.Mode()&os.ModeSymlink != 0,
					"permissions": fmt.Sprintf("%04o", info.Mode().Perm()),
				})
			},
		},
	})
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
	i.globals.Define("random", RegisterRandomBuiltins(i))
	i.globals.Define("os", RegisterOsBuiltins(i))
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
