package cmd

import (
	"fmt"
	"os"
)

func Execute(command string, flags []string) {
	switch command {
	case "init":
		Init(flags)
	case "hash-object":
		HashObject(flags)
	case "cat-file":
		CatFile(flags)
	case "add":
		fmt.Println("Git add placeholder...")
	case "commit":
		fmt.Println("Git commit placeholder...")
	}
}

func PrintUsage() {
	fmt.Println("Go-Git-Good Usage: gitgood <command> <args>")
	fmt.Println("Commands:")
	fmt.Println("init          Create an empty GGG repository")
	fmt.Println("add           Add file contents to the repository")
	fmt.Println("commit        Save changes to the repository")
	os.Exit(1)
}
