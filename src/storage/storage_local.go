package storage

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"../events"
	"github.com/cavaliercoder/grab"
)

type storageLocal struct {
	Path string
}

func NewStorageLocal() Storage {
	return &storageLocal{
		Path: `F:\Code\specialedge\storage\java\`,
	}
}

// DownloadArtifactToStorage : Download the artifact from the URI to Storage
func (s storageLocal) DownloadArtifactToStorage(uri string, id Identifier) {

	// Create a Custom Client
	client := grab.NewClient()
	client.UserAgent = "Hangar v0.0.1"

	// Create a Download Request
	req, err := grab.NewRequest(filepath.Join(s.Path, id.Key), uri)
	if err != nil {
		panic(err)
	}

	events.Debug("hangar.storage.download", uri)
	resp := client.Do(req)

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}
}

func (s storageLocal) ServeFile(w http.ResponseWriter, r *http.Request, id Identifier) {
	events.Debug("hangar.storage.serve", id.Key)
	http.ServeFile(w, r, filepath.Join(s.Path, id.Key))
}