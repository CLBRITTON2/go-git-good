package cmd

import (
	"fmt"

	"github.com/CLBRITTON2/go-git-good/objects"
)

func HashObject(flags []string) {
	// Hash-object requires a file for now at minimum
	// It also won't do anything more than print or write to the DB so more/less flags = bad
	if len(flags) > 2 || len(flags) < 1 {
		printHashObjectUsage()
	}

	write := false
	file := ""

	if flags[0] == "-w" {
		if len(flags) != 2 {
			printHashObjectUsage()
			return
		}
		write = true
		file = flags[1]
	} else {
		file = flags[0]
	}

	blob, err := objects.CreateBlobFromFile(file)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	serializedBlobData := blob.Serialize()
	blobHash := objects.CalculateHash(serializedBlobData)

	if !write {
		fmt.Printf("%v\n", blobHash)
		return
	}

	repository, err := objects.FindRepository(".")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	err = repository.WriteObject(blobHash, serializedBlobData)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%v\n", blobHash)
}

func printHashObjectUsage() {
	fmt.Println("Usage: gitgood hash-object <file>          Print the SHA1 hash")
	fmt.Println("Usage: gitgood hash-object -w <file>       Write the blob to the git DB")
}
