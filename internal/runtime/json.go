package runtime

import (
	"encoding/json"
)

func NewJsonModule() *MapInstance {
	return NewMapInstance(map[string]any{
		"encode": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				// Converter ListInstance ou DictInstance em dados Go nativos
				data := toGoValue(args[0])
				jsonBytes, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					i.Runtime.ReportRuntimeError(nil, "json.encode: "+err.Error())
					return nil
				}
				return string(jsonBytes)
			},
		},
		"decode": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				input, ok := args[0].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "json.decode expects a string")
					return nil
				}

				var result any
				if err := json.Unmarshal([]byte(input), &result); err != nil {
					i.Runtime.ReportRuntimeError(nil, "json.decode: "+err.Error())
					return nil
				}

				return fromGoValue(result)
			},
		},
	})
}

func toGoValue(v any) any {
	switch val := v.(type) {
	case *ListInstance:
		var list []any
		for _, item := range val.Elements {
			list = append(list, toGoValue(item))
		}
		return list
	case *DictInstance:
		m := map[string]any{}
		for k, v := range val.Entries {
			m[k] = toGoValue(v)
		}
		return m
	default:
		return val
	}
}

func fromGoValue(v any) any {
	switch val := v.(type) {
	case []any:
		var list []any
		for _, item := range val {
			list = append(list, fromGoValue(item))
		}
		return NewListInstance(list)
	case map[string]any:
		m := map[string]any{}
		for k, v := range val {
			m[k] = fromGoValue(v)
		}
		return NewDictInstance(m)
	default:
		return val
	}
}
