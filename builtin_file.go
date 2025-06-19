package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type FileObject struct {
	file   *os.File
	reader *bufio.Reader
}

var _ HasMethods = (*FileObject)(nil)

func (f *FileObject) TypeName() string {
	return "File"
}

func (f *FileObject) String() string {
	return fmt.Sprintf("<file %s>", f.file.Name())
}

func (f *FileObject) GetMethod(name string) any {
	//func (f *FileObject) GetMethod(name string) Callable {
	switch name {
	case "read":
		return &BuiltinFunction{
			arity: 0,
			call: func(i *Interpreter, args []any) any {
				data, err := io.ReadAll(f.file)
				if err != nil {
					i.runtime.ReportRuntimeError(&Token{Lexeme: "read"}, "read error: "+err.Error())
					return nil
				}
				return string(data)
			},
		}
	case "read_bytes":
		return &BuiltinFunction{
			arity: 0,
			call: func(i *Interpreter, args []any) any {
				data, err := io.ReadAll(f.file)
				if err != nil {
					i.runtime.ReportRuntimeError(&Token{Lexeme: "read_bytes"}, "read_bytes error: "+err.Error())
					return nil
				}
				return data
			},
		}
	case "readline":
		return &BuiltinFunction{
			arity: 0,
			call: func(i *Interpreter, args []any) any {
				if f.reader == nil {
					i.runtime.ReportRuntimeError(&Token{Lexeme: "readline"}, "File not readable.")
					return nil
				}

				line, err := f.reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						if line == "" {
							return nil // Fim de arquivo
						}
						// Última linha sem \n
						return strings.TrimRight(line, "\n\r")
					}
					i.runtime.ReportRuntimeError(&Token{Lexeme: "readline"}, "readline error: "+err.Error())
					return nil
				}

				return strings.TrimRight(line, "\n\r")
			},
		}

	case "write":
		return &BuiltinFunction{
			arity: 1,
			call: func(i *Interpreter, args []any) any {
				str, ok := args[0].(string)
				if !ok {
					i.runtime.ReportRuntimeError(&Token{Lexeme: "write"}, "write() expects a string")
					return nil
				}
				_, err := f.file.WriteString(str)
				if err != nil {
					i.runtime.ReportRuntimeError(&Token{Lexeme: "write"}, "write error: "+err.Error())
					return nil
				}
				return nil
			},
		}
	case "write_bytes":
		return &BuiltinFunction{
			arity: 1,
			call: func(i *Interpreter, args []any) any {
				str, ok := args[0].(string)
				if !ok {
					i.runtime.ReportRuntimeError(&Token{Lexeme: "write_bytes"}, "write_bytes() expects a string")
					return nil
				}
				bytes := []byte(str)
				if !ok {
					i.runtime.ReportRuntimeError(&Token{Lexeme: "write_bytes"}, "write_bytes() expects a byte slice")
					return nil
				}
				_, err := f.file.Write(bytes)
				if err != nil {
					i.runtime.ReportRuntimeError(&Token{Lexeme: "write_bytes"}, "write_bytes error: "+err.Error())
					return nil
				}
				return nil
			},
		}
	case "flush":
		return &BuiltinFunction{
			arity: 0,
			call: func(i *Interpreter, args []any) any {
				err := f.file.Sync()
				if err != nil {
					i.runtime.ReportRuntimeError(&Token{Lexeme: "flush"}, "flush error: "+err.Error())
					return nil
				}
				return nil
			},
		}
	case "tell":
		return &BuiltinFunction{
			arity: 0,
			call: func(i *Interpreter, args []any) any {
				pos, err := f.file.Seek(0, io.SeekCurrent)
				if err != nil {
					i.runtime.ReportRuntimeError(&Token{Lexeme: "tell"}, "tell error: "+err.Error())
					return nil
				}
				return pos
			},
		}
	case "seek":
		return &BuiltinFunction{
			arity: 2,
			call: func(i *Interpreter, args []any) any {
				offset, ok1 := args[0].(float64)
				whence, ok2 := args[1].(float64)
				if !ok1 || !ok2 {
					i.runtime.ReportRuntimeError(&Token{Lexeme: "seek"}, "seek(offset, whence) expects numbers")
					return nil
				}
				pos, err := f.file.Seek(int64(offset), int(whence))
				if err != nil {
					i.runtime.ReportRuntimeError(&Token{Lexeme: "seek"}, "seek error: "+err.Error())
					return nil
				}
				return pos
			},
		}
	case "exists":
		return &BuiltinFunction{
			arity: 0,
			call: func(i *Interpreter, args []any) any {
				_, err := os.Stat(f.file.Name())
				if os.IsNotExist(err) {
					return false
				} else if err != nil {
					i.runtime.ReportRuntimeError(&Token{Lexeme: "exists"}, "exists error: "+err.Error())
					return nil
				}
				return true
			},
		}
	case "close":
		return &BuiltinFunction{
			arity: 0,
			call: func(i *Interpreter, args []any) any {
				err := f.file.Close()
				if err != nil {
					i.runtime.ReportRuntimeError(&Token{Lexeme: "close"}, "close error: "+err.Error())
				}
				return nil
			},
		}
	default:
		return nil
	}
}

func parseFileMode(mode string) (int, error) {
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
