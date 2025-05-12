package objects

type Index struct {
	NumberOfEntries uint32
	Entries         []*IndexEntry
}

type IndexEntry struct {
	CreationTimeSec, CreationTimeNano uint32
	ModifiedTimeSec, ModifiedTimeNano uint32
	Hash                              [20]byte
	FileSize                          uint32
	FileMode                          uint32
	EntryPath                         string
}
