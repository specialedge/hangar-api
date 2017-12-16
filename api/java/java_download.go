package java

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/specialedge/hangar-api/index"
	"github.com/specialedge/hangar-api/storage"
	"github.com/spf13/viper"
)

// Endpoints : API for serving Java Requests
type Endpoints struct {
	ArtifactIndex   index.Index
	ArtifactStorage storage.Storage
}

// InitialiseJavaEndpoints : Initialise the object that contains the endpoints,
// bringing together the requested storage and indexing mechanism.
func InitialiseJavaEndpoints(r *mux.Router) {

	stor := storage.BuildStorage("java.storage")
	ind := index.BuildIndex("java.index")

	// Provided an option has been given for both...
	if stor != nil && ind != nil {

		// Create Endpoints object
		javaEndpoints := Endpoints{
			ArtifactIndex:   ind,
			ArtifactStorage: stor,
		}

		// Add all the endpoints for the Java API
		javaEndpoints.AppendEndpoints(r)

		// If required, re-populate the index.
		if viper.GetBool("java.index.reindex") {
			javaEndpoints.ReIndex()
		}
	}
}

// AppendEndpoints : In Java, Artifacts are saved with xml metadata at the artifact level as well as the version level
func (je Endpoints) AppendEndpoints(r *mux.Router, handlers ...func(w http.ResponseWriter, r *http.Request)) {

	var javaDarFunc func(w http.ResponseWriter, r *http.Request)
	var javaDacrFunc func(w http.ResponseWriter, r *http.Request)

	if handlers == nil {
		// If we don't submit any options, we should go for the default functions to handle these endpoints.
		javaDarFunc = je.javaDownloadArtifactRouter
		javaDacrFunc = je.javaDownloadArtifactChecksumRouter
	} else {
		// However, for perhaps test purposes - if we want to override these handlers then we should.
		javaDarFunc = handlers[0]
		javaDacrFunc = handlers[1]
	}

	r.HandleFunc("/java/{group:.+}/{artifact:.+}/{version:.+}/{filename:[^/]+}{type:\\.jar|\\.pom}", javaDarFunc)
	r.HandleFunc("/java/{group:.+}/{artifact:.+}/{version:.+}/{filename:[^/]+}{type:\\.jar|\\.pom}{checksumType:\\.md5|\\.sha1}", javaDacrFunc)
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
