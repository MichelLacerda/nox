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

	r.interpreter = NewInterpreter(r)
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
	return e.Message
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
	fmt.Printf("Runtime Error at %s: %s line %d\n", t.Lexeme, message, t.line)
	n.hadRuntimeError = true
}

func (n *Nox) RunFile(path string) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	fmt.Println("Running file:", path, " ", len(source), "bytes")

	if err := n.Run(string(source)); err != nil {
		fmt.Printf("Error: %v\n", err)
		if perr, ok := err.(ParserError); ok {
			n.ErrorAt(perr.Token.line, perr.Message)
			os.Exit(65)
		} else if rerr, ok := err.(RuntimeError); ok {
			n.ReportRuntimeError(rerr.Token, rerr.Message)
			os.Exit(70)
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
		if err := n.Run(src); err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			n.hadError = false
		}
	}
}

func (n *Nox) Run(source string) error {
	scanner := NewScanner([]rune(source))
	tokens := scanner.ScanTokens()

	parser := NewParser(tokens)
	statements, err := parser.Parse()

	if err != nil {
		return err
	}

	resolver := NewResolver(n.interpreter)
	resolver.ResolveStatements(statements)

	if n.hadError {
		return fmt.Errorf("parsing failed with errors")
	}

	n.interpreter.Interpret(statements)

	return nil
}
