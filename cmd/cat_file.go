package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/CLBRITTON2/go-git-good/common"
	"github.com/CLBRITTON2/go-git-good/objects"
)

func CatFile(flags []string) {
	if len(flags) != 1 {
		printCatFileUsage()
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
	switch objectType {
	case "blob":
		content := string(rawObjectData[nullIndex+1:])
		fmt.Printf("%v", content)
	case "tree":
		LsTree(flags)
	case "commit":
		content := rawObjectData[nullIndex:]
		commit, err := objects.ParseCommit(content)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		utcOffset := fmt.Sprintf("%d %s", commit.Timestamp.Unix(), commit.Timestamp.Format("-0700"))
		commitString := fmt.Sprintf("tree %s\nauthor %s %s\ncommitter %v %s\n", commit.Tree.Hash.String(), commit.Author, utcOffset, commit.Author, utcOffset)
		fmt.Println(commitString)

	}
}

func printCatFileUsage() {
	fmt.Println("Usage: gitgood cat-file <object-hash>          Print the file contents")
}
