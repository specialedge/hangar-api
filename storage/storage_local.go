package storage

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cavaliercoder/grab"
	log "github.com/sirupsen/logrus"
)

type storageLocal struct {
	Path string
}

// NewStorageLocal : Creates a new Storage module which uses the local disk.
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

	log.WithFields(log.Fields{"module": "storage", "action": "DownloadArtifactToStorage"}).Debug(uri)
	resp := client.Do(req)

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}
}

// ServeFile : Requests that the storage serve the user with the artifact.
func (s storageLocal) ServeFile(w http.ResponseWriter, r *http.Request, id Identifier) {
	log.WithFields(log.Fields{"module": "storage", "action": "ServeFile"}).Debug(id.Key)
	http.ServeFile(w, r, filepath.Join(s.Path, id.Key))
}

// GetArtifacts : Returns an array of Storage Identifiers that
func (s storageLocal) GetArtifacts() []Identifier {
	fileList := []Identifier{}

	filepath.Walk(s.Path, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			base := strings.Replace(path, s.Path, "", 1)
			fileList = append(fileList, Identifier{
				Key:       strings.Replace(base, "\\", "/", -1),
				Separator: "/",
			})
		}
		return nil
	})

	log.WithFields(log.Fields{"module": "storage", "action": "GetArtifacts"}).Info(strconv.Itoa(len(fileList)) + " entities retrieved from storage...")
	return fileList
}
