package main

import (
	"fmt"
	"os"
)

var Interpret *Nox = &Nox{}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: nox <script>")
		os.Exit(64) // Exit code 64 indicates a command line usage error.
	} else if len(os.Args) == 2 {
		Interpret.RunFile(os.Args[1])
	} else {
		Interpret.RunPrompt()
	}
}
