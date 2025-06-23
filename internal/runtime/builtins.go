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
	return "<builtin fn>"
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
				n, ok := args[0].(float64)
				if !ok {
					i.Runtime.ReportRuntimeError(token.NewToken(token.TokenType_IDENTIFIER, "math.abs", nil, 0), "Argument must be a number.")
					return nil
				}
				return math.Abs(n)
			},
		},
		"sqrt": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				n, ok := args[0].(float64)
				if !ok || n < 0 {
					i.Runtime.ReportRuntimeError(token.NewToken(token.TokenType_IDENTIFIER, "math.sqrt", nil, 0), "Argument must be a non-negative number.")
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
				n, ok := args[0].(float64)
				if !ok || n <= 0 {
					i.Runtime.ReportRuntimeError(token.NewToken(token.TokenType_IDENTIFIER, "math.log", nil, 0), "Argument must be a positive number.")
					return nil
				}
				return math.Log(n)
			},
		},
		"pow": &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(i *Interpreter, args []any) any {
				a, ok1 := args[0].(float64)
				b, ok2 := args[1].(float64)
				if !ok1 || !ok2 {
					i.Runtime.ReportRuntimeError(token.NewToken(token.TokenType_IDENTIFIER, "math.pow", nil, 0), "Arguments must be numbers.")
					return nil
				}
				return math.Pow(a, b)
			},
		},
		"max": &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(i *Interpreter, args []any) any {
				a, ok1 := args[0].(float64)
				b, ok2 := args[1].(float64)
				if !ok1 || !ok2 {
					i.Runtime.ReportRuntimeError(token.NewToken(token.TokenType_IDENTIFIER, "math.max", nil, 0), "Arguments must be numbers.")
					return nil
				}
				return math.Max(a, b)
			},
		},
		"min": &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(i *Interpreter, args []any) any {
				a, ok1 := args[0].(float64)
				b, ok2 := args[1].(float64)
				if !ok1 || !ok2 {
					i.Runtime.ReportRuntimeError(token.NewToken(token.TokenType_IDENTIFIER, "math.min", nil, 0), "Arguments must be numbers.")
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
			n, ok := args[0].(float64)
			if !ok {
				i.Runtime.ReportRuntimeError(token.NewToken(token.TokenType_IDENTIFIER, name, nil, 0), "Argument must be a number.")
				return nil
			}
			return fn(n)
		},
	}
}
