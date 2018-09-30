package java

import (
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// For checksums, we want to attempt to download something from the proxy - but if it doesn't exist, it's probably
// better if we generate a checksum using the artifact we've got as some files don't have a checksum.
// TODO: Maybe we make this choice configurable?
func (je endpoints) javaDownloadArtifactChecksumRouter(w http.ResponseWriter, r *http.Request) {

	ja := RequestToArtifact(mux.Vars(r))
	log.WithFields(log.Fields{"module": "api", "action": "downloadChecksum"}).Info(ja.ToString())
	je.javaProxiedArtifactAction(w, r, ja)
}

func (je endpoints) javaDownloadArtifactRouter(w http.ResponseWriter, r *http.Request) {

	ja := RequestToArtifact(mux.Vars(r))
	log.WithFields(log.Fields{"module": "api", "action": "downloadArtifact"}).Info(ja.ToString())
	je.javaProxiedArtifactAction(w, r, ja)
}

func (je endpoints) javaProxiedArtifactAction(w http.ResponseWriter, r *http.Request, ja Artifact) {

	// If the file does not exist in the index and cannot be served from within the system...
	if !je.ArtifactIndex.IsDownloadedArtifact(ja.GetIdentifier(), ja.Type) {

		// Cycle through the repositories that are available.
		for _, proxy := range je.Proxies {

			// Form the URL to hit in order to get an artifact
			uri := *proxy
			uri.Path = path.Join(uri.Path, strings.Replace(ja.Group, ".", "/", -1), ja.Artifact, ja.Version, ja.Filename)

			// Attempt to download and store the artefact.
			code, err := je.ArtifactStorage.DownloadArtifactToStorage(uri.String(), ja.GetStorageIdentifier(), 200)

			// If there's been a problem, try the next proxy - otherwise, add the artifact to the index and break.
			if code == 200 {
				addJavaArtifactToIndex(je, ja)
				break
			}
			if err != nil {
				log.WithFields(log.Fields{"module": "api", "action": "javaProxiedArtifactAction"}).Errorln(err)
			}
		}
	}

	// Attempt to serve the File to the User
	je.ArtifactStorage.ServeFile(w, r, ja.GetStorageIdentifier())
}

func addJavaArtifactToIndex(je endpoints, ja Artifact) {
	// Add to the index
	if !je.ArtifactIndex.IsArtifact(ja.GetIdentifier()) {
		je.ArtifactIndex.AddArtifact(ja.GetIdentifier(), NewJavaFileList())
	}

	je.ArtifactIndex.AddDownloadedArtifact(ja.GetIdentifier(), ja.Type)
}
