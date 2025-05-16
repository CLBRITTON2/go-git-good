package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/CLBRITTON2/go-git-good/common"
)

func Add(flags []string) {
	// Will only support single file or the entire work tree initially
	if len(flags) != 1 {
		printAddUsage()
		return
	}

	repository, err := common.FindRepository(".")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// Walk the entire work tree and all sub directories, add to the index
	if flags[0] == "." {
		root := repository.WorkTree
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			// Use repository.WorkTree as the base when determining paths
			return processEntry(repository.WorkTree, path, entry, err)
		})
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		// processEntry will add all discovered files to the index
		// Return here to allow single file handling
		return
	}

	UpdateIndex([]string{"-add", flags[0]})
}

func processEntry(rootPath string, path string, entry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if entry.IsDir() && (entry.Name() == ".gitgood" || entry.Name() == ".git") {
		return fs.SkipDir
	}

	if !entry.IsDir() && entry.Type().IsRegular() {
		// Calculate relative path from the repository root
		relativePath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return fmt.Errorf("error getting relative path %s: %v", path, err)
		}

		UpdateIndex([]string{"-add", relativePath})
	}

	return nil
}

func printAddUsage() {
	fmt.Println("Usage: gitgood add <filename>         Stage a file by adding it to the index")
	fmt.Println("Usage: gitgood add .                  Stage all files in the working directory")
}
