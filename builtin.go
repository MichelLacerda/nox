package main

type BuiltinFunction struct {
	arity int
	call  func(interpreter *Interpreter, args []any) any
}

// ✅ Agora implementa Callable corretamente
func (b *BuiltinFunction) Arity() int {
	return b.arity
}

func (b *BuiltinFunction) Call(interpreter *Interpreter, args []any) any {
	return b.call(interpreter, args)
}

func (b *BuiltinFunction) String() string {
	return "<builtin fn>"
}
