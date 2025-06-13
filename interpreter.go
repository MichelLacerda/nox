package main

import (
	"fmt"
	"log"
	"reflect"
)

type Interpreter struct {
	runtime      *Nox
	globals      *Environment
	locals       map[Expr]int
	environments *Environment
}

func NewInterpreter(r *Nox) *Interpreter {
	interpreter := &Interpreter{
		runtime:      r,
		globals:      NewEnvironment(r, nil),
		environments: nil, // Inicialmente nil
		locals:       map[Expr]int{},
	}
	interpreter.environments = interpreter.globals // Aponta para o global no início
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
	result := i.eval(stmt.Expression)
	if result != nil {
		i.runtime.hadRuntimeError = false
	}
	return result
}

func (i *Interpreter) VisitFunctionStmt(stmt *FunctionStmt) any {
	function := NewNoxFunction(stmt, i.environments)
	i.environments.Define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *IfStmt) any {
	condition := i.eval(stmt.Condition)

	if i.isTruthy(condition) {
		i.execute(stmt.Then)
	} else if stmt.Else != nil {
		i.execute(stmt.Else)
	}

	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *PrintStmt) any {
	result := i.eval(stmt.Expression)
	if result != nil {
		i.runtime.hadRuntimeError = false
	}
	fmt.Println(i.stringify(result))
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ReturnStmt) any {
	var value any

	if stmt.Value != nil {
		value = i.eval(stmt.Value)
	}

	if value != nil {
		i.runtime.hadRuntimeError = false
	}

	panic(Return{Value: value})
}

func (i *Interpreter) VisitVarStmt(stmt *VarStmt) any {
	var value any
	if stmt.Initializer != nil {
		value = i.eval(stmt.Initializer)
		if value != nil {
			i.runtime.hadRuntimeError = false
		}
	}
	i.environments.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *BlockStmt) any {
	i.executeBlock(stmt.Statements, NewEnvironment(i.runtime, i.environments))
	return nil
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) error {
	previous := i.environments
	i.environments = environment
	defer func() {
		i.environments = previous
	}()

	for _, statement := range statements {
		if err := i.execute(statement); err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) VisitLiteralExpr(expr *LiteralExpr) any {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr *GroupingExpr) any {
	return i.eval(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr *UnaryExpr) any {
	right := i.eval(expr.Right)

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
	left := i.eval(expr.Left)
	right := i.eval(expr.Right)

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
			if r, ok := right.(string); ok {
				return l + r
			}
		}
		i.runtime.ReportRuntimeError(expr.Operator, "Operands must be two numbers or two strings.")
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
		return i.environments.GetAt(depth, t.Lexeme)
	}
	return i.globals.Get(t)
}

func (i *Interpreter) VisitAssignExpr(expr *AssignExpr) any {
	value := i.eval(expr.Value)

	if d, ok := i.locals[expr]; ok {
		i.environments.AssignAt(d, expr.Name, value)
	} else {
		i.globals.Assign(expr.Name, value)
	}

	return nil
}

func (i *Interpreter) VisitCallExpr(expr *CallExpr) any {
	callee := i.eval(expr.Callee)

	arguments := make([]any, len(expr.Arguments))
	for idx, arg := range expr.Arguments {
		arguments[idx] = i.eval(arg)
	}

	function, ok := callee.(NoxCallable)
	if !ok {
		i.runtime.ReportRuntimeError(expr.Parenthesis, "Can only call functions and classes.")
		return nil // Corrige panic ao tentar chamar valor não callable
	}

	if len(arguments) != function.Arity() {
		i.runtime.ReportRuntimeError(expr.Parenthesis, fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments)))
		return nil
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) VisitGetExpr(expr *GetExpr) any {
	log.Panic("VisitGetExpr not implemented yet.")
	return nil
}

func (i *Interpreter) VisitSetExpr(expr *SetExpr) any {
	log.Panic("VisitSetExpr not implemented yet.")
	return nil
}

func (i *Interpreter) VisitLogicalExpr(expr *LogicalExpr) any {
	left := i.eval(expr.Left)

	if expr.Operator.Type == TokenType_OR {
		if i.isTruthy(left) {
			return left
		}
	} else if expr.Operator.Type == TokenType_AND {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.eval(expr.Right)
}

func (i *Interpreter) VisitSuperExpr(expr *SuperExpr) any {
	log.Panic("VisitSuperExpr not implemented yet.")
	return nil
}

func (i *Interpreter) VisitThisExpr(expr *ThisExpr) any {
	log.Panic("VisitThisExpr not implemented yet.")
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *WhileStmt) any {
	previous := i.environments
	defer func() { i.environments = previous }()

	for i.isTruthy(i.eval(stmt.Condition)) {
		i.execute(stmt.Body)
		if i.runtime.hadRuntimeError {
			return nil
		}
	}
	return nil
}

// ===== Helpers =====

func (i *Interpreter) eval(expr Expr) any {
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
