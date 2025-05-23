package objects

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
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

func BuildTreeFromIndex(index *common.Index) (*Tree, map[string]*Tree, error) {
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
			_, exists := trees[directoryPath]
			if !exists {
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
		// Don't add the root tree to itself
		if directory == "" {
			continue
		}
		tree := trees[directory]
		serializedTreeData := tree.Serialize()
		hash, err := common.HashObject(serializedTreeData)
		if err != nil {
			return nil, nil, err
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

	// Calculate hash for the root tree
	rootTree := trees[""]
	serializedRootTreeData := rootTree.Serialize()
	rootHash, err := common.HashObject(serializedRootTreeData)
	if err != nil {
		return nil, nil, err
	}
	rootTree.Hash = rootHash

	// Return root and trees map so each tree can be written to the object DB
	return rootTree, trees, nil
}

// Convert "." to an empty string for relative root directories (hashing consistency with Git)
func normalizeDirectoryPath(path string) string {
	if path == "." {
		return ""
	}
	return path
}

// This creates the format for tree entries and the root tree to be written to the DB
func (tree *Tree) Serialize() []byte {
	// Git sorts entries by name so this allows us to cross compare with git tree hashes
	sort.Slice(tree.Entries, func(i, j int) bool {
		return tree.Entries[i].Name < tree.Entries[j].Name
	})

	var buffer bytes.Buffer
	for _, entry := range tree.Entries {
		// Format for tree entries: [mode] [name]\0[hash]
		entryHeader := fmt.Sprintf("%o %s\x00", entry.FileMode, entry.Name)
		buffer.Write([]byte(entryHeader))
		buffer.Write(entry.Hash[:])
	}

	// Add tree header: "tree <length>\x00"
	entryData := buffer.Bytes()
	header := fmt.Sprintf("tree %d\x00", len(entryData))
	return append([]byte(header), entryData...)
}

func ParseTree(rawTreeData []byte) (*Tree, error) {
	// Expected input data format: "tree <length>\x00[mode] [name]\x00[hash]..."
	// Process tree header
	nullIndex := bytes.IndexByte(rawTreeData, byte('\x00'))
	if nullIndex == -1 {
		return nil, fmt.Errorf("invalid tree format: no null byte found parsing tree object")
	}

	header := string(rawTreeData[:nullIndex])
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid tree header format expected tree <data length> got: %s", header)
	}

	if parts[0] != "tree" {
		return nil, fmt.Errorf("non-tree object passed to ParseTree function: %v", parts[0])
	}

	// Process tree entries
	entries := rawTreeData[nullIndex+1:]
	tree := &Tree{Entries: []*TreeEntry{}}
	offset := 0

	for offset < len(entries) {
		// Find end of mode+name section
		entryNullIndex := bytes.IndexByte(entries[offset:], byte('\x00'))
		if entryNullIndex == -1 {
			return nil, fmt.Errorf("malformed tree entry: missing null byte")
		}

		modeAndName := string(entries[offset : offset+entryNullIndex])
		modeAndNameParts := strings.SplitN(modeAndName, " ", 2)
		if len(modeAndNameParts) != 2 {
			return nil, fmt.Errorf("malformed tree entry: invalid mode/name format")
		}
		fileModeString := modeAndNameParts[0]
		fileMode64, err := strconv.ParseUint(fileModeString, 8, 32)
		if err != nil {
			return nil, err
		}
		fileMode := uint32(fileMode64)
		name := modeAndNameParts[1]

		offset += entryNullIndex + 1
		// Ensure we have enough bytes for the hash
		if offset+20 > len(entries) {
			return nil, fmt.Errorf("malformed tree entry: insufficient bytes for hash")
		}

		hash := common.Hash(entries[offset : offset+20])
		tree.Entries = append(tree.Entries, &TreeEntry{
			Name:     name,
			FileMode: fileMode,
			Hash:     hash,
		})

		// Move offset past the hash to the next entry
		offset += 20
	}
	return tree, nil
}
