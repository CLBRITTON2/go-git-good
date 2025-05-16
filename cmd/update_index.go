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
		return
	}
	if flags[0] != "-add" && flags[0] != "-remove" {
		fmt.Println("Unsupported flag...")
		printUpdateIndexUsage()
		return
	}

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

	// Find the repository, find the index, create the index entry path
	// Which is the object's path relative to the work tree
	repository, err := common.FindRepository(".")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	currentIndex, err := common.FindIndex(repository)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	indexEntryRelativePath, err := filepath.Rel(repository.WorkTree, absolutePath)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// We have the entry's path, if the -remove flag was passed we're safe to skip
	// metadata gathering for writing files to the index and just remove the entry
	if flags[0] == "-remove" {
		currentIndex.RemoveEntry(indexEntryRelativePath)
		err = common.WriteIndex(repository, currentIndex)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		return
	}

	blob, err := objects.CreateBlobFromFile(file)
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
		EntryPath:    indexEntryRelativePath,
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
	fmt.Println("Usage: gitgood update-index -remove <filename>      Remove a file from the staging area (index)")
}
