package main

import "fmt"

type RangeCallable struct{}

func (r RangeCallable) Arity() int {
	// aceita 1, 2 ou 3 argumentos
	return -1 // Arity flexível
}

func (r RangeCallable) Call(interpreter *Interpreter, arguments []any) any {
	var start, end, step float64

	switch len(arguments) {
	case 1:
		start = 0
		end = toFloat(arguments[0])
		step = 1
	case 2:
		start = toFloat(arguments[0])
		end = toFloat(arguments[1])
		step = 1
	case 3:
		start = toFloat(arguments[0])
		end = toFloat(arguments[1])
		step = toFloat(arguments[2])
	default:
		interpreter.runtime.ReportRuntimeError(nil, "range() expects 1 to 3 arguments.")
		return nil
	}

	if step == 0 {
		interpreter.runtime.ReportRuntimeError(nil, "range() step must not be zero.")
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

	return result
}

func (r RangeCallable) String() string {
	return "<builtin fn range>"
}

// Utilitário para conversão com fallback
func toFloat(v any) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	default:
		panic(fmt.Sprintf("Invalid type for range argument: %T", val))
	}
}
