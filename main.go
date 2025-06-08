package main

import (
	"fmt"
	"os"
)

var runtime *Nox = NewNox()

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: nox <script>")
		os.Exit(64) // Exit code 64 indicates a command line usage error.
	} else if len(os.Args) == 2 {
		runtime.RunFile(os.Args[1])
	} else {
		runtime.RunPrompt()
	}
}
