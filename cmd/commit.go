package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/CLBRITTON2/go-git-good/common"
	"github.com/CLBRITTON2/go-git-good/objects"
)

func Commit(flags []string) {
	if len(flags) < 1 || flags[0] != "-m" {
		printCommitUsage()
		return
	}
	message := flags[1]

	repository, err := common.FindRepository(".")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	index, err := common.GetIndex(repository)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	exists, err := index.Exists(repository)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	if !exists {
		fmt.Println("No files have been staged to commit. Use update-index or add to stage files.")
		return
	}

	rootTree, _, err := objects.BuildTreeFromIndex(index)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	// Ensure the tree is written to the object DB
	WriteTree([]string{"-q"})
	author, email := parseGitConfig()
	if author == "" {
		author = "local user"
	}
	authorString := fmt.Sprintf("%s <%s>", author, email)
	commit := &objects.Commit{
		Tree:      rootTree,
		Parents:   []common.Hash{},
		Author:    authorString,
		Timestamp: time.Now(),
		Message:   message,
	}

	branch, err := repository.GetBranch()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	// Empty ref is returned if this is the first commit
	ref, err := repository.FindRef(branch)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	if !ref.Hash.Empty() {
		commit.Parents = []common.Hash{ref.Hash}
	}

	serializedCommitData := commit.Serialize()
	commitHash, err := common.HashObject(serializedCommitData)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	ref.Hash = commitHash
	repository.WriteRef(ref, branch)
	repository.WriteObject(commitHash.String(), serializedCommitData)
}

// This might need to be moved to a more accessible location
func parseGitConfig() (name, email string) {
	userHomeDirectory, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error locating user home directory: %v\n", err)
		return
	}

	configLocations := []string{
		filepath.Join(userHomeDirectory, ".gitconfig"),
		filepath.Join(userHomeDirectory, ".config", "git", "config"),
	}

	var fileContent []byte
	found := false

	for _, config := range configLocations {
		fileInfo, err := os.ReadFile(config)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			fmt.Printf("error reading config file %s: %v\n", config, err)
			continue
		}
		fileContent = fileInfo
		found = true
		break
	}

	if !found {
		return "", ""
	}

	nameRegex := regexp.MustCompile(`name\s*=\s*(.+)`)
	emailRegex := regexp.MustCompile(`email\s*=\s*(.+)`)

	nameMatch := nameRegex.FindStringSubmatch(string(fileContent))
	if len(nameMatch) > 1 {
		name = strings.TrimSpace(nameMatch[1])
	}

	emailMatch := emailRegex.FindStringSubmatch(string(fileContent))
	if len(emailMatch) > 1 {
		email = strings.TrimSpace(emailMatch[1])
	}

	return name, email
}

func printCommitUsage() {
	fmt.Println("Usage: gitgood commit -m <message>          Record changes to the repository")
}
