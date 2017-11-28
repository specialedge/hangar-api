package java

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cavaliercoder/grab"
	"github.com/gorilla/mux"
)

func AppendJavaDownloadMetadataRouter(r *mux.Router) {
	// Hmm, still missing the type group :  (?:\\.){type:\\w*}
	r.HandleFunc("/{group:.+}/{artifact:.+}/maven-metadata.xml", javaDownloadMetadataHandler)
}

func javaDownloadMetadataHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"metadata": {"group": "`+vars["group"]+`", "artifact": "`+vars["artifact"]+`"}}`)
	javaDownloadAction()
}

func javaDownloadAction() {

	// create a custom client
	client := grab.NewClient()
	client.UserAgent = "Hangar v0.0.1"

	// Need to add this as a proxy configuration rather than having it hardcoded.
	mavenCentral := "https://repo.maven.apache.org/maven2/" + "activemq/activemq-core/3.2/activemq-core-3.2.jar"

	// create a download request
	req, err := grab.NewRequest(".", mavenCentral)
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
