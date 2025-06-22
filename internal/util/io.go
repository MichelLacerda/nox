package util

import (
	"fmt"
	"os"
)

func ParseFileMode(mode string) (int, error) {
	switch mode {
	case "r":
		return os.O_RDONLY, nil
	case "w":
		return os.O_WRONLY | os.O_CREATE | os.O_TRUNC, nil
	case "a":
		return os.O_WRONLY | os.O_CREATE | os.O_APPEND, nil
	case "r+":
		return os.O_RDWR, nil
	case "w+":
		return os.O_RDWR | os.O_CREATE | os.O_TRUNC, nil
	case "a+":
		return os.O_RDWR | os.O_CREATE | os.O_APPEND, nil
	case "rb":
		return os.O_RDONLY, nil
	case "wb":
		return os.O_WRONLY | os.O_CREATE | os.O_TRUNC, nil
	case "ab":
		return os.O_WRONLY | os.O_CREATE | os.O_APPEND, nil
	case "rb+", "r+b":
		return os.O_RDWR, nil
	case "wb+", "w+b":
		return os.O_RDWR | os.O_CREATE | os.O_TRUNC, nil
	case "ab+", "a+b":
		return os.O_RDWR | os.O_CREATE | os.O_APPEND, nil
	default:
		return 0, fmt.Errorf("unsupported file mode: %s", mode)
	}
}
