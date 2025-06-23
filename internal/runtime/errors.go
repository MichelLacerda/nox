package runtime

import (
	"fmt"

	"github.com/MichelLacerda/nox/internal/token"
)

type RuntimeError struct {
	Token   *token.Token
	Message string
}

func (r *RuntimeError) Error() string {
	if r.Token != nil {
		return fmt.Sprintf("[line %d] RuntimeError at '%s': %s",
			r.Token.Line, r.Token.Lexeme, r.Message)
	}
	return fmt.Sprintf("RuntimeError: %s", r.Message)
}

func NewRuntimeError(token *token.Token, message string) *RuntimeError {
	return &RuntimeError{
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
