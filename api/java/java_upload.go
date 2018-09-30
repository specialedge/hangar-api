package java

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (je endpoints) javaUploadSnapshotArtifactRouter(w http.ResponseWriter, r *http.Request) {

	ja := RequestToArtifact(mux.Vars(r))
	log.WithFields(log.Fields{"module": "api", "action": "uploadSnapshot"}).Info(ja.ToString())
	je.javaUploadArtifactAction(w, r, ja)
}

func (je endpoints) javaUploadSnapshotArtifactChecksumRouter(w http.ResponseWriter, r *http.Request) {

	ja := RequestToArtifact(mux.Vars(r))
	log.WithFields(log.Fields{"module": "api", "action": "uploadSnapshotChecksum"}).Info(ja.ToString())
	je.javaUploadArtifactAction(w, r, ja)
}

func (je endpoints) javaUploadArtifactAction(w http.ResponseWriter, r *http.Request, ja Artifact) {

	je.ArtifactStorage.SaveArtifact(w, r, ja.GetStorageIdentifier())
}
