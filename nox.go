package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Nox struct {
	hadError        bool
	hadRuntimeError bool
	interpreter     *Interpreter
}

func NewNox() *Nox {
	r := &Nox{
		hadError:        false,
		hadRuntimeError: false,
	}
	return r
}

type ParserError struct {
	Token   *Token
	Message string
}

func (e ParserError) Error() string {
	return fmt.Sprintf("Parser Error at %s: %s", e.Token.Lexeme, e.Message)
}

type RuntimeError struct {
	Token   *Token
	Message string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("Runtime Error at %s: %s", e.Token.Lexeme, e.Message)
}

func NewRuntimeError(token *Token, message string) RuntimeError {
	return RuntimeError{
		Token:   token,
		Message: message,
	}
}

func (n *Nox) ReportError(line int, where, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
}

func (n *Nox) ErrorAt(line int, message string) {
	n.ReportError(line, "", message)
}

func (n *Nox) ReportRuntimeError(t *Token, message string) {
	panic(RuntimeError{Token: t, Message: message})
}

func (n *Nox) RunFile(path string) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	fmt.Println("Running file:", path, " ", len(source), "bytes")

	interpreter := NewInterpreter(n, StringifyCompact)

	err = n.Run(string(source), interpreter)
	if err != nil {
		// fmt.Printf("Error: %v\n", err)
		switch err := err.(type) {
		case ParserError:
			n.ErrorAt(err.Token.line, err.Message)
			os.Exit(65)
		case RuntimeError:
			// já foi exibido em Run()
			os.Exit(70)
		default:
			os.Exit(1)
		}
	} else {
		n.hadError = false
	}

	return nil
}

func (n *Nox) RunPrompt() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to Nox! Type 'exit', 'quit' or '\\q' to exit.")
	fmt.Println("Press ENTER twice to execute multiline input.")

	interpreter := NewInterpreter(n, func(val any) string {
		return StringifyColor(val, "")
	})

	for {
		var lines []string
		for {
			fmt.Print(">> ")
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Println("Exiting Nox.")
					return
				}
				fmt.Println("Error reading input:", err)
				return
			}

			text := strings.TrimSpace(line)

			if text == "exit" || text == "quit" || text == "\\q" {
				fmt.Println("Exiting Nox.")
				return
			}

			// Break if the user hits enter twice
			if text == "" {
				break
			}

			lines = append(lines, line)
		}

		src := strings.Join(lines, "")
		if err := n.Run(src, interpreter); err != nil {
			if _, ok := err.(RuntimeError); !ok {
				// Só imprime erros que NÃO são RuntimeError (ex: ParserError, etc.)
				fmt.Printf("Error: %v\n", err)
			}
		} else {
			n.hadError = false
		}
	}
}

func (n *Nox) Run(source string, interpreter *Interpreter) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(RuntimeError); ok {
				fmt.Println(runtimeErr.Error()) // só a mensagem bonita
				n.hadRuntimeError = true
				err = runtimeErr // permite tratamento externo se necessário
				return
			}
			// panics inesperados continuam
			panic(r)
		}
	}()

	scanner := NewScanner([]rune(source))
	tokens := scanner.ScanTokens()

	parser := NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		return err
	}

	resolver := NewResolver(interpreter)
	resolver.ResolveStatements(statements)

	if n.hadError {
		return fmt.Errorf("parsing failed with errors")
	}

	interpreter.Interpret(statements)

	return nil
}
