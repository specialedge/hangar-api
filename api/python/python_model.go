package python

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/specialedge/hangar-api/index"
	"github.com/specialedge/hangar-api/storage"
)

// Artifact : Represents the metadata for a Python Artifact
type Artifact struct {
	PackageName string `json:"packageName"`
	Version     string `json:"package"`
	Letter      string `json:"letter"`
	PackageType string `json:"packageType"`
	PackageFile string `json:"packageFile"`
}

// NewPythonFileList : Generates a new list of acceptable files to store for this
// type of endpoint.
func NewPythonFileList() index.FileList {

	types := map[string]bool{
		"egg":    false,
		"whl":    false,
		"tar.gz": false,
		"zip":    false,
		"exe":    false,
		"msi":    false,
	}

	return index.FileList{
		FileTypes: types,
	}
}

// RequestToArtifact : Convert a map of strings into an Artifact.
func RequestToArtifact(vars map[string]string) (pa Artifact) {
	return Artifact{
		PackageName: vars["packageName"],
		Version:     vars["version"],
		Letter:      string(vars["packageName"][0]),
		PackageType: vars["packageType"],
		PackageFile: vars["packageFile"],
	}
}

// RequestToVersions should return all the available versions (INCLUDING filetype) for a given package
func RequestToVersions(pakcageName string) []string {
	versions := []string{"0.1", "0.2"}
	return versions
}

// GetFilename constructs the full filename from a package's name and available version
func GetFilename(version string, packageName string) string {
	packageFilename := fmt.Sprintf("%s-%s", packageName, version)
	return packageFilename
}

// StorageIdentifierToArtifact : Convert a storage identifier into a Python Artifact.
func StorageIdentifierToArtifact(id storage.Identifier) (pa Artifact) {

	// Remove the initial slash if there is one (by mistake)
	sanitized := strings.TrimPrefix(id.Key, id.Separator)

	// Split the path into a set of
	a := strings.Split(sanitized, id.Separator)

	// General GAV parameters
	version, a := a[len(a)-1], a[:len(a)-1]
	packageName := "D"
	letter := string(packageName[0])
	packageType := "source"
	packageFile := "source.tar.gz"

	return Artifact{
		PackageName: packageName,
		Version:     version,
		Letter:      letter,
		PackageType: packageType,
		PackageFile: packageFile,
	}
}

// GetIdentifier : Return a unique key for this artifact to identify it by
func (a Artifact) GetIdentifier() index.Identifier {
	return index.Identifier{
		Key: strings.Join([]string{"PYTHON", a.PackageType, a.PackageName, a.Version}, ":"),
	}
}

// GetStorageIdentifier : At this point, we return a slash-delimited "path" for the Artifact
func (a Artifact) GetStorageIdentifier() storage.Identifier {
	return storage.Identifier{
		Key:       filepath.Join("/", a.PackageName, a.Letter, a.Version, a.PackageFile),
		Separator: "/",
	}
}

// ToString : Prints out the Identifier in a easy to understand format.
func (a Artifact) ToString() string {
	output := "P(" + a.PackageName + ") L(" + a.Letter + ") V(" + a.Version + ") F(" + a.PackageFile + ") T(" + a.PackageType + ")"
	return output
}
