package cmd

import (
	"fmt"
	"io/fs"
	"os"

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
		fileSystem := os.DirFS(root)
		err := fs.WalkDir(fileSystem, ".", processEntry)
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

func processEntry(path string, entry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if entry.IsDir() && (entry.Name() == ".gitgood" || entry.Name() == ".git") {
		return fs.SkipDir
	}

	if !entry.IsDir() {
		if entry.Type().IsRegular() {
			fmt.Printf("Found regular file: %v\n", path)
			UpdateIndex([]string{"-add", path})
		}
	}
	return nil
}

func printAddUsage() {
	fmt.Println("Usage: gitgood add <filename>         Stage a file by adding it to the index")
	fmt.Println("Usage: gitgood add .                  Stage all files in the working directory")
}
