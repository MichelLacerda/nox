package main

import (
	"fmt"
	"log"
	"reflect"
)

type Interpreter struct {
	runtime     *Nox
	globals     *Environment
	locals      map[Expr]int
	environment *Environment
}

func NewInterpreter(r *Nox) *Interpreter {
	interpreter := &Interpreter{
		runtime:     r,
		globals:     NewEnvironment(r, nil),
		environment: nil, // Inicialmente nil
		locals:      map[Expr]int{},
	}
	interpreter.environment = interpreter.globals // Aponta para o global no início
	interpreter.globals.Define("clock", ClockCallable{})
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
	result := i.evaluate(stmt.Expression)
	if result != nil {
		i.runtime.hadRuntimeError = false
	}
	fmt.Println(i.stringify(result))
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
	i.environment.Define(stmt.Name.Lexeme, nil) // Define a classe antes de instanciá-la

	methods := MethodType{}
	for _, method := range stmt.Methods {
		fn := NewFunction(i.runtime, method, i.environment, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = fn
	}
	class := NewClass(stmt.Name.Lexeme, methods)
	i.environment.Assign(stmt.Name, class)
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
	case TokenType_BANG:
		return !i.isTruthy(right)
	}
	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr *BinaryExpr) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
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

	var arguments []any
	for _, argument := range expr.Arguments {
		arguments = append(arguments, i.evaluate(argument))
	}

	callable, ok := callee.(Callable)
	if !ok {
		i.runtime.ReportRuntimeError(expr.Parenthesis, "Can only call functions and classes.")
		return nil
	}

	if len(arguments) != callable.Arity() {
		i.runtime.ReportRuntimeError(expr.Parenthesis, fmt.Sprintf(
			"Expected %d arguments but got %d.",
			callable.Arity(), len(arguments),
		))
		return nil
	}

	return callable.Call(i, arguments)
}

func (i *Interpreter) VisitGetExpr(expr *GetExpr) any {
	object := i.evaluate(expr.Object)
	if instance, ok := object.(*Instance); ok {
		value := instance.Get(expr.Name)
		return value
	}

	if instance, ok := object.(*Instance); ok {
		return instance.Get(expr.Name)
	}

	i.runtime.ReportRuntimeError(expr.Name, "Only instances have properties.")
	return nil
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

	if expr.Operator.Type == TokenType_OR {
		if i.isTruthy(left) {
			return left
		}
	} else if expr.Operator.Type == TokenType_AND {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitSuperExpr(expr *SuperExpr) any {
	log.Panic("VisitSuperExpr not implemented yet.")
	return nil
}

func (i *Interpreter) VisitSelfExpr(expr *SelfExpr) any {
	value := i.lookUpVariable(expr.Keyword, expr)
	return value
}

func (i *Interpreter) VisitWhileStmt(stmt *WhileStmt) any {
	previous := i.environment
	defer func() { i.environment = previous }()

	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
		if i.runtime.hadRuntimeError {
			return nil
		}
	}
	return nil
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

func (i *Interpreter) stringify(value any) string {
	switch v := value.(type) {
	case nil:
		return "nil"
	case float64:
		return fmt.Sprintf("%g", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
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
