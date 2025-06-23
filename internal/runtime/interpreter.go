package runtime

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/MichelLacerda/nox/internal/ast"
	"github.com/MichelLacerda/nox/internal/token"
)

type Interpreter struct {
	Runtime      *Nox
	globals      *Environment
	locals       map[ast.Expr]int
	environment  *Environment
	silentErrors bool
	debug        bool // Modo de depuração
	Colored      bool // Se deve usar cores na saída
}

type HasMethods interface {
	GetMethod(name string) any
}

func NewInterpreter(r *Nox, colored bool) *Interpreter {
	interpreter := &Interpreter{
		Runtime:      r,
		globals:      NewEnvironment(r, nil),
		environment:  nil, // Inicialmente nil
		locals:       map[ast.Expr]int{},
		silentErrors: true,    // Inicialmente não silencioso
		debug:        false,   // Modo de depuração desativado por padrão
		Colored:      colored, // Cores ativadas por padrão
	}
	interpreter.environment = interpreter.globals // Aponta para o global no início
	interpreter.globals.Define("clock", RegisterClockBuiltin(interpreter))
	interpreter.globals.Define("len", RegisterLenBuiltin(interpreter))
	interpreter.globals.Define("range", RegisterRangeBuiltin(interpreter))
	interpreter.globals.Define("assert", RegisterAssertBuiltin(interpreter))
	interpreter.globals.Define("open", RegisterIoBuiltins(interpreter))
	interpreter.globals.Define("math", RegisterMathBuiltin(interpreter))
	interpreter.globals.Define("fmt", RegisterFmtBuiltin(interpreter))

	return interpreter
}

func (i *Interpreter) Interpret(expr []ast.Stmt) {
	for _, statement := range expr {
		i.execute(statement)
	}
}

func (i *Interpreter) execute(s ast.Stmt) error {
	s.Accept(i)
	return nil
}

func (i *Interpreter) Resolve(expr ast.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) ExecuteBlock(statements []ast.Stmt, environment *Environment) error {
	previous := i.environment
	i.environment = environment
	defer func() {
		i.environment = previous
	}()

	for _, stmt := range statements {
		i.execute(stmt)
	}

	return nil
}

func (i *Interpreter) lookUpVariable(t *token.Token, expr ast.Expr) any {
	if depth, ok := i.locals[expr]; ok {
		return i.environment.GetAt(depth, t.Lexeme)
	}
	return i.globals.Get(t)
}

// ===== Helpers =====

func (i *Interpreter) evaluate(expr ast.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) isTruthy(value any) bool {
	switch v := value.(type) {
	case nil:
		return false
	case bool:
		return v
	default:
		return true
	}
}

func (i *Interpreter) isEqual(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

// ===== Operand validation =====

func (i *Interpreter) mustBeNumber(op *token.Token, val any) bool {
	if _, ok := val.(float64); !ok {
		i.Runtime.ReportRuntimeError(op, "Operand must be a number.")
		return false
	}
	return true
}

func (i *Interpreter) mustBeNumbers(op *token.Token, left, right any) bool {
	if _, ok := left.(float64); !ok {
		i.Runtime.ReportRuntimeError(op, "Left operand must be a number.")
		return false
	}
	if _, ok := right.(float64); !ok {
		i.Runtime.ReportRuntimeError(op, "Right operand must be a number.")
		return false
	}
	return true
}

// ===== Stringify helpers =====

func (i *Interpreter) stringify(value any) string {
	if i.Colored {
		return StringifyColor(value, "")
	}
	return StringifyCompact(value)
}

func StringifyCompact(value any) string {
	switch v := value.(type) {
	case *ListInstance:
		items := make([]string, len(v.Elements))
		for i, el := range v.Elements {
			items[i] = StringifyCompact(el)
		}
		return "[" + strings.Join(items, ", ") + "]"

	case *DictInstance:
		items := []string{}
		for k, v := range v.Entries {
			items = append(items, fmt.Sprintf("%q: %s", k, StringifyCompact(v)))
		}
		return "{" + strings.Join(items, ", ") + "}"

	default:
		return fmt.Sprintf("%v", v)
	}
}

func StringifyColor(value any, indent string) string {
	switch v := value.(type) {
	case *ListInstance:
		if len(v.Elements) == 0 {
			return "[]"
		}
		builder := strings.Builder{}
		builder.WriteString("[\n")
		for _, el := range v.Elements {
			builder.WriteString(indent + "  " + StringifyColor(el, indent+"  ") + ",\n")
		}
		builder.WriteString(indent + "]")
		return builder.String()

	case *DictInstance:
		if len(v.Entries) == 0 {
			return "{}"
		}
		builder := strings.Builder{}
		builder.WriteString("{\n")
		for k, val := range v.Entries {
			builder.WriteString(fmt.Sprintf("%s  \033[36m%q\033[0m: %s,\n", indent, k, StringifyColor(val, indent+"  ")))
		}
		builder.WriteString(indent + "}")
		return builder.String()

	default:
		// número: amarelo | string: verde
		switch v := value.(type) {
		case float64:
			return fmt.Sprintf("\033[33m%v\033[0m", v)
		case string:
			return fmt.Sprintf("\033[32m%q\033[0m", v)
		default:
			return fmt.Sprintf("%v", v)
		}
	}
}
