package cmd

import (
	"fmt"

	"github.com/CLBRITTON2/go-git-good/common"
	"github.com/CLBRITTON2/go-git-good/objects"
)

func Log() {
	repository, err := common.FindRepository(".")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
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
	if ref.Hash.Empty() {
		fmt.Printf("fatal: your current branch '%s' does not have any commits yet\n", branch)
		return
	}

	// Start recursive commit printing
	printCommitHistory(repository, ref, branch, true)
}

func printCommitHistory(repository *common.Repository, ref *common.Ref, branch string, isHead bool) {
	rawCommitData, err := repository.ReadObject(ref.Hash.String())
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	commit, err := objects.ParseCommit(rawCommitData)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	// Format the commit header: include (HEAD -> branch) only for the HEAD commit
	commitHeader := fmt.Sprintf("commit %s", ref.Hash.String())
	if isHead {
		commitHeader += fmt.Sprintf(" (HEAD -> %s)", branch)
	}
	fmt.Printf("%s\nAuthor: %s\nDate: %v\n\n    %s\n\n", commitHeader, commit.Author, commit.Timestamp.Format("Mon Jan 02 15:04:05 2006 -0700"), commit.Message)

	// Recursively process the first parent (if any)
	if len(commit.Parents) > 0 {
		parentRef := &common.Ref{
			Name: branch,
			Hash: commit.Parents[0],
		}
		printCommitHistory(repository, parentRef, branch, false)
	}
}
