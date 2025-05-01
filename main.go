package main

import (
	"os"

	"github.com/CLBRITTON2/go-git-good/cmd"
)

func main() {
	if len(os.Args) < 2 {
		cmd.PrintUsage()
	}

	gitCommand := os.Args[1]
	cmd.Execute(gitCommand)
}
