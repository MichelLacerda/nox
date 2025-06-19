package main

import (
	"fmt"
	"strings"
)

// Compacto: sem cor, ideal para scripts
func StringifyCompact(value any) string {
	switch v := value.(type) {
	case *ListInstance:
		items := make([]string, len(v.Elements))
		for i, el := range v.Elements {
			items[i] = StringifyCompact(el)
		}
		return "[" + strings.Join(items, ", ") + "]"

	case *DictInstance:
		items := []string{}
		for k, v := range v.Entries {
			items = append(items, fmt.Sprintf("%q: %s", k, StringifyCompact(v)))
		}
		return "{" + strings.Join(items, ", ") + "}"

	default:
		return fmt.Sprintf("%v", v)
	}
}

// Expandido: colorido e indentado para REPL
func StringifyColor(value any, indent string) string {
	switch v := value.(type) {
	case *ListInstance:
		if len(v.Elements) == 0 {
			return "[]"
		}
		builder := strings.Builder{}
		builder.WriteString("[\n")
		for _, el := range v.Elements {
			builder.WriteString(indent + "  " + StringifyColor(el, indent+"  ") + ",\n")
		}
		builder.WriteString(indent + "]")
		return builder.String()

	case *DictInstance:
		if len(v.Entries) == 0 {
			return "{}"
		}
		builder := strings.Builder{}
		builder.WriteString("{\n")
		for k, val := range v.Entries {
			builder.WriteString(fmt.Sprintf("%s  \033[36m%q\033[0m: %s,\n", indent, k, StringifyColor(val, indent+"  ")))
		}
		builder.WriteString(indent + "}")
		return builder.String()

	default:
		// número: amarelo | string: verde
		switch v := value.(type) {
		case float64:
			return fmt.Sprintf("\033[33m%v\033[0m", v)
		case string:
			return fmt.Sprintf("\033[32m%q\033[0m", v)
		default:
			return fmt.Sprintf("%v", v)
		}
	}
}
