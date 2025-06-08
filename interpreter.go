package main

import (
	"fmt"
	"log"
	"reflect"
)

type Interpreter struct {
	vm *Nox
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

// Interpret executa uma expressão e imprime o resultado.
func (i *Interpreter) Interpret(vm *Nox, expr Expr) {
	i.vm = vm
	result := i.eval(expr)
	if result != nil {
		i.vm.hadRuntimeError = false
	}
	fmt.Println("Result:", i.stringify(result))
}

// ===== Visitor methods =====

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
			i.vm.ReportRuntimeError(expr.Operator, "Division by zero.")
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
		i.vm.ReportRuntimeError(expr.Operator, "Operands must be two numbers or two strings.")
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
	log.Panic("VisitVariableExpr not implemented yet.")
	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *AssignExpr) any {
	log.Panic("VisitAssignExpr not implemented yet.")
	return nil
}

func (i *Interpreter) VisitCallExpr(expr *CallExpr) any {
	log.Panic("VisitCallExpr not implemented yet.")
	return nil
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
	log.Panic("VisitLogicalExpr not implemented yet.")
	return nil
}

func (i *Interpreter) VisitSuperExpr(expr *SuperExpr) any {
	log.Panic("VisitSuperExpr not implemented yet.")
	return nil
}

func (i *Interpreter) VisitThisExpr(expr *ThisExpr) any {
	log.Panic("VisitThisExpr not implemented yet.")
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
		i.vm.ReportRuntimeError(op, "Operand must be a number.")
		return false
	}
	return true
}

func (i *Interpreter) mustBeNumbers(op *Token, left, right any) bool {
	if _, ok := left.(float64); !ok {
		i.vm.ReportRuntimeError(op, "Left operand must be a number.")
		return false
	}
	if _, ok := right.(float64); !ok {
		i.vm.ReportRuntimeError(op, "Right operand must be a number.")
		return false
	}
	return true
}
