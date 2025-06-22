package runtime

type Callable interface {
	Call(interpreter *Interpreter, args []any) any
	Arity() int
	String() string
}
