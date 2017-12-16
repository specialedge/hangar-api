package storage

import "net/http"

// Identifier : Basic building block of the index.
type Identifier struct {
	Key       string `json:"key"`
	Separator string `json:"separator"`
}

// Storage : Interface for the Storage Core Module
type Storage interface {
	DownloadArtifactToStorage(uri string, id Identifier)
	ServeFile(w http.ResponseWriter, r *http.Request, id Identifier)
	GetArtifacts() []Identifier
}
