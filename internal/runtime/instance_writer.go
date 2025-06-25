package runtime

import (
	"fmt"
	"net/http"

	"github.com/MichelLacerda/nox/internal/token"
)

type WriterInstance struct {
	w http.ResponseWriter
}

func (w *WriterInstance) Get(name *token.Token) any {
	switch name.Lexeme {
	case "write":
		return &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				fmt.Fprint(w.w, args[0])
				return nil
			},
		}
	default:
		return nil
	}
}

func (w *WriterInstance) String() string {
	return "<writer>"
}
