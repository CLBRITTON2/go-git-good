package cmd

import (
	"fmt"

	"github.com/CLBRITTON2/go-git-good/common"
)

func LsFiles(flags []string) {
	if len(flags) > 1 || (len(flags) == 1 && flags[0] != "-s") {
		fmt.Println("Unsupported flag...")
		printLsFilesUsage()
		return
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

	if flags[0] == "-s" {
		for _, entry := range index.Entries {
			fmt.Printf("%o %v         %v\n", entry.FileMode, entry.Hash.String(), entry.EntryPath)
			return
		}

		for _, entry := range index.Entries {
			fmt.Printf("%v\n", entry.EntryPath)
		}
	}
}

func printLsFilesUsage() {
	fmt.Println("Usage: gitgood ls-files          Show information about files in the index")
	fmt.Println("Usage: gitgood ls-files -s       Show staged contents' mode bits, object hash, and index entry path")
}
