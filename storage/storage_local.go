package storage

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type storageLocal struct {
	Path       string
	FileSystem afero.Fs
}

// NewStorageLocal : Creates a new Storage module which uses the local disk.
func NewStorageLocal(path string) Storage {

	// We are using afero in order to make mocking easier. To switch to in-memory, use afero.NewMemMapFs()
	var AppFs = afero.NewOsFs()

	// If the path doesn't exist, create it.
	if _, err := AppFs.Stat(path); os.IsNotExist(err) {
		log.WithFields(log.Fields{"module": "storage", "action": "NewStorageLocal"}).Info("Creating Empty Directory : " + path)
		AppFs.Mkdir(path, 0755)
	}

	return &storageLocal{
		Path:       path,
		FileSystem: AppFs,
	}
}

// DownloadArtifactToStorage : Download the artifact from the URI to Storage.
// We pass in a set of whitelisted status codes to accept from the proxy server.
func (s storageLocal) DownloadArtifactToStorage(uri string, id Identifier, codes ...int) (int, error) {

	// Create a Custom Client
	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	filename := filepath.Join(s.Path, id.Key)

	// We choose to attempt to get the data first.
	// It is far more likely that this request will fail (as some repositories may
	// not actually have the artefact we are looking for) so we want to find that out
	// before creating or saving a file.
	log.WithFields(log.Fields{"module": "storage", "action": "DownloadArtifactToStorage"}).Info(uri)
	log.WithFields(log.Fields{"module": "storage", "action": "DownloadArtifactToStorage"}).Debug(filename)

	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Set("User-Agent", "Hangar v0.0.1")
	resp, err := client.Do(req)

	// If we couldn't form the request, return with an error.
	if err != nil {
		log.WithFields(log.Fields{"module": "storage", "action": "DownloadArtifactToStorage"}).Error("Could not execute request : " + err.Error())
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	// Once we know the request (and response) was valid, we want to check for any status codes that were not whitelisted.
	if resp != nil {
		if !intExists(resp.StatusCode, codes) {
			log.WithFields(log.Fields{"module": "storage", "action": "DownloadArtifactToStorage"}).Error("Could not download file : " + resp.Status)
			return resp.StatusCode, err
		}
	}

	// So we know it's a valid response code. We now want to save the file to disk. Create the directory then the file within it.
	if _, err := s.FileSystem.Stat(filepath.Dir(filename)); os.IsNotExist(err) {
		s.FileSystem.MkdirAll(filepath.Dir(filename), 644)
	}
	out, err := s.FileSystem.Create(filename)

	if err != nil {
		log.WithFields(log.Fields{"module": "storage", "action": "DownloadArtifactToStorage"}).Error("Could not create file at request : " + err.Error())
		return http.StatusInternalServerError, err
	}
	defer out.Close()

	// Write the body to the file.
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return resp.StatusCode, nil
}

// I have no idea why a function like this does not exist within
// the base Go libraries, but there you go.
func intExists(val int, arr []int) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

// ServeFile : Requests that the storage serve the user with the artifact.
func (s storageLocal) ServeFile(w http.ResponseWriter, r *http.Request, id Identifier) {

	filename := filepath.Join(s.Path, id.Key)
	file, err := s.FileSystem.Open(filename)
	if err != nil {
		log.WithFields(log.Fields{"module": "storage", "action": "ServeFile"}).Error("Could not get file at " + filename)
		log.WithFields(log.Fields{"module": "storage", "action": "ServeFile"}).Error(err)
	}
	defer file.Close()

	// Switched to ServeContent to support afero.
	http.ServeContent(w, r, filename, time.Now(), file)
}

// GetArtifacts : Returns an array of Storage Identifiers by traversing the local storage filesystem.
func (s storageLocal) GetArtifacts() []Identifier {
	fileList := []Identifier{}

	afero.Walk(s.FileSystem, s.Path, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			base := strings.Replace(path, filepath.Clean(s.Path), "", 1)
			log.WithFields(log.Fields{"module": "storage", "action": "GetArtifacts"}).Debug(base)
			fileList = append(fileList, Identifier{
				Key:       strings.Replace(base, "\\", "/", -1),
				Separator: "/",
			})
		}
		return nil
	})

	log.WithFields(log.Fields{"module": "storage", "action": "GetArtifacts"}).Info(strconv.Itoa(len(fileList)) + " entities retrieved from storage (" + s.Path + ")")
	return fileList
}

func (s storageLocal) SaveArtifact(w http.ResponseWriter, r *http.Request, id Identifier) {
	filename := filepath.Join(s.Path, id.Key)
	afero.WriteReader(s.FileSystem, filename, r.Body)
}
