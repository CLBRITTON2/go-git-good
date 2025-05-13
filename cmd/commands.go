package cmd

import (
	"fmt"
)

func Execute(command string, flags []string) {
	switch command {
	case "init":
		Init(flags)
	case "hash-object":
		HashObject(flags)
	case "cat-file":
		CatFile(flags)
	case "update-index":
		UpdateIndex(flags)
	case "add":
		fmt.Println("Git add placeholder...")
	case "commit":
		fmt.Println("Git commit placeholder...")
	default:
		fmt.Println("Unsupported command...")
		PrintUsage()
	}
}

func PrintUsage() {
	fmt.Println("Go-Git-Good Usage: gitgood <command> <args>")
	fmt.Println("Commands:")
	fmt.Println("init          Create an empty gitgood repository")
	fmt.Println("hash-object   Compute object ID and optionally write an object to the DB")
	fmt.Println("cat-file      Print the contents of an object in the DB")
	fmt.Println("add           Add file contents to the repository")
	fmt.Println("commit        Save changes to the repository")
}
