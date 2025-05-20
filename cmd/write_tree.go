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

	tree, err := objects.BuildTreeFromIndex(index)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	// Serializing this tree again is redundant - look for a better way to avoid repeat
	// tree serialization since serialize requires a sort
	serializedTreeData := tree.Serialize()
	repository.WriteObject(tree.Hash.String(), serializedTreeData)
	fmt.Printf("%v\n", tree.Hash)
}

func printWriteTreeUsage() {
	fmt.Println("Usage: gitgood write-tree         Add a file to the staging area (index)")
}
