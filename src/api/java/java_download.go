package java

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cavaliercoder/grab"
	"github.com/gorilla/mux"
)

// AppendJavaDownloadTopLevelMetadataRouter : In Java, Artifacts are saved with xml metadata at the artifact level as well as the version level
func AppendJavaDownloadTopLevelMetadataRouter(r *mux.Router) {
	// Hmm, still missing the type group :  (?:\\.){type:\\w*}
	r.HandleFunc("/{group:.+}/{artifact:.+}/maven-metadata.xml", javaDownloadMetadataHandler)
}

// AppendJavaDownloadArtifactRouter : This is the main endpoint for retrieving java artifacts through Hangar.
func AppendJavaDownloadArtifactRouter(r *mux.Router) {
	r.HandleFunc("/{group:.+}/{artifact:.+}/{version:.+}/{filename:[^/]+}", javaDownloadArtifactRouter)
}

func javaDownloadMetadataHandler(w http.ResponseWriter, r *http.Request) {

	artifact := RequestToArtifact(r)
	artifact.Filename = "maven-metadata.xml"

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Metadata Download : JAVA, G(" + artifact.Group + ") A(" + artifact.Artifact + ") F(" + artifact.Filename + ")")
}

func javaDownloadArtifactRouter(w http.ResponseWriter, r *http.Request) {

	artifact := RequestToArtifact(r)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Artifact Download : JAVA, G(" + artifact.Group + ") A(" + artifact.Artifact + ") V(" + artifact.Version + ") F(" + artifact.Filename + ")")
	javaDownloadAction(artifact)
}

func javaDownloadAction(ja Artifact) {

	// create a custom client
	client := grab.NewClient()
	client.UserAgent = "Hangar v0.0.1"

	// Need to add this as a proxy configuration rather than having it hardcoded.
	mavenCentral := "https://repo.maven.apache.org/maven2/" + strings.Replace(ja.Group, ".", "/", -1) + "/" + ja.Artifact + "/" + ja.Version + "/" + ja.Filename
	storage := filepath.Dir(`F:\Code\specialedge\storage\java\`)
	path := filepath.Join(storage, strings.Replace(ja.Group, ".", "/", -1)+"/", ja.Artifact, ja.Version, ja.Filename)

	fmt.Println(mavenCentral)
	fmt.Println(path)

	// create a download request
	req, err := grab.NewRequest(path, mavenCentral)
	if err != nil {
		panic(err)
	}

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)

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

	fmt.Printf("Download saved to ./%v \n", resp.Filename)
}
