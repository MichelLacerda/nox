package runtime

import (
	"fmt"

	"github.com/MichelLacerda/nox/internal/token"
)

type RuntimeError struct {
	Token   *token.Token
	Message string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("Runtime Error at %s: %s", e.Token.Lexeme, e.Message)
}

func NewRuntimeError(token *token.Token, message string) RuntimeError {
	return RuntimeError{
		Token:   token,
		Message: message,
	}
}

func (n *Interpreter) ReportError(line int, where, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
}

func (n *Interpreter) ErrorAt(line int, message string) {
	n.ReportError(line, "", message)
}
