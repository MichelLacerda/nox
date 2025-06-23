package runtime

import (
	"fmt"
	"strings"

	"github.com/MichelLacerda/nox/internal/token"
)

type MapInstance struct {
	Entries map[string]any
}

func NewMapInstance(entries map[string]any) *MapInstance {
	return &MapInstance{Entries: entries}
}

func (m *MapInstance) Get(name *token.Token) any {
	val, ok := m.Entries[name.Lexeme]
	if !ok {
		return nil
	}
	return val
}

func (m *MapInstance) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	i := 0
	for k, v := range m.Entries {
		sb.WriteString(fmt.Sprintf("%q: %v", k, v))
		if i < len(m.Entries)-1 {
			sb.WriteString(", ")
		}
		i++
	}
	sb.WriteString("}")
	return sb.String()
}
