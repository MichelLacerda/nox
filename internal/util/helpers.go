package util

import "fmt"

func ToFloat(v any) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	default:
		panic(fmt.Sprintf("Invalid type for range argument: %T", val))
	}
}
