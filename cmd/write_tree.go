package cmd

import (
	"fmt"

	"github.com/CLBRITTON2/go-git-good/common"
	"github.com/CLBRITTON2/go-git-good/objects"
)

func WriteTree(flags []string) {
	if len(flags) != 0 {
		printWriteTreeUsage()
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
	fmt.Printf("%v\n", rootTree.Hash)
}

func printWriteTreeUsage() {
	fmt.Println("Usage: gitgood write-tree         Creates a tree object from the current index and writes it to the DB")
}
