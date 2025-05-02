package cmd

import (
	"fmt"
	"os"

	"github.com/CLBRITTON2/go-git-good/internal"
)

func Execute(argument string, path string) {
	switch argument {
	case "init":
		_, err := internal.CreateRepository(path)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	case "add":
		fmt.Println("Git add placeholder...")
	case "commit":
		fmt.Println("Git commit placeholder...")
	}
}

func PrintUsage() {
	fmt.Println("Go-Git-Good Usage: ggg <command> <args>")
	fmt.Println("Commands:")
	fmt.Println("init          Create an empty GGG repository")
	fmt.Println("add           Add file contents to the repository")
	fmt.Println("commit        Save changes to the repository")
	os.Exit(1)
}
