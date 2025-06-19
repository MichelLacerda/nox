package main

import (
	"time"
)

type ClockCallable struct{}

func (c ClockCallable) Arity() int {
	return 0
}

func (c ClockCallable) Call(interpreter *Interpreter, arguments []any) any {
	return float64(time.Now().UnixNano()) / 1e9 // segundos com fração
}

func (c ClockCallable) String() string {
	return "<native func>"
}

func (c ClockCallable) Bind(instance *Instance) Callable {
	return c
}
