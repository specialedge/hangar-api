package index

import (
	"strconv"

	"../events"
)

// Note - maps are not safe for concurrent use (https://golang.org/doc/faq#atomic_maps)
type inMemory struct {
	artifacts map[string]FileList
}

func NewInMemory() Index {
	return &inMemory{
		artifacts: make(map[string]FileList),
	}
}

// IsArtifact : Does this artifact identifier exist in the system?
func (i inMemory) IsArtifact(key Identifier) bool {
	_, ok := i.artifacts[key.Key]
	return ok
}

// AddArtifact : Add an identifier for this artifact.
func (i inMemory) AddArtifact(key Identifier, filetypes FileList) {
	i.artifacts[key.Key] = filetypes
}

// IsDownloadedArtifact : Does this artifact exist in the system?
func (i inMemory) IsDownloadedArtifact(key Identifier, filetype string) bool {
	events.Debug("hangar.index.isDownloadedArtifact", key.Key+":"+strconv.FormatBool(i.artifacts[key.Key].FileTypes[filetype]))
	return i.artifacts[key.Key].FileTypes[filetype]
}

// AddDownloadedArtifact : Mark this artifact as downloaded.
func (i inMemory) AddDownloadedArtifact(key Identifier, filetype string) {
	events.Debug("hangar.index.AddDownloadedArtifact", key.Key+":"+filetype)
	i.artifacts[key.Key].FileTypes[filetype] = true
}
