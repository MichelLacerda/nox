package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
)

type Interpreter struct {
	runtime      *Nox
	globals      *Environment
	locals       map[Expr]int
	environment  *Environment
	stringify    func(value any) string
	silentErrors bool
	debug        bool // Modo de depuração
}

type BreakSignal struct{}

type ContinueSignal struct{}

type HasMethods interface {
	GetMethod(name string) any
}

func NewInterpreter(r *Nox, stringifyFn func(value any) string) *Interpreter {
	interpreter := &Interpreter{
		runtime:      r,
		globals:      NewEnvironment(r, nil),
		environment:  nil, // Inicialmente nil
		locals:       map[Expr]int{},
		stringify:    stringifyFn,
		silentErrors: true,  // Inicialmente não silencioso
		debug:        false, // Modo de depuração desativado por padrão
	}
	interpreter.environment = interpreter.globals // Aponta para o global no início
	interpreter.globals.Define("clock", ClockCallable{})
	interpreter.globals.Define("len", LenCallable{})
	interpreter.globals.Define("range", RangeCallable{})
	interpreter.globals.Define("assert", &BuiltinFunction{
		arity: 2,
		call: func(i *Interpreter, args []any) any {
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
			i.runtime.ReportRuntimeError(&Token{Lexeme: "assert"}, fmt.Sprintf("Assertion failed: %v", message))
			return nil
		},
	})
	interpreter.globals.Define("open", &BuiltinFunction{
		arity: 2,
		call: func(i *Interpreter, args []any) any {
			if len(args) != 2 {
				i.runtime.ReportRuntimeError(nil, "open() expects 2 arguments.")
				return nil
			}

			path, ok1 := args[0].(string)
			mode, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				i.runtime.ReportRuntimeError(&Token{Lexeme: "open"}, "open(path, mode) expects strings")
				return nil
			}

			flags, err := parseFileMode(mode)
			if err != nil {
				i.runtime.ReportRuntimeError(&Token{Lexeme: "open"}, err.Error())
				return nil
			}

			f, err := os.OpenFile(path, flags, 0666)
			if err != nil {
				i.runtime.ReportRuntimeError(&Token{Lexeme: "open"}, "failed to open file: "+err.Error())
				return nil
			}
			return &FileObject{
				file:   f,
				reader: bufio.NewReader(f),
			}
		},
	})
	return interpreter
}

func (i *Interpreter) Interpret(expr []Stmt) {
	for _, statement := range expr {
		i.execute(statement)
	}
}

func (i *Interpreter) execute(s Stmt) error {
	s.Accept(i)
	return nil
}

func (i *Interpreter) Resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

// ===== Visitor methods =====

func (i *Interpreter) VisitExpressionStmt(stmt *ExpressionStmt) any {
	result := i.evaluate(stmt.Expression)
	if result != nil {
		i.runtime.hadRuntimeError = false
	}
	return result
}

func (i *Interpreter) VisitFunctionStmt(stmt *FunctionStmt) any {
	function := NewFunction(i.runtime, stmt, i.environment, false)
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *IfStmt) any {
	condition := i.evaluate(stmt.Condition)

	if i.isTruthy(condition) {
		i.execute(stmt.Then)
	} else if stmt.Else != nil {
		i.execute(stmt.Else)
	}

	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *PrintStmt) any {
	var parts []string
	for _, expr := range stmt.Expressions {
		value := i.evaluate(expr)
		parts = append(parts, i.stringify(value))
	}
	fmt.Println(strings.Join(parts, " "))
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ReturnStmt) any {
	var value any

	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
		if value != nil {
			i.runtime.hadRuntimeError = false
		}
	}

	// Encerra a execução da função com um "Return" (que será capturado via recover)
	panic(Return{Value: value})
}

func (i *Interpreter) VisitVarStmt(stmt *VarStmt) any {
	var value any
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
		if value != nil {
			i.runtime.hadRuntimeError = false
		}
	}
	i.environment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *BlockStmt) any {
	i.executeBlock(stmt.Statements, NewEnvironment(i.runtime, i.environment))
	return nil
}

func (i *Interpreter) VisitClassStmt(stmt *ClassStmt) any {
	var superclass *Class

	if stmt.Superclass != nil {
		evaluatedSuperclass := i.evaluate(stmt.Superclass)
		if sc, ok := evaluatedSuperclass.(*Class); ok {
			superclass = sc
		} else {
			i.runtime.ReportRuntimeError(stmt.Superclass.Name, "Superclass must be a class.")
			return nil
		}
	}

	i.environment.Define(stmt.Name.Lexeme, nil) // Define a classe antes de instanciá-la

	if stmt.Superclass != nil {
		i.environment = NewEnvironment(i.runtime, i.environment) // Cria um novo ambiente para a classe
		i.environment.Define("super", superclass)                // Define a variável 'super' no ambiente da classe
	}

	methods := MethodType{}
	for _, method := range stmt.Methods {
		fn := NewFunction(i.runtime, method, i.environment, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = fn
	}

	class := NewClass(stmt.Name.Lexeme, superclass, methods)
	if stmt.Superclass != nil {
		i.environment = i.environment.Enclosing // Retorna ao ambiente anterior após definir a classe
	}
	i.environment.Assign(stmt.Name, class)
	return nil
}

func (i *Interpreter) VisitBreakStmt(stmt *BreakStmt) any {
	panic(BreakSignal{})
}

func (i *Interpreter) VisitContinueStmt(stmt *ContinueStmt) any {
	panic(ContinueSignal{})
}

func (i *Interpreter) VisitWithStmt(stmt *WithStmt) any {
	resource := i.evaluate(stmt.Resource)
	fmt.Printf("[WITH] Defining alias %s as: %T => %v\n", stmt.Alias.Lexeme, resource, resource)

	env := NewEnvironment(i.runtime, i.environment)
	env.Define(stmt.Alias.Lexeme, resource)

	defer func() {
		if file, ok := resource.(*FileObject); ok {
			if closeFn := file.GetMethod("close"); closeFn != nil {
				if callable, ok := closeFn.(Callable); ok {
					defer func() { recover() }()
					callable.Call(i, []any{})
				}
			}
		}
	}()

	i.executeBlock([]Stmt{stmt.Body}, env)
	return nil
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) error {
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

func (i *Interpreter) VisitLiteralExpr(expr *LiteralExpr) any {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr *GroupingExpr) any {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr *UnaryExpr) any {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case TokenType_MINUS:
		if !i.mustBeNumber(expr.Operator, right) {
			return nil
		}
		return -right.(float64)
	case TokenType_BANG, TokenType_NOT:
		return !i.isTruthy(right)
	}
	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr *BinaryExpr) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case TokenType_PERCENT:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		if right.(float64) == 0 {
			i.runtime.ReportRuntimeError(expr.Operator, "Division by zero.")
			return nil
		}
		return float64(int(left.(float64)) % int(right.(float64)))
	case TokenType_DOUBLE_STAR:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		if right.(float64) < 0 {
			i.runtime.ReportRuntimeError(expr.Operator, "Exponent must be a non-negative number.")
			return nil
		}
		if left.(float64) == 0 && right.(float64) == 0 {
			i.runtime.ReportRuntimeError(expr.Operator, "0 raised to the power of 0 is undefined.")
			return nil
		}
		return math.Pow(left.(float64), right.(float64))
	case TokenType_MINUS:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		return left.(float64) - right.(float64)

	case TokenType_SLASH:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		if right.(float64) == 0 {
			i.runtime.ReportRuntimeError(expr.Operator, "Division by zero.")
			return nil
		}
		return left.(float64) / right.(float64)

	case TokenType_STAR:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		return left.(float64) * right.(float64)

	case TokenType_PLUS:
		switch l := left.(type) {
		case float64:
			if r, ok := right.(float64); ok {
				return l + r
			}
		case string:
			return l + i.stringify(right)
		case nil:
			if r, ok := right.(string); ok {
				return "nil" + r
			}
		default:
			if r, ok := right.(string); ok {
				return i.stringify(left) + r
			}
		}

		if _, lok := left.(string); lok {
			return i.stringify(left) + i.stringify(right)
		}

		if _, rok := right.(string); rok {
			return i.stringify(left) + i.stringify(right)
		}

		i.runtime.ReportRuntimeError(expr.Operator, fmt.Sprintf(
			"Operands must be two numbers or two strings, but got [%T] and [%T].", left, right))
		return nil

	case TokenType_GREATER:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		return left.(float64) > right.(float64)

	case TokenType_GREATER_EQUAL:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		return left.(float64) >= right.(float64)

	case TokenType_LESS:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		return left.(float64) < right.(float64)

	case TokenType_LESS_EQUAL:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		return left.(float64) <= right.(float64)

	case TokenType_EQUAL_EQUAL:
		return i.isEqual(left, right)

	case TokenType_BANG_EQUAL:
		return !i.isEqual(left, right)
	}

	return nil
}

// Placeholder visitors (ainda não implementados)
func (i *Interpreter) VisitVariableExpr(expr *VariableExpr) any {
	// return i.environments.Get(expr.Name.Lexeme)
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) lookUpVariable(t *Token, expr Expr) any {
	if depth, ok := i.locals[expr]; ok {
		return i.environment.GetAt(depth, t.Lexeme)
	}
	return i.globals.Get(t)
}

func (i *Interpreter) VisitAssignExpr(expr *AssignExpr) any {
	value := i.evaluate(expr.Value)

	if d, ok := i.locals[expr]; ok {
		i.environment.AssignAt(d, expr.Name, value)
	} else {
		i.globals.Assign(expr.Name, value)
	}

	return nil
}

func (i *Interpreter) VisitCallExpr(expr *CallExpr) any {
	callee := i.evaluate(expr.Callee)

	if callee == nil {
		i.runtime.ReportRuntimeError(expr.Parenthesis, "Attempt to call method on nil.")
		return nil
	}

	var arguments []any
	for _, argument := range expr.Arguments {
		arguments = append(arguments, i.evaluate(argument))
	}

	callable, ok := callee.(Callable)
	if !ok {
		i.runtime.ReportRuntimeError(expr.Parenthesis, fmt.Sprintf("Can only call functions and classes. %T", callee))
		return nil
	}

	arity := callable.Arity()
	if arity >= 0 && len(arguments) != arity {
		i.runtime.ReportRuntimeError(expr.Parenthesis, fmt.Sprintf(
			"Expected %d arguments but got %d.",
			arity, len(arguments),
		))
		return nil
	}

	return callable.Call(i, arguments)
}

func (i *Interpreter) VisitGetExpr(expr *GetExpr) any {
	object := i.evaluate(expr.Object)

	switch obj := object.(type) {
	case *Instance:
		return obj.Get(expr.Name)

	case *ListInstance:
		return obj.Get(expr.Name)

	case *DictInstance:
		return obj.Get(expr.Name)

	case *FileObject:
		method := obj.GetMethod(expr.Name.Lexeme)
		if method != nil {
			return method
		}
		i.runtime.ReportRuntimeError(expr.Name,
			fmt.Sprintf("Undefined property '%s' for file object.", expr.Name.Lexeme))
		return nil

	default:
		i.runtime.ReportRuntimeError(expr.Name,
			"Only instances, lists, or dicts have properties.")
		return nil
	}
}

func (i *Interpreter) VisitSetExpr(expr *SetExpr) any {
	object := i.evaluate(expr.Object)

	if instance, ok := object.(*Instance); ok {
		value := i.evaluate(expr.Value)
		instance.Set(expr.Name, value)
		return value
	}

	i.runtime.ReportRuntimeError(expr.Name, "Only instances have fields.")
	return nil
}

func (i *Interpreter) VisitLogicalExpr(expr *LogicalExpr) any {
	left := i.evaluate(expr.Left)

	switch expr.Operator.Type {
	case TokenType_OR:
		if i.isTruthy(left) {
			return left
		}
	case TokenType_AND:
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitSuperExpr(expr *SuperExpr) any {

	// Recupera a distância do 'super' na resolução de variáveis (enquanto resolver)
	distance := i.locals[expr]

	// Busca a classe pai (superclasse) a partir da distância
	superclass, ok := i.environment.GetAt(distance, "super").(*Class)

	if !ok {
		i.runtime.ReportRuntimeError(expr.Keyword, "Invalid superclass.")
		return nil
	}

	// Busca a instância atual (this/self), que está um escopo acima
	object, ok := i.environment.GetAt(distance-1, "self").(*Instance)
	if !ok {
		i.runtime.ReportRuntimeError(expr.Keyword, "Invalid instance for 'super'.")
		return nil
	}

	// Tenta localizar o método na superclasse
	method, found := superclass.FindMethod(expr.Method.Lexeme)
	if !found {
		i.runtime.ReportRuntimeError(expr.Method, fmt.Sprintf(
			"Undefined property '%s'.", expr.Method.Lexeme))
		return nil
	}

	// Retorna o método ligado à instância (bind)
	return method.Bind(object)
}

func (i *Interpreter) VisitSelfExpr(expr *SelfExpr) any {
	value := i.lookUpVariable(expr.Keyword, expr)
	return value
}

// func (i *Interpreter) VisitWhileStmt(stmt *WhileStmt) any {
// 	previous := i.environment
// 	defer func() { i.environment = previous }()

// 	for i.isTruthy(i.evaluate(stmt.Condition)) {
// 		i.execute(stmt.Body)
// 		if i.runtime.hadRuntimeError {
// 			return nil
// 		}
// 	}
// 	return nil
// }

func (i *Interpreter) VisitWhileStmt(stmt *WhileStmt) any {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					switch r.(type) {
					case BreakSignal:
						panic(r) // quebra o laço externo
					case ContinueSignal:
						// simplesmente ignora, continua o loop
					default:
						panic(r)
					}
				}
			}()
			i.executeBlock([]Stmt{stmt.Body}, i.environment)
		}()
	}
	return nil
}

func (i *Interpreter) VisitForInStmt(stmt *ForInStmt) any {
	iterable := i.evaluate(stmt.Iterable)

	switch coll := iterable.(type) {
	case []any: // list
		for index, value := range coll {
			env := NewEnvironment(i.runtime, i.environment)

			if stmt.IndexVar != nil {
				env.Define(stmt.IndexVar.Lexeme, float64(index))
			}
			env.Define(stmt.ValueVar.Lexeme, value)

			func() {
				defer func() {
					if r := recover(); r != nil {
						switch r.(type) {
						case BreakSignal:
							panic(r) // repassa pro loop pai
						case ContinueSignal:
							// ignora, continua o próximo item
						default:
							panic(r)
						}
					}
				}()

				i.executeBlock([]Stmt{stmt.Body}, env)
			}()
		}

	case map[string]any: // dict
		for key, value := range coll {
			env := NewEnvironment(i.runtime, i.environment)

			if stmt.IndexVar != nil {
				env.Define(stmt.IndexVar.Lexeme, key)
			}
			env.Define(stmt.ValueVar.Lexeme, value)

			func() {
				defer func() {
					if r := recover(); r != nil {
						switch r.(type) {
						case BreakSignal:
							panic(r)
						case ContinueSignal:
							// ignora
						default:
							panic(r)
						}
					}
				}()

				i.executeBlock([]Stmt{stmt.Body}, env)
			}()
		}

	case bool: // for { ... } → loop infinito
		if coll {
			for {
				env := NewEnvironment(i.runtime, i.environment)
				func() {
					defer func() {
						if r := recover(); r != nil {
							switch r.(type) {
							case BreakSignal:
								panic(r)
							case ContinueSignal:
								// ignora
							default:
								panic(r)
							}
						}
					}()

					i.executeBlock([]Stmt{stmt.Body}, env)
				}()
			}
		}

	default:
		i.runtime.ReportRuntimeError(stmt.ValueVar, "Object is not iterable.")
	}

	return nil
}

func (i *Interpreter) VisitListExpr(expr *ListExpr) any {
	var result []any
	for _, element := range expr.Elements {
		value := i.evaluate(element)
		result = append(result, value)
	}
	return NewListInstance(result)
}

func (i *Interpreter) VisitIndexExpr(expr *IndexExpr) any {
	object := i.evaluate(expr.List)
	index := i.evaluate(expr.Index)

	switch obj := object.(type) {
	case []any:
		intIndex, ok := index.(float64)
		if !ok {
			i.runtime.ReportRuntimeError(&Token{Type: TokenType_Unknown, Lexeme: "[]"}, "List index must be a number.")
			return nil
		}
		idx := int(intIndex)
		if idx < 0 || idx >= len(obj) {
			i.runtime.ReportRuntimeError(&Token{Type: TokenType_Unknown, Lexeme: "[]"}, fmt.Sprintf("List index out of range: %d", idx))
			return nil
		}
		return obj[idx]
	case map[string]any: // dicionário
		key, ok := index.(string)
		if !ok {
			i.runtime.ReportRuntimeError(&Token{Type: TokenType_Unknown, Lexeme: "{}"}, "Dictionary keys must be strings.")
			return nil
		}
		val, exists := obj[key]
		if !exists {
			i.runtime.ReportRuntimeError(&Token{Type: TokenType_Unknown, Lexeme: "{}"}, fmt.Sprintf("Key '%s' not found in dictionary.", key))
			return nil
		}
		return val
	case *ListInstance:
		intIndex, ok := index.(float64)
		if !ok {
			i.runtime.ReportRuntimeError(&Token{Type: TokenType_Unknown, Lexeme: "[]"}, "List index must be a number.")
			return nil
		}
		idx := int(intIndex)
		if idx < 0 || idx >= len(obj.Elements) {
			i.runtime.ReportRuntimeError(&Token{Type: TokenType_Unknown, Lexeme: "[]"}, fmt.Sprintf("List index out of range: %d", idx))
			return nil
		}
		return obj.Elements[idx]

	case *DictInstance:
		key, ok := index.(string)
		if !ok {
			i.runtime.ReportRuntimeError(&Token{Type: TokenType_Unknown, Lexeme: "{}"}, "Dictionary keys must be strings.")
			return nil
		}
		val, exists := obj.Entries[key]
		if !exists {
			i.runtime.ReportRuntimeError(&Token{Type: TokenType_Unknown, Lexeme: "{}"}, fmt.Sprintf("Key '%s' not found in dictionary.", key))
			return nil
		}
		return val
	default:
		i.runtime.ReportRuntimeError(&Token{Type: TokenType_Unknown, Lexeme: ""}, "Only lists and dictionaries support indexing.")
		return nil

	}
}

func (i *Interpreter) VisitDictExpr(expr *DictExpr) any {
	dict := map[string]any{}

	for _, pair := range expr.Pairs {
		key := i.evaluate(pair.Key)
		value := i.evaluate(pair.Value)

		if keyStr, ok := key.(string); ok {
			dict[keyStr] = value
		} else {
			i.runtime.ReportRuntimeError(&Token{Type: TokenType_Unknown, Lexeme: ""}, "Dictionary keys must be strings.")
			return nil
		}
	}
	return NewDictInstance(dict)
}

func (i *Interpreter) VisitSafeExpr(expr *SafeExpr) any {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(RuntimeError); ok {
				// erro controlado — retorna nil silenciosamente
				return
			}
			panic(r) // outros panics inesperados continuam
		}
	}()

	return i.evaluate(expr.Expr)
}

// ===== Helpers =====

func (i *Interpreter) evaluate(expr Expr) any {
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

// ===== Validação de operandos =====

func (i *Interpreter) mustBeNumber(op *Token, val any) bool {
	if _, ok := val.(float64); !ok {
		i.runtime.ReportRuntimeError(op, "Operand must be a number.")
		return false
	}
	return true
}

func (i *Interpreter) mustBeNumbers(op *Token, left, right any) bool {
	if _, ok := left.(float64); !ok {
		i.runtime.ReportRuntimeError(op, "Left operand must be a number.")
		return false
	}
	if _, ok := right.(float64); !ok {
		i.runtime.ReportRuntimeError(op, "Right operand must be a number.")
		return false
	}
	return true
}
