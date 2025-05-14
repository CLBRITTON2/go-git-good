package cmd

import (
	"fmt"

	"github.com/CLBRITTON2/go-git-good/common"
)

func LsFiles(flags []string) {
	if len(flags) > 0 {
		fmt.Println("Unsupported flag...")
		printLsFilesUsage()
	}

	repository, err := common.FindRepository(".")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	index, err := common.FindIndex(repository)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	for _, entry := range index.Entries {
		fmt.Printf("%v\n", entry.EntryPath)
	}
}

func printLsFilesUsage() {
	fmt.Println("Usage: gitgood ls-files          Show information about files in the index")
}
