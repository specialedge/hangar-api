package python

import (
	"fmt"
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

// endpoints : API for serving Python Requests
type endpoints struct {
	ArtifactIndex   index.Index
	ArtifactStorage storage.Storage
	Proxies         []*url.URL
}

// InitialisePythonEndpoints : Initialise the object that contains the endpoints,
// bringing together the requested storage and indexing mechanism.
func InitialisePythonEndpoints(r *mux.Router) {

	stor := storage.BuildStorage("python.storage")
	ind := index.BuildIndex("python.index")

	// Provided an option has been given for both
	if stor != nil && ind != nil {

		// Create Endpoints object
		pythonEndpoints := endpoints{
			ArtifactIndex:   ind,
			ArtifactStorage: stor,
		}

		// Add all the endpoints for the Python API
		pythonEndpoints.AppendEndpoints(r)

		// If required, re-populate the index.
		if viper.GetBool("python.index.reindex") {
			pythonEndpoints.ReIndex()
		}

		// Set a default for the proxies
		viper.SetDefault("python.proxies", []string{"https://pypi.python.org/packages/"})
	}
}

// AppendEndpoints : Append the artefact download endpoints
func (pe endpoints) AppendEndpoints(r *mux.Router, handlers ...func(w http.ResponseWriter, r *http.Request)) {

	var pythonDarFunc func(w http.ResponseWriter, r *http.Request)
	var pythonVerFunc func(w http.ResponseWriter, r *http.Request)

	if handlers == nil {
		pythonDarFunc = pe.pythonDownloadArtifactRouter
		pythonVerFunc = pe.pythonVersionArtifactRouter
	} else {
		pythonDarFunc = handlers[0]
		pythonVerFunc = handlers[1]
	}

	r.HandleFunc("/python/{packageType:.+}/{letter:.+}/{packageName:.+}/{packageFile:[^/]+}{type:\\.whl|\\.egg|\\.exe|\\.msi|\\.tar.gz|\\.zip}", pythonDarFunc)
	r.HandleFunc("/python/{packageName:.+}/", pythonVerFunc)
}

func (pe endpoints) pythonDownloadArtifactRouter(w http.ResponseWriter, r *http.Request) {
	pa := RequestToArtifact(mux.Vars(r))
	log.WithFields(log.Fields{"module": "api", "action": "downloadArtifact"}).Info(pa.ToString())
	pe.pythonProxiedArtifactAction(w, r, pa)
}

func (pe endpoints) pythonVersionArtifactRouter(w http.ResponseWriter, r *http.Request) {
	packageName := mux.Vars(r)["packageName"]
	localVersions := RequestToVersions(packageName)
	remoteVersions := pe.FetchRemoteVersions(packageName)
	log.WithFields(log.Fields{"module": "api", "action": "listingArtifactVersions"}).Info(packageName)
	pe.returnVersions(w, r, versions, packageName)
}

func (pe endpoints) FetchRemoteVersions(packageName string) []string {
	var remoteVersions []string
	for _, proxy := range pe.Proxies {
		// Scrape the versions from each proxy in turn here and merge the results
	}
	return remoteVersions
}

func (pe endpoints) returnVersions(w http.ResponseWriter, r *http.Request, versions []string, packageName string) {
	// HTML templates for pypi style version page
	header := "<html><head><title>python: links for %s</title></head>\n<body><h1>python: links for %s</h1>"
	item := "<a href=\"%s/%s\">%s</a><br/>"
	footer := "</body></html>"
	if len(versions) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not Found (%s does not have any releases)", packageName)
	} else {
		var res []string
		res = append(res, fmt.Sprintf(header, packageName, packageName))
		// Loop through all the items we know about and add them to the list
		for _, ver := range versions {
			filename := GetFilename(ver, packageName)
			res = append(res, fmt.Sprintf(item, packageName, filename, filename))
		}
		res = append(res, footer)
		fmt.Fprintf(w, strings.Join(res, "\n"))
	}
}

func (pe endpoints) pythonProxiedArtifactAction(w http.ResponseWriter, r *http.Request, pa Artifact) {

	// If file exists in index - attempt to serve,
	if !pe.ArtifactIndex.IsDownloadedArtifact(pa.GetIdentifier(), pa.PackageFile) {
		// Cycle through the repositories that are available.
		proxies := viper.GetStringSlice("python.proxies")
		for _, proxy := range proxies {

			u, _ := url.Parse(proxy)
			u.Path = path.Join(u.Path, pa.PackageType, pa.Letter, pa.PackageName, pa.PackageFile)

			pe.ArtifactStorage.DownloadArtifactToStorage(u.String(), pa.GetStorageIdentifier())
		}

		addPythonArtifactToIndex(pe, pa)
	}

	// Serve the File to the User
	pe.ArtifactStorage.ServeFile(w, r, pa.GetStorageIdentifier())
}

// ReIndex : Has the index populate itself from the storage using the model for Python Artifacts
func (pe endpoints) ReIndex() {
	idents := pe.ArtifactStorage.GetArtifacts()

	for _, file := range idents {
		pa := StorageIdentifierToArtifact(file)
		addPythonArtifactToIndex(pe, pa)
	}

	log.WithFields(log.Fields{"module": "api", "action": "ReIndex"}).Info(strconv.Itoa(pe.ArtifactIndex.CountAll()) + " current artifacts registered.")
}

func addPythonArtifactToIndex(pe endpoints, pa Artifact) {
	// Add to the index
	if !pe.ArtifactIndex.IsArtifact(pa.GetIdentifier()) {
		pe.ArtifactIndex.AddArtifact(pa.GetIdentifier(), NewPythonFileList())
	}

	pe.ArtifactIndex.AddDownloadedArtifact(pa.GetIdentifier(), pa.PackageFile)
}
