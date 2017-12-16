package index

import (
	"strings"

	"github.com/spf13/viper"
)

// Index : Interface for the Index Core Module
type Index interface {
	IsArtifact(id Identifier) bool
	IsDownloadedArtifact(id Identifier, filetype string) bool
	AddArtifact(id Identifier, filetypes FileList)
	AddDownloadedArtifact(id Identifier, filetype string)
	CountAll() int
}

// Identifier : Basic building block of the index.
type Identifier struct {
	Key string `json:"key"`
}

// FileList : List of FileTypes that are allowed for this artefact and if
// they have been downloaded to the filesystem yet.
type FileList struct {
	FileTypes map[string]bool
}

// BuildIndex : Returns an initialised index of the type at the config key.
func BuildIndex(indexConfigKey string) Index {

	// Initialises and returns an instance of InMemory
	if strings.Compare(viper.GetString(indexConfigKey+".type"), "in-memory") == 0 {
		return newInMemory()
	}

	// If index configuration is not complete or cannot be instantiated, return nil
	return nil
}
