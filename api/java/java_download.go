package java

import (
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

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

// ReIndex : Has the index populate itself from the storage using the model for Java Artifacts
func (je endpoints) ReIndex() {
	idents := je.ArtifactStorage.GetArtifacts()

	for _, file := range idents {
		ja := StorageIdentifierToArtifact(file)
		addJavaArtifactToIndex(je, ja)
	}

	log.WithFields(log.Fields{"module": "api", "action": "ReIndex"}).Info(strconv.Itoa(je.ArtifactIndex.CountAll()) + " current artifacts registered.")
}

func addJavaArtifactToIndex(je endpoints, ja Artifact) {
	// Add to the index
	if !je.ArtifactIndex.IsArtifact(ja.GetIdentifier()) {
		je.ArtifactIndex.AddArtifact(ja.GetIdentifier(), NewJavaFileList())
	}

	je.ArtifactIndex.AddDownloadedArtifact(ja.GetIdentifier(), ja.Type)
}
