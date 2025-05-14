package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CLBRITTON2/go-git-good/common"
	"github.com/CLBRITTON2/go-git-good/objects"
)

func UpdateIndex(flags []string) {
	if len(flags) < 2 {
		printUpdateIndexUsage()
	}
	if flags[0] != "-add" {
		fmt.Println("Unsupported flag...")
		printUpdateIndexUsage()
		return
	}

	// Get absolute path from the file and ensure we have a file to make a blob
	file := flags[1]
	absolutePath, err := filepath.Abs(file)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fileInfo, err := os.Stat(absolutePath)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	blob, err := objects.CreateBlobFromFile(file)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	repository, err := common.FindRepository(".")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	indexRelativePath, err := filepath.Rel(repository.WorkTree, absolutePath)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// Start getting metadata for the index file
	modifiedTime := fileInfo.ModTime()
	fileSize := uint32(fileInfo.Size())
	mode := fileInfo.Mode()
	var fileModeInt uint32
	// Keeping it simple for now - normal files and executable are the only accepted modes
	if mode.IsRegular() {
		if mode&0111 != 0 {
			// Executable
			fileModeInt = 0100755
		} else {
			// Regular file
			fileModeInt = 0100644
		}
	}

	// Just file permissions 644/755
	if fileModeInt != 0100644 && fileModeInt != 0100755 {
		fmt.Println("Unsupported file mode discovered at update-index command")
		return
	}

	indexEntry := &common.IndexEntry{
		ModifiedTime: modifiedTime,
		Hash:         blob.Hash,
		FileSize:     fileSize,
		FileMode:     fileModeInt,
		EntryPath:    indexRelativePath,
	}

	currentIndex, err := common.FindIndex(repository)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	currentIndex.AddEntry(indexEntry)

	err = common.WriteIndex(repository, currentIndex)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}

func printUpdateIndexUsage() {
	fmt.Println("Usage: gitgood update-index -add <filename>         Add a file to the staging area (index)")
}
