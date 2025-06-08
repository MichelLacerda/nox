package main

import (
	"bufio"
	"fmt"
	"os"
)

type Nox struct {
	hadError bool
}

func (n *Nox) ReportError(line int, where, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
}

func (n *Nox) ErrorAt(line int, message string) {
	n.ReportError(line, "", message)
}

func (n *Nox) RunFile(path string) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	fmt.Println("Running file:", path, " ", len(source), "bytes")

	n.Run(string(source))
	fmt.Println("File execution completed.")
	if n.hadError {
		fmt.Println("Errors encountered during execution.")
		os.Exit(65) // Exit code 65 indicates a runtime error.
	}
	return nil
}

func (n *Nox) RunPrompt() {
	input := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to Nox! Type 'exit' or 'quit' to exit.")
	for {
		fmt.Print(">> ")
		line := input.Scan()
		if !line {
			if input.Err() != nil {
				fmt.Println("Error reading input:", input.Err())
			}
			break
		}
		text := input.Text()
		if text == "exit" || text == "quit" {
			fmt.Println("Exiting Nox.")
			break
		}

		if err := n.Run(text); err != nil {
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
	expr, err := parser.Parse()

	if err != nil {
		return err
	}

	fmt.Println("Parsed expression:", expr)

	for _, token := range tokens {
		fmt.Println(token)
	}

	return nil
}
