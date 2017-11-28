package java

import (
	"net/http"
	"strings"

	"../../index"
	"github.com/gorilla/mux"
)

// Artifact : Represents the metadata for a Java Artifact
type Artifact struct {
	Group    string `json:"group"`
	Artifact string `json:"artifact"`
	Version  string `json:"version"`
	Filename string `json:"filename"`
	Type     string `json:"type"`
}

func NewJavaFileList() index.FileList {

	types := map[string]bool{
		"pom":  false,
		"jar":  false,
		"sha1": false,
		"md5":  false,
	}

	return index.FileList{
		FileTypes: types,
	}
}

// RequestToArtifact : Convert a map of strings into an Artifact.
func RequestToArtifact(r *http.Request) (ja Artifact) {

	// Grab the variables from the request.
	vars := mux.Vars(r)

	// The variables for the group might be slash-delimited, we need them
	// to be dot-delimited to be accurate in Java terminology.
	return Artifact{
		Group:    strings.Replace(vars["group"], "/", ".", -1),
		Artifact: vars["artifact"],
		Version:  vars["version"],
		Filename: vars["filename"] + vars["type"],
		Type:     strings.Replace(vars["type"], ".", "", -1),
	}
}

// GetIdentifier : Return a unique key for this artifact to identify it by
func (a Artifact) GetIdentifier() index.Identifier {
	return index.Identifier{
		Key: strings.Join([]string{"JAVA", a.Group, a.Artifact, a.Version}, ":"),
	}
}
