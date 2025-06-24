package util

import (
	"fmt"
	"os"
)

// ParseFileMode converts a string representation of a file mode to an os file mode.
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

// PathExists checks if a file or directory exists at the given path.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
