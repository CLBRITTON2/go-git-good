package cmd

import (
	"fmt"
	"os"
)

func Execute(argument string) {
	switch argument {
	case "init":
		fmt.Println("Git init placeholder...")
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
