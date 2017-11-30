package java

import (
	"path/filepath"
	"strings"

	"../../index"
	"../../storage"
)

// Artifact : Represents the metadata for a Java Artifact
type Artifact struct {
	Group    string `json:"group"`
	Artifact string `json:"artifact"`
	Version  string `json:"version"`
	Filename string `json:"filename"`
	Type     string `json:"type"`
	Checksum string `json:"checksum"`
}

func NewJavaFileList() index.FileList {

	types := map[string]bool{
		"pom":      false,
		"pom.sha1": false,
		"pom.md5":  false,
		"jar":      false,
		"jar.sha1": false,
		"jar.md5":  false,
	}

	return index.FileList{
		FileTypes: types,
	}
}

// RequestToArtifact : Convert a map of strings into an Artifact.
func RequestToArtifact(vars map[string]string) (ja Artifact) {

	// The variables for the group might be slash-delimited, we need them
	// to be dot-delimited to be accurate in Java terminology.
	return Artifact{
		Group:    strings.Replace(vars["group"], "/", ".", -1),
		Artifact: vars["artifact"],
		Version:  vars["version"],
		Filename: vars["filename"] + vars["type"] + vars["checksum"],
		Type:     strings.Replace(vars["type"], ".", "", -1) + vars["checksum"],
		Checksum: strings.Replace(vars["checksum"], ".", "", -1),
	}
}

// GetIdentifier : Return a unique key for this artifact to identify it by
func (a Artifact) GetIdentifier() index.Identifier {
	return index.Identifier{
		Key: strings.Join([]string{"JAVA", a.Group, a.Artifact, a.Version}, ":"),
	}
}

// GetStorageIdentifier : At this point, we return a slash-delimited "path" for the Artifact
func (a Artifact) GetStorageIdentifier() storage.Identifier {
	return storage.Identifier{
		Key: filepath.Join(strings.Replace(a.Group, ".", "/", -1)+"/", a.Artifact, a.Version, a.Filename),
	}
}

// ToString : Prints out the Identifier in a easy to understand format.
func (a Artifact) ToString() string {
	output := "G(" + a.Group + ") A(" + a.Artifact + ") V(" + a.Version + ") F(" + a.Filename + ") T(" + a.Type + ")"
	if len(a.Checksum) > 0 {
		output += ", C(" + a.Checksum + ")"
	}
	return output
}
