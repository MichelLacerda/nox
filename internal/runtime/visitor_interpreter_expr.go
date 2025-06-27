package runtime

import (
	"fmt"
	"math"

	"github.com/MichelLacerda/nox/internal/ast"
	"github.com/MichelLacerda/nox/internal/token"
)

func (i *Interpreter) VisitAssignExpr(expr *ast.AssignExpr) any {
	value := i.evaluate(expr.Value)

	if d, ok := i.locals[expr]; ok {
		i.environment.AssignAt(d, expr.Name, value)
	} else {
		i.globals.Assign(expr.Name, value)
	}

	return nil
}

func (i *Interpreter) VisitCallExpr(expr *ast.CallExpr) any {
	callee := i.evaluate(expr.Callee)

	// fmt.Printf("Visiting CallExpr: callee=%T, arguments=%v\n", callee, expr.Arguments)
	if callee == nil {
		i.Runtime.ReportRuntimeError(expr.Parenthesis, "Attempt to call method on nil.")
		return nil
	}

	var arguments []any
	for _, argument := range expr.Arguments {
		arguments = append(arguments, i.evaluate(argument))
	}

	callable, ok := callee.(Callable)
	if !ok {
		i.Runtime.ReportRuntimeError(expr.Parenthesis, fmt.Sprintf("Can only call functions and classes. %T", callee))
		return nil
	}

	return callable.Call(i, arguments)
}

func (i *Interpreter) VisitGetExpr(expr *ast.GetExpr) any {
	object := i.evaluate(expr.Object)
	switch obj := object.(type) {
	case *Instance:
		return obj.Get(expr.Name)
	case *ListInstance:
		return obj.Get(expr.Name)
	case *DictInstance:
		if val, ok := obj.Entries[expr.Name.Lexeme]; ok {
			return val
		}
		return obj.Get(expr.Name)
	case string:
		o := &StringInstance{Value: obj}
		if method := o.GetMethod(expr.Name.Lexeme); method != nil {
			return method
		}
		i.Runtime.ReportRuntimeError(
			expr.Name,
			fmt.Sprintf("Undefined property '%s' for string object.", expr.Name.Lexeme),
		)
		return nil
	case *StringInstance:
		if method := obj.GetMethod(expr.Name.Lexeme); method != nil {
			return method
		}
		i.Runtime.ReportRuntimeError(
			expr.Name,
			fmt.Sprintf("Undefined property '%s' for string object.", expr.Name.Lexeme),
		)
		return nil
	case *FileObject:
		if method := obj.GetMethod(expr.Name.Lexeme); method != nil {
			return method
		}
		i.Runtime.ReportRuntimeError(
			expr.Name,
			fmt.Sprintf("Undefined property '%s' for file object.", expr.Name.Lexeme),
		)
		return nil
	case *EnvironmentWrapper:
		if val, ok := obj.Env.Values[expr.Name.Lexeme]; ok {
			return val
		}
		i.Runtime.ReportRuntimeError(
			expr.Name,
			fmt.Sprintf("Undefined property '%s' in module.", expr.Name.Lexeme),
		)
		return nil

	case *MapInstance:
		if method := obj.Get(expr.Name); method != nil {
			return method
		}
		i.Runtime.ReportRuntimeError(
			expr.Name,
			fmt.Sprintf("Undefined property '%s' for map object.", expr.Name.Lexeme),
		)
		return nil

	case *WriterInstance:
		if method := obj.Get(expr.Name); method != nil {
			return method
		}
		i.Runtime.ReportRuntimeError(
			expr.Name,
			fmt.Sprintf("Undefined property '%s' for writer object.", expr.Name.Lexeme),
		)
		return nil

	default:
		i.Runtime.ReportRuntimeError(
			expr.Name,
			"Only instances, lists, dicts, or modules have properties.",
		)
		return nil
	}
}

func (i *Interpreter) VisitSetExpr(expr *ast.SetExpr) any {
	object := i.evaluate(expr.Object)

	if instance, ok := object.(*Instance); ok {
		value := i.evaluate(expr.Value)
		instance.Set(expr.Name, value)
		return value
	}

	i.Runtime.ReportRuntimeError(expr.Name, "Only instances have fields.")
	return nil
}

func (i *Interpreter) VisitSetIndexExpr(expr *ast.SetIndexExpr) any {
	object := i.evaluate(expr.Object)
	index := i.evaluate(expr.Index)
	value := i.evaluate(expr.Value)
	switch obj := object.(type) {
	case []any:
		intIndex, ok := index.(float64)
		if !ok {
			i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: "[]"}, "List index must be a number.")
			return nil
		}
		idx := int(intIndex)
		if idx < 0 || idx >= len(obj) {
			i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: "[]"}, fmt.Sprintf("List index out of range: %d", idx))
			return nil
		}
		obj[idx] = value
		return value
	case map[string]any: // dicionário
		key, ok := index.(string)
		if !ok {
			i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: "{}"}, "Dictionary keys must be strings.")
			return nil
		}
		obj[key] = value
		return value
	case *ListInstance:
		intIndex, ok := index.(float64)
		if !ok {
			i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: "[]"}, "List index must be a number.")
			return nil
		}
		idx := int(intIndex)
		if idx < 0 || idx >= len(obj.Elements) {
			i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: "[]"}, fmt.Sprintf("List index out of range: %d", idx))
			return nil
		}
		obj.Elements[idx] = value
		return value
	case *DictInstance:
		key, ok := index.(string)
		if !ok {
			i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: "{}"}, "Dictionary keys must be strings.")
			return nil
		}
		obj.Entries[key] = value
		return value
	default:
		i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: "[]"}, "Only lists and dictionaries support indexing.")
		return nil
	}
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.LogicalExpr) any {
	left := i.evaluate(expr.Left)

	switch expr.Operator.Type {
	case token.TokenType_OR:
		if i.isTruthy(left) {
			return left
		}
	case token.TokenType_AND:
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitSuperExpr(expr *ast.SuperExpr) any {

	// Recupera a distância do 'super' na resolução de variáveis (enquanto resolver)
	distance := i.locals[expr]

	// Busca a classe pai (superclasse) a partir da distância
	superclass, ok := i.environment.GetAt(distance, "super").(*Class)

	if !ok {
		i.Runtime.ReportRuntimeError(expr.Keyword, "Invalid superclass.")
		return nil
	}

	// Busca a instância atual (this/self), que está um escopo acima
	object, ok := i.environment.GetAt(distance-1, "self").(*Instance)
	if !ok {
		i.Runtime.ReportRuntimeError(expr.Keyword, "Invalid instance for 'super'.")
		return nil
	}

	// Tenta localizar o método na superclasse
	method, found := superclass.FindMethod(expr.Method.Lexeme)
	if !found {
		i.Runtime.ReportRuntimeError(expr.Method, fmt.Sprintf(
			"Undefined property '%s'.", expr.Method.Lexeme))
		return nil
	}

	// Retorna o método ligado à instância (bind)
	return method.Bind(object)
}

func (i *Interpreter) VisitSelfExpr(expr *ast.SelfExpr) any {
	value := i.lookUpVariable(expr.Keyword, expr)
	return value
}

func (i *Interpreter) VisitListExpr(expr *ast.ListExpr) any {
	var result []any
	for _, element := range expr.Elements {
		value := i.evaluate(element)
		result = append(result, value)
	}
	return NewListInstance(result)
}

func (i *Interpreter) VisitIndexExpr(expr *ast.IndexExpr) any {
	object := i.evaluate(expr.Object)
	index := i.evaluate(expr.Index)

	switch obj := object.(type) {
	case []any:
		intIndex, ok := index.(float64)
		if !ok {
			i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: "[]"}, "List index must be a number.")
			return nil
		}
		idx := int(intIndex)
		if idx < 0 || idx >= len(obj) {
			i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: "[]"}, fmt.Sprintf("List index out of range: %d", idx))
			return nil
		}
		return obj[idx]
	case map[string]any: // dicionário
		key, ok := index.(string)
		if !ok {
			i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: "{}"}, "Dictionary keys must be strings.")
			return nil
		}
		val, exists := obj[key]
		if !exists {
			i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: "{}"}, fmt.Sprintf("Key '%s' not found in dictionary.", key))
			return nil
		}
		return val
	case *ListInstance:
		intIndex, ok := index.(float64)
		if !ok {
			i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: "[]"}, "List index must be a number.")
			return nil
		}
		idx := int(intIndex)
		if idx < 0 || idx >= len(obj.Elements) {
			i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: "[]"}, fmt.Sprintf("List index out of range: %d", idx))
			return nil
		}
		return obj.Elements[idx]

	case *DictInstance:
		key, ok := index.(string)
		if !ok {
			i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: "{}"}, "Dictionary keys must be strings.")
			return nil
		}
		val, exists := obj.Entries[key]
		if !exists {
			i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: "{}"}, fmt.Sprintf("Key '%s' not found in dictionary.", key))
			return nil
		}
		return val
	default:
		i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: ""}, "Only lists and dictionaries support indexing.")
		return nil

	}
}

func (i *Interpreter) VisitDictExpr(expr *ast.DictExpr) any {
	dict := map[string]any{}

	for _, pair := range expr.Pairs {
		key := i.evaluate(pair.Key)
		value := i.evaluate(pair.Value)

		if keyStr, ok := key.(string); ok {
			dict[keyStr] = value
		} else {
			i.Runtime.ReportRuntimeError(&token.Token{Type: token.TokenType_Unknown, Lexeme: ""}, "Dictionary keys must be strings.")
			return nil
		}
	}
	return NewDictInstance(dict)
}

func (i *Interpreter) VisitSafeExpr(expr *ast.SafeExpr) any {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(*RuntimeError); ok {
				// erro controlado — retorna nil silenciosamente
				return
			}
			panic(r) // outros panics inesperados continuam
		}
	}()

	return i.evaluate(expr.Expr)
}

func (i *Interpreter) VisitVariableExpr(expr *ast.VariableExpr) any {
	// return i.environments.Get(expr.Name.Lexeme)
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) any {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) any {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) any {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.TokenType_MINUS:
		if !i.mustBeNumber(expr.Operator, right) {
			return nil
		}
		return -right.(float64)
	case token.TokenType_BANG, token.TokenType_NOT:
		return !i.isTruthy(right)
	}
	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.TokenType_PERCENT:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		if right.(float64) == 0 {
			i.Runtime.ReportRuntimeError(expr.Operator, "Division by zero.")
			return nil
		}
		return float64(int(left.(float64)) % int(right.(float64)))
	case token.TokenType_DOUBLE_STAR:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		if right.(float64) < 0 {
			i.Runtime.ReportRuntimeError(expr.Operator, "Exponent must be a non-negative number.")
			return nil
		}
		if left.(float64) == 0 && right.(float64) == 0 {
			i.Runtime.ReportRuntimeError(expr.Operator, "0 raised to the power of 0 is undefined.")
			return nil
		}
		return math.Pow(left.(float64), right.(float64))
	case token.TokenType_MINUS:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		return left.(float64) - right.(float64)

	case token.TokenType_SLASH:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		if right.(float64) == 0 {
			i.Runtime.ReportRuntimeError(expr.Operator, "Division by zero.")
			return nil
		}
		return left.(float64) / right.(float64)

	case token.TokenType_STAR:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		return left.(float64) * right.(float64)

	case token.TokenType_PLUS:
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

		i.Runtime.ReportRuntimeError(expr.Operator, fmt.Sprintf(
			"Operands must be two numbers or two strings, but got [%T] and [%T].", left, right))
		return nil

	case token.TokenType_GREATER:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		return left.(float64) > right.(float64)

	case token.TokenType_GREATER_EQUAL:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		return left.(float64) >= right.(float64)

	case token.TokenType_LESS:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		return left.(float64) < right.(float64)

	case token.TokenType_LESS_EQUAL:
		if !i.mustBeNumbers(expr.Operator, left, right) {
			return nil
		}
		return left.(float64) <= right.(float64)

	case token.TokenType_EQUAL_EQUAL:
		return i.isEqual(left, right)

	case token.TokenType_BANG_EQUAL:
		return !i.isEqual(left, right)
	}

	return nil
}
