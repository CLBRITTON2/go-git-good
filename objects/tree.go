package objects

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/CLBRITTON2/go-git-good/common"
)

type Tree struct {
	Entries []*TreeEntry
	Hash    common.Hash
}

type TreeEntry struct {
	Name     string
	FileMode uint32
	Hash     common.Hash
}

func BuildTreeFromIndex(index *common.Index) (*Tree, error) {
	// Holds the hierarchy structure of the root tree and all of its subtrees with their associated files and subtrees
	trees := make(map[string]*Tree)

	for _, entry := range index.Entries {
		// Build empty trees to represent each directory and subdirectory within each index entry filepath
		// Each entry filepath represents its relative location in the working tree
		directoryPath := entry.EntryPath
		for {
			directoryPath = filepath.Dir(directoryPath)
			directoryPath = normalizeDirectoryPath(directoryPath)
			// Create tree for this directory if it doesn't exist
			if _, exists := trees[directoryPath]; !exists {
				trees[directoryPath] = &Tree{Entries: []*TreeEntry{}}
			}

			// Stop when we reach the root
			if directoryPath == "" {
				break
			}
		}

		// Add the file entry to its immediate parent directory tree
		directory := filepath.Dir(entry.EntryPath)
		directory = normalizeDirectoryPath(directory)

		trees[directory].Entries = append(trees[directory].Entries, &TreeEntry{
			Name:     filepath.Base(entry.EntryPath),
			FileMode: entry.FileMode,
			Hash:     entry.Hash,
		})
	}

	// Sort trees by depth - trees have to be built from the bottom up because parent
	// trees need the hash of their children to compute their own hash
	var directoriesByDepth []string
	for directory := range trees {
		directoriesByDepth = append(directoriesByDepth, directory)
	}

	sort.Slice(directoriesByDepth, func(i, j int) bool {
		// Count path separators to determine depth
		depthI := strings.Count(directoriesByDepth[i], string(filepath.Separator))
		depthJ := strings.Count(directoriesByDepth[j], string(filepath.Separator))

		// If depths are equal, force root to be last to enforce bottom-up tree building
		// This is required because if 0 separators are present entries are sorted based on the order they appear in the slice
		if depthI == depthJ {
			if directoriesByDepth[i] == "" {
				return false
			}
			if directoriesByDepth[j] == "" {
				return true
			}
			return directoriesByDepth[i] < directoriesByDepth[j]
		}
		return depthI > depthJ
	})

	// Build trees from the bottom up: start
	for _, directory := range directoriesByDepth {
		tree := trees[directory]
		serializedTreeData := tree.Serialize()
		hash, err := common.HashObject("tree", serializedTreeData)
		if err != nil {
			return nil, err
		}
		tree.Hash = hash

		// Add this tree to its parent
		parent := filepath.Dir(directory)
		parent = normalizeDirectoryPath(parent)

		trees[parent].Entries = append(trees[parent].Entries, &TreeEntry{
			Name:     filepath.Base(directory),
			FileMode: 040000,
			Hash:     hash,
		})
	}

	// Return root
	return trees[""], nil
}

// Convert "." to an empty string for relative root directories for hashing consistency
func normalizeDirectoryPath(path string) string {
	if path == "." {
		return ""
	}
	return path
}

func (tree *Tree) Serialize() []byte {
	// Git sorts entries by name so this allows us to cross compare with git tree hashes
	sort.Slice(tree.Entries, func(i, j int) bool {
		return tree.Entries[i].Name < tree.Entries[j].Name
	})

	var buffer bytes.Buffer
	for _, entry := range tree.Entries {
		// Format for tree entries: [mode] [name]\0[hash]
		// Note that this is different from the format of the root tree which is handled by common.HashObject()
		entryHeader := fmt.Sprintf("%o %s\x00", entry.FileMode, entry.Name)
		buffer.Write([]byte(entryHeader))
		buffer.Write(entry.Hash[:])
	}

	return buffer.Bytes()
}
