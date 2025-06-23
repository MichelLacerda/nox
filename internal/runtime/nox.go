package runtime

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/MichelLacerda/nox/internal/parser"
	"github.com/MichelLacerda/nox/internal/scanner"
	"github.com/MichelLacerda/nox/internal/signal"
	"github.com/MichelLacerda/nox/internal/token"
)

type Nox struct {
	HadError        bool
	HadRuntimeError bool
	Interpreter     *Interpreter
	WorkingDir      string         // pasta onde o script foi carregado ou "." no REPL
	Modules         map[string]any // cache de módulos importados
}

func NewNox() *Nox {
	r := &Nox{
		HadError:        false,
		HadRuntimeError: false,
		Interpreter:     nil,
		WorkingDir:      ".",
		Modules:         map[string]any{},
	}
	return r
}

func (n *Nox) ReportRuntimeError(t *token.Token, message string) {
	panic(&RuntimeError{Token: t, Message: message}) // ← usa &
}

func (n *Nox) RunFile(path string) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	if absDir, err := filepath.Abs(filepath.Dir(path)); err == nil {
		n.WorkingDir = absDir
	}

	fmt.Println("Running file:", path, " ", len(source), "bytes")

	interpreter := NewInterpreter(n, false)

	err = n.Run(string(source), interpreter)
	if err != nil {
		// fmt.Printf("Error: %v\n", err)
		switch err := err.(type) {
		case parser.ParserError:
			interpreter.ErrorAt(err.Token.Line, err.Message)
			os.Exit(65)
		case *RuntimeError:
			// já foi exibido em Run()
			os.Exit(70)
		default:
			os.Exit(1)
		}
	} else {
		n.HadError = false
	}

	return nil
}

func (n *Nox) RunPrompt() {
	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		os.Exit(1)
	}
	n.WorkingDir = currentPath // Define o diretório de trabalho como atual
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to Nox! Type 'exit', 'quit' or '\\q' to exit.")
	fmt.Println("Press ENTER twice to execute multiline input.")

	interpreter := NewInterpreter(n, true)

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
			if _, ok := err.(*RuntimeError); !ok {
				// Só imprime erros que NÃO são RuntimeError (ex: ParserError, etc.)
				fmt.Printf("Error: %v\n", err)
			}
		} else {
			n.HadError = false
		}
	}
}

func (n *Nox) Run(source string, interpreter *Interpreter) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(*RuntimeError); ok {
				fmt.Println(runtimeErr.Error()) // só a mensagem bonita
				n.HadRuntimeError = true
				err = runtimeErr // permite tratamento externo se necessário
				return
			} else {
				panic(r)
			}
			// panics inesperados continuam
		}
	}()

	scanner := scanner.NewScanner([]rune(source))
	tokens, err := scanner.ScanTokens()
	if err != nil {
		n.HadError = true
		return fmt.Errorf("failed to scan tokens: %w", err)
	}

	parser := parser.NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		return err
	}

	resolver := NewResolver(interpreter)
	resolver.ResolveStatements(statements)

	if n.HadError {
		return fmt.Errorf("parsing failed with errors")
	}

	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case signal.BreakSignal, signal.ContinueSignal:
				// Silenciar: break/continue usados fora de escopo válido ou fora de loop
				// já tratados no resolver. Aqui são só resíduos que podemos ignorar.
				return
			default:
				panic(r) // repassa qualquer erro real
			}
		}
	}()

	interpreter.Interpret(statements)

	return nil
}
