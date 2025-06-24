package main

import (
	"fmt"
	"os"

	"github.com/MichelLacerda/nox/internal/runtime"
)

func main() {
	nox := runtime.NewNox()
	if len(os.Args) > 2 {
		fmt.Println("Usage: nox <script>")
		os.Exit(64) // Exit code 64 indicates a command line usage error.
	} else if len(os.Args) == 2 {
		nox.RunFile(os.Args[1])
	} else {
		nox.RunPrompt()
	}
}
