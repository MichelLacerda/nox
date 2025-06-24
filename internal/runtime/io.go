package runtime

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/MichelLacerda/nox/internal/token"
)

type FileObject struct {
	File   *os.File
	Reader *bufio.Reader
}

var _ HasMethods = (*FileObject)(nil)

func (f *FileObject) TypeName() string {
	return "File"
}

func (f *FileObject) String() string {
	return fmt.Sprintf("<file %s>", f.File.Name())
}

func (f *FileObject) GetMethod(name string) any {
	switch name {
	case "read":
		return &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(i *Interpreter, args []any) any {
				data, err := io.ReadAll(f.File)
				if err != nil {
					i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "read"}, "read error: "+err.Error())
					return nil
				}
				return string(data)
			},
		}
	case "read_bytes":
		return &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(i *Interpreter, args []any) any {
				data, err := io.ReadAll(f.File)
				if err != nil {
					i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "read_bytes"}, "read_bytes error: "+err.Error())
					return nil
				}
				return data
			},
		}
	case "readline":
		return &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(i *Interpreter, args []any) any {
				if f.Reader == nil {
					i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "readline"}, "File not readable.")
					return nil
				}

				line, err := f.Reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						if line == "" {
							return nil // Fim de arquivo
						}
						// Última linha sem \n
						return strings.TrimRight(line, "\n\r")
					}
					i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "readline"}, "readline error: "+err.Error())
					return nil
				}

				return strings.TrimRight(line, "\n\r")
			},
		}

	case "write":
		return &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				str, ok := args[0].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "write"}, "write() expects a string")
					return nil
				}
				_, err := f.File.WriteString(str)
				if err != nil {
					i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "write"}, "write error: "+err.Error())
					return nil
				}
				return nil
			},
		}
	case "write_bytes":
		return &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				str, ok := args[0].(string)
				if !ok {
					i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "write_bytes"}, "write_bytes() expects a string")
					return nil
				}
				bytes := []byte(str)
				if !ok {
					i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "write_bytes"}, "write_bytes() expects a byte slice")
					return nil
				}
				_, err := f.File.Write(bytes)
				if err != nil {
					i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "write_bytes"}, "write_bytes error: "+err.Error())
					return nil
				}
				return nil
			},
		}
	case "flush":
		return &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(i *Interpreter, args []any) any {
				err := f.File.Sync()
				if err != nil {
					i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "flush"}, "flush error: "+err.Error())
					return nil
				}
				return nil
			},
		}
	case "tell":
		return &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(i *Interpreter, args []any) any {
				pos, err := f.File.Seek(0, io.SeekCurrent)
				if err != nil {
					i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "tell"}, "tell error: "+err.Error())
					return nil
				}
				return pos
			},
		}
	case "seek":
		return &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(i *Interpreter, args []any) any {
				offset, ok1 := args[0].(float64)
				whence, ok2 := args[1].(float64)
				if !ok1 || !ok2 {
					i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "seek"}, "seek(offset, whence) expects numbers")
					return nil
				}
				pos, err := f.File.Seek(int64(offset), int(whence))
				if err != nil {
					i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "seek"}, "seek error: "+err.Error())
					return nil
				}
				return pos
			},
		}
	case "exists":
		return &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(i *Interpreter, args []any) any {
				_, err := os.Stat(f.File.Name())
				if os.IsNotExist(err) {
					return false
				} else if err != nil {
					i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "exists"}, "exists error: "+err.Error())
					return nil
				}
				return true
			},
		}
	case "close":
		return &BuiltinFunction{
			ArityValue: 0,
			CallFunc: func(i *Interpreter, args []any) any {
				err := f.File.Close()
				if err != nil {
					i.Runtime.ReportRuntimeError(&token.Token{Lexeme: "close"}, "close error: "+err.Error())
				}
				return nil
			},
		}
	default:
		return nil
	}
}
