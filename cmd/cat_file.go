package cmd

import (
	"fmt"

	"github.com/CLBRITTON2/go-git-good/objects"
)

func CatFile(flags []string) {
	if len(flags) != 1 {
		PrintCatFileUsage()
		return
	}
	objectHash := flags[0]
	if len(objectHash) != 40 {
		fmt.Printf("invalid hash length\n")
		return
	}
	repository, err := objects.FindRepository(".")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fileContents, err := repository.ReadObject(objectHash)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%v", fileContents)
}

func PrintCatFileUsage() {
	fmt.Println("Usage: gitgood cat-file <object-hash>          Print the file contents")
}
