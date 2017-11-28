package java

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"../../index"
	"github.com/cavaliercoder/grab"
	"github.com/gorilla/mux"
)

// JavaEndpoints : API for serving Java Requests
type JavaEndpoints struct {
	Ind index.Index
}

// AppendEndpoints : In Java, Artifacts are saved with xml metadata at the artifact level as well as the version level
func (je JavaEndpoints) AppendEndpoints(r *mux.Router) {
	// Hmm, still missing the type group :  (?:\\.){type:\\w*}
	r.HandleFunc("/java/{group:.+}/{artifact:.+}/{version:.+}/{filename:[^/]+}{type:\\.md5|\\.sha1|\\.jar|\\.pom}", je.javaDownloadArtifactRouter)
	r.HandleFunc("/java/{group:.+}/{artifact:.+}/maven-metadata.xml", je.javaDownloadTopLevelMetadataHandler)
	r.HandleFunc("/java/{group:.+}/{artifact:.+}/{version:([.\\d\\.]*\\-SNAPSHOT)+}/maven-metadata.xml{type:((?:\\.)(*:\\w))}", je.javaDownloadMetadataHandler)
}

func (je JavaEndpoints) javaDownloadTopLevelMetadataHandler(w http.ResponseWriter, r *http.Request) {

	artifact := RequestToArtifact(r)
	artifact.Filename = "maven-metadata.xml"

	fmt.Println("Top Level Metadata Download : JAVA, G(" + artifact.Group + ") A(" + artifact.Artifact + ") F(" + artifact.Filename + ")")
}

func (je JavaEndpoints) javaDownloadMetadataHandler(w http.ResponseWriter, r *http.Request) {

	artifact := RequestToArtifact(r)
	artifact.Filename = "maven-metadata.xml"

	fmt.Println("Metadata Download : JAVA, G(" + artifact.Group + ") A(" + artifact.Artifact + ") F(" + artifact.Filename + ")")
}

func (je JavaEndpoints) javaDownloadArtifactRouter(w http.ResponseWriter, r *http.Request) {

	ja := RequestToArtifact(r)

	fmt.Println("Artifact Download : JAVA, G(" + ja.Group + ") A(" + ja.Artifact + ") V(" + ja.Version + ") F(" + ja.Filename + ") T(" + ja.Type + ")")

	storage := filepath.Dir(`F:\Code\specialedge\storage\java\`)
	path := filepath.Join(storage, strings.Replace(ja.Group, ".", "/", -1)+"/", ja.Artifact, ja.Version, ja.Filename)

	// If file exists in index - attempt to serve,
	if je.Ind.IsDownloadedArtifact(ja.GetIdentifier(), ja.Type) {
		http.ServeFile(w, r, path)
	} else {
		// if it doesn't exist on filesystem, download and serve.
		javaDownloadAction(ja)

		// Add to the index
		je.Ind.AddArtifact(ja.GetIdentifier(), NewJavaFileList())
		je.Ind.AddDownloadedArtifact(ja.GetIdentifier(), ja.Type)

		// Serve the File to the User
		http.ServeFile(w, r, path)
	}
}

func javaDownloadAction(ja Artifact) {

	// create a custom client
	client := grab.NewClient()
	client.UserAgent = "Hangar v0.0.1"

	// Need to add this as a proxy configuration rather than having it hardcoded.
	mavenCentral := "https://repo.maven.apache.org/maven2/" + strings.Replace(ja.Group, ".", "/", -1) + "/" + ja.Artifact + "/" + ja.Version + "/" + ja.Filename
	storage := filepath.Dir(`F:\Code\specialedge\storage\java\`)
	path := filepath.Join(storage, strings.Replace(ja.Group, ".", "/", -1)+"/", ja.Artifact, ja.Version, ja.Filename)

	// Debug Statements
	//fmt.Println(mavenCentral)
	//fmt.Println(path)

	// create a download request
	req, err := grab.NewRequest(path, mavenCentral)
	if err != nil {
		panic(err)
	}

	// start download
	//fmt.Printf("Downloading %v...\n", req.URL())
	fmt.Println("Artifact Download : JAVA, G(" + ja.Group + ") A(" + ja.Artifact + ") V(" + ja.Version + ") F(" + ja.Filename + ")")
	resp := client.Do(req)
	//fmt.Printf("  %v\n", resp.HTTPResponse.Status)

	// set expected file size if known
	//req.Size = 1024

	// set expected file checksum if known
	//b, _ := hex.DecodeString("b982505fc48ea2221d163730c1856770dc6579af9eb73c997541c4ac6ecf20a9")
	//req.SetChecksum("sha256", b)

	// delete the downloaded file if it fails checksum validation
	//req.RemoveOnError = true

	// request a notification when the download is completed (successfully or
	// otherwise)

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("Download saved to ./%v \n", resp.Filename)
}
