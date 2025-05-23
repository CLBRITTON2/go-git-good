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
	case "ls-files":
		LsFiles(flags)
	case "add":
		Add(flags)
	case "write-tree":
		WriteTree(flags)
	case "ls-tree":
		LsTree(flags)
	case "commit":
		Commit(flags)
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
	fmt.Println("update-index  Register file contents in the working tree to the index")
	fmt.Println("ls-files      Show information about files in the index")
	fmt.Println("add           Add file contents to the index and DB")
	fmt.Println("write-tree    Create a tree object from the current index and write it to the DB")
	fmt.Println("ls-tree       List the contents of a tree object")
}
