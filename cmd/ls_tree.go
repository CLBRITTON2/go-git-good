package cmd

import (
	"fmt"

	"github.com/CLBRITTON2/go-git-good/common"
	"github.com/CLBRITTON2/go-git-good/objects"
)

func LsTree(flags []string) {
	if len(flags) != 1 {
		printLsTreeUsage()
		return
	}
	objectHash := flags[0]
	if len(objectHash) != 40 {
		fmt.Printf("invalid hash length: expected length 40 got %v\n", len(objectHash))
		return
	}
	repository, err := common.FindRepository(".")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	rawObjectData, err := repository.ReadObject(objectHash)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	tree, err := objects.ParseTree(rawObjectData)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	// Create the same format that Git uses for ls-tree and cat-file -p with a tree hash
	for _, entry := range tree.Entries {
		objectType := ""
		if entry.FileMode == 040000 {
			objectType = "tree"
			fmt.Printf("%#o %s %s   %s\n", entry.FileMode, objectType, entry.Hash.String(), entry.Name)
		} else {
			objectType = "blob"
			fmt.Printf("%o %s %s   %s\n", entry.FileMode, objectType, entry.Hash.String(), entry.Name)
		}
	}
}

func printLsTreeUsage() {
	fmt.Println("Usage: gitgood ls-tree <tree-hash>          Print the tree contents")
}
