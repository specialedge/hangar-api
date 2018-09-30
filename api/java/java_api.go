package java

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/specialedge/hangar-api/index"
	"github.com/specialedge/hangar-api/storage"
	"github.com/spf13/viper"
)

// endpoints : API for serving Java Requests
type endpoints struct {
	ArtifactIndex   index.Index
	ArtifactStorage storage.Storage
	Proxies         []*url.URL
}

// InitialiseJavaEndpoints : Initialise the object that contains the endpoints,
// bringing together the requested storage and indexing mechanism.
func InitialiseJavaEndpoints(r *mux.Router) {

	stor := storage.BuildStorage("java.storage")
	ind := index.BuildIndex("java.index")

	// Provided an option has been given for both...
	if stor != nil && ind != nil {

		// Create Endpoints object
		javaEndpoints := endpoints{
			ArtifactIndex:   ind,
			ArtifactStorage: stor,
			Proxies:         prepareProxies(),
		}

		// Add all the endpoints for the Java API
		javaEndpoints.AppendEndpoints(r)

		// If required, re-populate the index.
		if viper.GetBool("java.index.reindex") {
			javaEndpoints.ReIndex()
		}
	}
}

// prepareProxies : Cycle through the submitted proxies and sanitize them.
func prepareProxies() []*url.URL {

	// Cycle through the repositories that are available.
	configProxies := viper.GetStringSlice("java.proxies")
	var proxies []*url.URL

	if len(configProxies) > 0 {
		for _, proxy := range configProxies {
			url, err := url.ParseRequestURI(proxy)
			if err != nil {
				log.WithFields(log.Fields{"module": "api", "action": "InitialiseJavaEndpoints"}).Error(proxy + " is not a valid URI, ignoring...")
			} else {
				proxies = append(proxies, url)
				log.WithFields(log.Fields{"module": "api", "action": "InitialiseJavaEndpoints"}).Info(url.String() + " configured as Java Proxy.")
			}
		}
	}

	// If we've not managed to find any actual URLs to use as a proxy.
	if len(proxies) == 0 {
		uri, _ := url.Parse("https://repo.maven.apache.org/maven2/")
		proxies = []*url.URL{uri}
		log.WithFields(log.Fields{"module": "api", "action": "InitialiseJavaEndpoints"}).Info("Using default proxy of " + uri.String())
	}

	return proxies
}

// AppendEndpoints : In Java, Artifacts are saved with xml metadata at the artifact level as well as the version level
func (je endpoints) AppendEndpoints(r *mux.Router, handlers ...func(w http.ResponseWriter, r *http.Request)) {

	var javaDarFunc func(w http.ResponseWriter, r *http.Request)
	var javaDacrFunc func(w http.ResponseWriter, r *http.Request)
	var javaUsarFunc func(w http.ResponseWriter, r *http.Request)
	var javaUsacrFunc func(w http.ResponseWriter, r *http.Request)

	if handlers == nil {
		// If we don't submit any options, we should go for the default functions to handle these endpoints.
		javaDarFunc = je.javaDownloadArtifactRouter
		javaDacrFunc = je.javaDownloadArtifactChecksumRouter
		javaUsarFunc = je.javaUploadSnapshotArtifactRouter
		javaUsacrFunc = je.javaUploadSnapshotArtifactChecksumRouter
	} else {
		// However, for perhaps test purposes - if we want to override these handlers then we should.
		javaDarFunc = handlers[0]
		javaDacrFunc = handlers[1]
		javaUsarFunc = handlers[2]
		javaUsacrFunc = handlers[3]
	}

	// Main download APIs
	r.HandleFunc("/java/{group:.+}/{artifact:.+}/{version:.+}/{filename:[^/]+}{type:\\.jar|\\.pom}", javaDarFunc).Methods("GET")
	r.HandleFunc("/java/{group:.+}/{artifact:.+}/{version:.+}/{filename:[^/]+}{type:\\.jar|\\.pom}{checksumType:\\.md5|\\.sha1}", javaDacrFunc).Methods("GET")

	// Snapshot Upload APIs
	r.HandleFunc("/java/snapshots/{group:.+}/{artifact:.+}/{version:.+-SNAPSHOT}/{filename:[^/]+}{type:\\.jar|\\.pom}", javaUsarFunc).Methods("PUT")
	r.HandleFunc("/java/snapshots/{group:.+}/{artifact:.+}/{version:.+-SNAPSHOT}/{filename:[^/]+}{type:\\.jar|\\.pom}{checksumType:\\.md5|\\.sha1}", javaUsacrFunc).Methods("PUT")
}

// ReIndex : Has the index populate itself from the storage using the model for Java Artifacts
func (je endpoints) ReIndex() {
	idents := je.ArtifactStorage.GetArtifacts()

	for _, file := range idents {
		ja := StorageIdentifierToArtifact(file)
		addJavaArtifactToIndex(je, ja)
	}

	log.WithFields(log.Fields{"module": "api", "action": "ReIndex"}).Info(strconv.Itoa(je.ArtifactIndex.CountAll()) + " current artifacts registered.")
}
