package java

import (
	"net/http"
	"strings"

	"../../events"
	"../../index"
	"../../storage"
	"github.com/gorilla/mux"
)

// JavaEndpoints : API for serving Java Requests
type JavaEndpoints struct {
	ArtifactIndex   index.Index
	ArtifactStorage storage.Storage
}

// AppendEndpoints : In Java, Artifacts are saved with xml metadata at the artifact level as well as the version level
func (je JavaEndpoints) AppendEndpoints(r *mux.Router) {
	// Hmm, still missing the type group :  (?:\\.){type:\\w*}
	r.HandleFunc("/java/{group:.+}/{artifact:.+}/{version:.+}/{filename:[^/]+}{type:\\.jar|\\.pom}", je.javaDownloadArtifactRouter)
	r.HandleFunc("/java/{group:.+}/{artifact:.+}/{version:.+}/{filename:[^/]+}{type:\\.jar|\\.pom}{checksum:\\.md5|\\.sha1}", je.javaDownloadArtifactChecksumRouter)
	// r.HandleFunc("/java/{group:.+}/{artifact:.+}/maven-metadata.xml", je.javaDownloadTopLevelMetadataHandler)
	// r.HandleFunc("/java/{group:.+}/{artifact:.+}/{version:([.\\d\\.]*\\-SNAPSHOT)+}/maven-metadata.xml{type:((?:\\.)(*:\\w))}", je.javaDownloadMetadataHandler)
}

// func (je JavaEndpoints) javaDownloadTopLevelMetadataHandler(w http.ResponseWriter, r *http.Request) {

// 	artifact := RequestToArtifact(r)
// 	artifact.Filename = "maven-metadata.xml"

// 	fmt.Println("Top Level Metadata Download : JAVA, G(" + artifact.Group + ") A(" + artifact.Artifact + ") F(" + artifact.Filename + ")")
// }

// func (je JavaEndpoints) javaDownloadMetadataHandler(w http.ResponseWriter, r *http.Request) {

// 	artifact := RequestToArtifact(r)
// 	artifact.Filename = "maven-metadata.xml"

// 	fmt.Println("Metadata Download : JAVA, G(" + artifact.Group + ") A(" + artifact.Artifact + ") F(" + artifact.Filename + ")")
// }

// For checksums, we want to attempt to download something from the proxy - but if it doesn't exist, it's probably
// better if we generate a checksum using the artifact we've got as some files don't have a checksum.
// TODO: Maybe we make this choice configurable?
func (je JavaEndpoints) javaDownloadArtifactChecksumRouter(w http.ResponseWriter, r *http.Request) {

	ja := RequestToArtifact(r)
	events.Info("hangar.request.java.downloadChecksum", ja.ToString())
	je.javaProxiedArtifactAction(w, r, ja)
}

func (je JavaEndpoints) javaDownloadArtifactRouter(w http.ResponseWriter, r *http.Request) {

	ja := RequestToArtifact(r)
	events.Info("hangar.request.java.downloadArtifact", ja.ToString())
	je.javaProxiedArtifactAction(w, r, ja)
}

func (je JavaEndpoints) javaProxiedArtifactAction(w http.ResponseWriter, r *http.Request, ja Artifact) {

	// If file exists in index - attempt to serve,
	if !je.ArtifactIndex.IsDownloadedArtifact(ja.GetIdentifier(), ja.Type) {

		events.Debug("hangar.request.java.download.ifIndex", ja.Type)

		mavenCentral := "https://repo.maven.apache.org/maven2/" + strings.Replace(ja.Group, ".", "/", -1) + "/" + ja.Artifact + "/" + ja.Version + "/" + ja.Filename
		je.ArtifactStorage.DownloadArtifactToStorage(mavenCentral, ja.GetStorageIdentifier())

		// Add to the index
		if !je.ArtifactIndex.IsArtifact(ja.GetIdentifier()) {
			je.ArtifactIndex.AddArtifact(ja.GetIdentifier(), NewJavaFileList())
		}
		je.ArtifactIndex.AddDownloadedArtifact(ja.GetIdentifier(), ja.Type)
	}

	// Serve the File to the User
	je.ArtifactStorage.ServeFile(w, r, ja.GetStorageIdentifier())
}
