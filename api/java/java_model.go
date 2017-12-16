package java

import (
	"path/filepath"
	"strings"

	"github.com/specialedge/hangar-api/index"
	"github.com/specialedge/hangar-api/storage"
)

// Artifact : Represents the metadata for a Java Artifact
type Artifact struct {
	Group        string `json:"group"`
	Artifact     string `json:"artifact"`
	Version      string `json:"version"`
	Filename     string `json:"filename"`
	Type         string `json:"type"`
	ChecksumType string `json:"checksumType"`
}

// NewJavaFileList : Generates a new list of acceptable files to store for this
// type of endpoint.
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
		Group:        strings.Replace(vars["group"], "/", ".", -1),
		Artifact:     vars["artifact"],
		Version:      vars["version"],
		Filename:     vars["filename"] + vars["type"] + vars["checksumType"],
		Type:         strings.Replace(vars["type"], ".", "", -1) + vars["checksumType"],
		ChecksumType: strings.Replace(vars["checksumType"], ".", "", -1),
	}
}

// StorageIdentifierToArtifact : Convert a storage identifier into an Artifact.
func StorageIdentifierToArtifact(id storage.Identifier) (ja Artifact) {

	// Remove the initial slash if there is one (by mistake)
	sanitized := strings.TrimPrefix(id.Key, id.Separator)

	// Split the path into a set of
	a := strings.Split(sanitized, id.Separator)

	// General GAV parameters
	filename, a := a[len(a)-1], a[:len(a)-1]
	version, a := a[len(a)-1], a[:len(a)-1]
	artifact, a := a[len(a)-1], a[:len(a)-1]
	group := strings.Join(a, ".")

	// Attempted this using string manipulation but couldn't - got scuppered
	// by some weird version numbers (like 1 which got in the way of sha1)
	// This seems a more robust way.
	typeVar := ""
	checksum := ""

	idents := NewJavaFileList().FileTypes

	for filetype := range idents {
		if strings.HasSuffix(filename, filetype) {
			checksumSlice := strings.Split(filetype, ".")
			typeVar = checksumSlice[0]
			if len(checksumSlice) > 1 {
				typeVar = typeVar + "." + checksumSlice[1]
				checksum = checksumSlice[1]
			}
			break
		}
	}

	return Artifact{
		Group:        group,
		Artifact:     artifact,
		Version:      version,
		Filename:     filename,
		Type:         typeVar,
		ChecksumType: checksum,
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
		Key:       filepath.Join(strings.Replace(a.Group, ".", "/", -1)+"/", a.Artifact, a.Version, a.Filename),
		Separator: "/",
	}
}

// ToString : Prints out the Identifier in a easy to understand format.
func (a Artifact) ToString() string {
	output := "G(" + a.Group + ") A(" + a.Artifact + ") V(" + a.Version + ") F(" + a.Filename + ") T(" + a.Type + ")"
	if len(a.ChecksumType) > 0 {
		output += " C(" + a.ChecksumType + ")"
	}
	return output
}
