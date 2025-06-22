package scanner

import "fmt"

type ScannerError struct {
	Message string
	Line    int
}

func (e ScannerError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.Line, e.Message)
}

func NewScannerError(line int, message string) ScannerError {
	return ScannerError{
		Message: message,
		Line:    line,
	}
}
