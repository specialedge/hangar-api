package storage

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Identifier : Basic building block of the index.
type Identifier struct {
	Key       string `json:"key"`
	Separator string `json:"separator"`
}

// Storage : Interface for the Storage Core Module
type Storage interface {
	DownloadArtifactToStorage(uri string, id Identifier, codes ...int) (int, error)
	ServeFile(w http.ResponseWriter, r *http.Request, id Identifier)
	GetArtifacts() []Identifier
}

// BuildStorage : Returns an initialised storage based on the config key.
func BuildStorage(storageConfigKey string) Storage {

	// Initialises and returns an instance of LocalStorage
	if strings.Compare(viper.GetString(storageConfigKey+".type"), "local") == 0 {
		return NewStorageLocal(viper.GetString(storageConfigKey + ".path"))
	}

	// If storage configuration is not complete or cannot be instantiated, return nil
	log.WithFields(log.Fields{"module": "storage", "action": "BuildStorage"}).Error("Could not instantiate storage!")
	return nil
}
