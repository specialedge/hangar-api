package index

// Note - maps are not safe for concurrent use (https://golang.org/doc/faq#atomic_maps)
type inMemory struct {
	artifacts map[Identifier]FileList
}

func NewInMemory() Index {
	return &inMemory{
		artifacts: make(map[Identifier]FileList),
	}
}

// IsArtifact : Does this artifact exist in the system?
func (i inMemory) IsArtifact(key Identifier) bool {
	_, ok := i.artifacts[key]
	return ok
}

// IsDownloadedArtifact : Does this artifact exist in the system?
func (i inMemory) IsDownloadedArtifact(key Identifier, filetype string) bool {
	return i.artifacts[key].FileTypes[filetype]
}

func (i inMemory) AddArtifact(key Identifier, filetypes FileList) {
	i.artifacts[key] = filetypes
}

func (i inMemory) AddDownloadedArtifact(key Identifier, filetype string) {
	i.artifacts[key].FileTypes[filetype] = true
}
