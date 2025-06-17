package main

import "fmt"

type LenCallable struct{}

func (c LenCallable) Arity() int {
	return 1
}

// func (c LenCallable) Call(interpreter *Interpreter, args []any) any {
// 	if list, ok := args[0].([]any); ok {
// 		return float64(len(list))
// 	}
// 	return nil
// }

func (l LenCallable) Call(i *Interpreter, args []any) any {
	arg := args[0]
	switch v := arg.(type) {
	case string:
		return float64(len(v))
	case []any:
		return float64(len(v))
	default:
		i.runtime.ReportRuntimeError(&Token{
			Lexeme: "len",
			Type:   TokenType_NIL,
		}, fmt.Sprintf("len() not supported for type %T", arg))
		return nil
	}
}

func (c LenCallable) String() string {
	return "<built-in function len>"
}

func (c LenCallable) Bind(instance *Instance) Callable {
	return c
}
