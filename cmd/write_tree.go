package cmd

import (
	"fmt"

	"github.com/CLBRITTON2/go-git-good/common"
	"github.com/CLBRITTON2/go-git-good/objects"
)

func WriteTree(flags []string) {
	// This is just going to be an internal flag to avoid printing the root tree hash when committing
	// It would probably be better to just return the root tree hash from this command and print it
	// One level up but that would require restructuring all commands to function the same... TODO maybe?
	quiet := false
	if len(flags) == 1 && flags[0] == "-q" {
		quiet = true
	} else if len(flags) > 0 {
		printWriteTreeUsage()
		return
	}

	repository, err := common.FindRepository(".")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	index, err := common.GetIndex(repository)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	exists, err := index.Exists(repository)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	if !exists {
		fmt.Println("No files have been staged in the index. Use update-index or add to stage files to write a tree.")
		return
	}

	rootTree, trees, err := objects.BuildTreeFromIndex(index)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	for _, tree := range trees {
		serializedTreeData := tree.Serialize()
		err = repository.WriteObject(tree.Hash.String(), serializedTreeData)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}
	if !quiet {
		fmt.Printf("%v\n", rootTree.Hash)
	}
}

func printWriteTreeUsage() {
	fmt.Println("Usage: gitgood write-tree         Creates a tree object from the current index and writes it to the DB")
}
