package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/CLBRITTON2/go-git-good/common"
)

// Store repository root for directory walk
// Not sure if this is the best way to do this...
var currentRepoRoot string

func Add(flags []string) {
	if len(flags) != 1 {
		printAddUsage()
		return
	}

	repository, err := common.FindRepository(".")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	flagIsDirectory, err := isDir(flags[0])
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// Walk the entire work tree and all sub directories, add to the index
	if flags[0] == "." {
		// Set the global repository root before walking
		currentRepoRoot = repository.WorkTree

		err := filepath.WalkDir(currentRepoRoot, processEntry)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		return
	}

	// Only walk the specified directory
	if flagIsDirectory {
		err := filepath.WalkDir(flags[0], processEntry)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		return
	}

	UpdateIndex([]string{"-add", flags[0]})
}

func processEntry(path string, entry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if entry.IsDir() && (entry.Name() == ".gitgood" || entry.Name() == ".git") {
		return fs.SkipDir
	}

	if !entry.IsDir() && entry.Type().IsRegular() {
		// Calculate relative path from the repository root
		relativePath, err := filepath.Rel(currentRepoRoot, path)
		if err != nil {
			return fmt.Errorf("error getting relative path %s: %v", path, err)
		}

		UpdateIndex([]string{"-add", relativePath})
	}

	return nil
}

func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func printAddUsage() {
	fmt.Println("Usage: gitgood add <filename>         Stage a file by adding it to the index")
	fmt.Println("Usage: gitgood add .                  Stage all files in the working directory")
}
