package objects

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/CLBRITTON2/go-git-good/common"
)

type Commit struct {
	Tree      *Tree
	Parents   []common.Hash
	Author    string
	Timestamp time.Time
	Message   string
}

// Ref https://stackoverflow.com/questions/22968856/what-is-the-file-format-of-a-git-commit-object-data-structure
func (commit *Commit) Serialize() []byte {
	// Generate the content string so we can pull the size for the top level commit header [commit {size}\0{content}]
	content := fmt.Sprintf("tree %v\n", commit.Tree.Hash)
	for _, parent := range commit.Parents {
		if !parent.Empty() {
			content += fmt.Sprintf("parent %v\n", parent)
		}
	}
	utcOffset := fmt.Sprintf("%d %s", commit.Timestamp.Unix(), commit.Timestamp.Format("-0700"))
	// Committer will be the author for simplicity
	content += fmt.Sprintf("author %s %s\n", commit.Author, utcOffset)
	content += fmt.Sprintf("committer %s %s\n\n", commit.Author, utcOffset)
	content += commit.Message

	header := fmt.Sprintf("commit %d\x00", len(content))
	return append([]byte(header), []byte(content)...)
}

func ParseCommit(rawCommitData []byte) (*Commit, error) {
	// Find the null byte that separates header from content
	nullIndex := bytes.IndexByte(rawCommitData, byte('\x00'))
	if nullIndex == -1 {
		return nil, fmt.Errorf("invalid commit object: no null separator found")
	}
	// Extract content after null byte
	content := string(rawCommitData[nullIndex+1:])
	lines := strings.Split(content, "\n")

	i := 0

	// Parse tree
	var treeHash common.Hash
	if i < len(lines) && strings.HasPrefix(lines[i], "tree ") {
		treeHashStr := strings.TrimPrefix(lines[i], "tree ")
		hashBytes, _ := hex.DecodeString(treeHashStr)
		copy(treeHash[:], hashBytes)
		i++
	}

	// Parse parents
	var parents []common.Hash
	for i < len(lines) && strings.HasPrefix(lines[i], "parent ") {
		parentHashStr := strings.TrimPrefix(lines[i], "parent ")
		var parentHash common.Hash
		hashBytes, _ := hex.DecodeString(parentHashStr)
		copy(parentHash[:], hashBytes)
		parents = append(parents, parentHash)
		i++
	}

	// Parse author and timestamp
	var author string
	var timestamp time.Time
	if i < len(lines) && strings.HasPrefix(lines[i], "author ") {
		authorLine := strings.TrimPrefix(lines[i], "author ")
		parts := strings.Fields(authorLine)
		if len(parts) >= 3 {
			timestampStr := parts[len(parts)-2]
			author = strings.Join(parts[:len(parts)-2], " ")
			ts, err := strconv.ParseInt(timestampStr, 10, 64)
			if err == nil {
				timestamp = time.Unix(ts, 0)
			}
		}
		i++
	}

	// Skip committer line
	if i < len(lines) && strings.HasPrefix(lines[i], "committer ") {
		i++
	}

	// Skip empty line
	if i < len(lines) && lines[i] == "" {
		i++
	}

	// Rest is commit message
	var message string
	if i < len(lines) {
		message = strings.Join(lines[i:], "\n")
	}

	return &Commit{
		Tree:      &Tree{Hash: treeHash},
		Parents:   parents,
		Author:    author,
		Timestamp: timestamp,
		Message:   message,
	}, nil
}
