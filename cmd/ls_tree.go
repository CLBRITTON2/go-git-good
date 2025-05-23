package cmd

import (
	"bytes"
	"fmt"
	"strings"

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

	nullIndex := bytes.IndexByte(rawObjectData, byte('\x00'))
	if nullIndex == -1 {
		fmt.Printf("invalid object format: no null byte found")
		return
	}

	header := string(rawObjectData[:nullIndex])
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		fmt.Printf("invalid object header format expected <type> <data length> got: %s", header)
		return
	}

	objectType := parts[0]
	var tree *objects.Tree
	switch objectType {
	case "tree":
		tree, err = objects.ParseTree(rawObjectData)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	case "commit":
		commit, err := objects.ParseCommit(rawObjectData)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		tree = commit.Tree
		rawTreeData, err := repository.ReadObject(tree.Hash.String())
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		tree, err = objects.ParseTree(rawTreeData)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
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
