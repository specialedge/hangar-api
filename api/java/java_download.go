package java

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/specialedge/hangar-api/index"
	"github.com/specialedge/hangar-api/storage"
)

// Endpoints : API for serving Java Requests
type Endpoints struct {
	ArtifactIndex   index.Index
	ArtifactStorage storage.Storage
}

// AppendEndpoints : In Java, Artifacts are saved with xml metadata at the artifact level as well as the version level
func (je Endpoints) AppendEndpoints(r *mux.Router) {
	// Hmm, still missing the type group :  (?:\\.){type:\\w*}
	r.HandleFunc("/java/{group:.+}/{artifact:.+}/{version:.+}/{filename:[^/]+}{type:\\.jar|\\.pom}", je.javaDownloadArtifactRouter)
	r.HandleFunc("/java/{group:.+}/{artifact:.+}/{version:.+}/{filename:[^/]+}{type:\\.jar|\\.pom}{checksum:\\.md5|\\.sha1}", je.javaDownloadArtifactChecksumRouter)
}

// For checksums, we want to attempt to download something from the proxy - but if it doesn't exist, it's probably
// better if we generate a checksum using the artifact we've got as some files don't have a checksum.
// TODO: Maybe we make this choice configurable?
func (je Endpoints) javaDownloadArtifactChecksumRouter(w http.ResponseWriter, r *http.Request) {

	ja := RequestToArtifact(mux.Vars(r))
	log.WithFields(log.Fields{"module": "api", "action": "downloadChecksum"}).Info(ja.ToString())
	je.javaProxiedArtifactAction(w, r, ja)
}

func (je Endpoints) javaDownloadArtifactRouter(w http.ResponseWriter, r *http.Request) {

	ja := RequestToArtifact(mux.Vars(r))
	log.WithFields(log.Fields{"module": "api", "action": "downloadArtifact"}).Info(ja.ToString())
	je.javaProxiedArtifactAction(w, r, ja)
}

func (je Endpoints) javaProxiedArtifactAction(w http.ResponseWriter, r *http.Request, ja Artifact) {

	// If file exists in index - attempt to serve,
	if !je.ArtifactIndex.IsDownloadedArtifact(ja.GetIdentifier(), ja.Type) {

		mavenCentral := "https://repo.maven.apache.org/maven2/" + strings.Replace(ja.Group, ".", "/", -1) + "/" + ja.Artifact + "/" + ja.Version + "/" + ja.Filename
		je.ArtifactStorage.DownloadArtifactToStorage(mavenCentral, ja.GetStorageIdentifier())

		addJavaArtifactToIndex(je, ja)
	}

	// Serve the File to the User
	je.ArtifactStorage.ServeFile(w, r, ja.GetStorageIdentifier())
}

// ReIndex : Has the index populate itself from the storage using the model for Java Artifacts
func (je Endpoints) ReIndex() {
	idents := je.ArtifactStorage.GetArtifacts()

	for _, file := range idents {
		ja := StorageIdentifierToArtifact(file)
		addJavaArtifactToIndex(je, ja)
	}

	log.WithFields(log.Fields{"module": "api", "action": "ReIndex"}).Info(strconv.Itoa(je.ArtifactIndex.CountAll()) + " current artifacts registered.")
}

func addJavaArtifactToIndex(je Endpoints, ja Artifact) {
	// Add to the index
	if !je.ArtifactIndex.IsArtifact(ja.GetIdentifier()) {
		je.ArtifactIndex.AddArtifact(ja.GetIdentifier(), NewJavaFileList())
	}

	je.ArtifactIndex.AddDownloadedArtifact(ja.GetIdentifier(), ja.Type)
}
