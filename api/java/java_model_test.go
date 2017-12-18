package java

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/specialedge/hangar-api/storage"
)

func createArtifact() Artifact {
	return Artifact{
		Group:        "com.specialedge.hangar",
		Artifact:     "test-artifact",
		Version:      "1.2.3",
		Filename:     "test-artifact-1.2.3.jar",
		Type:         "jar",
		ChecksumType: "",
	}
}

func createArtifactChecksum() Artifact {
	return Artifact{
		Group:        "com.specialedge.hangar",
		Artifact:     "test-artifact",
		Version:      "1.2.3",
		Filename:     "test-artifact-1.2.3.pom.sha1",
		Type:         "pom.sha1",
		ChecksumType: "sha1",
	}
}

func TestToString(t *testing.T) {

	a := createArtifact()
	result := "G(com.specialedge.hangar) A(test-artifact) V(1.2.3) F(test-artifact-1.2.3.jar) T(jar)"

	if strings.Compare(a.ToString(), result) != 0 {
		t.Error("ToString for Artifact is incorrect :" + a.ToString())
	}

	a = createArtifactChecksum()
	result = "G(com.specialedge.hangar) A(test-artifact) V(1.2.3) F(test-artifact-1.2.3.pom.sha1) T(pom.sha1) C(sha1)"

	if strings.Compare(a.ToString(), result) != 0 {
		t.Error("ToString for Artifact is incorrect :" + a.ToString())
	}
}

func TestIdentifier(t *testing.T) {

	a := createArtifact()
	result := "JAVA:com.specialedge.hangar:test-artifact:1.2.3"

	if strings.Compare(a.GetIdentifier().Key, result) != 0 {
		t.Error("Identifier Key for Artifact is incorrect" + a.GetIdentifier().Key)
	}

	a = createArtifactChecksum()
	result = "JAVA:com.specialedge.hangar:test-artifact:1.2.3"

	if strings.Compare(a.GetIdentifier().Key, result) != 0 {
		t.Error("Identifier Key for Artifact is incorrect" + a.GetIdentifier().Key)
	}
}

func TestStorageIdentifier(t *testing.T) {

	a := createArtifact()
	result := filepath.Join("com", "specialedge", "hangar", "test-artifact", "1.2.3", "test-artifact-1.2.3.jar")

	if strings.Compare(a.GetStorageIdentifier().Key, result) != 0 {
		t.Error("Storage Key for Artifact is incorrect" + a.GetStorageIdentifier().Key)
	}

	a = createArtifactChecksum()
	result = filepath.Join("com", "specialedge", "hangar", "test-artifact", "1.2.3", "test-artifact-1.2.3.pom.sha1")

	if strings.Compare(a.GetStorageIdentifier().Key, result) != 0 {
		t.Error("Storage Key for Artifact is incorrect" + a.GetStorageIdentifier().Key)
	}
}

func TestRequestToArtifact(t *testing.T) {

	// I'd have rather passed in the whole request and have mux sort it out
	// but I couldn't determine how to pass in context. Seems it's quite laborious but there's a Issue raised
	// with mux to fix : https://github.com/gorilla/mux/issues/233

	context := make(map[string]string)

	context["group"] = "com/specialedge/hangar"
	context["artifact"] = "test-artifact"
	context["version"] = "1.2.3"
	context["filename"] = "test-artifact-1.2.3"
	context["type"] = ".jar"

	a := RequestToArtifact(context)
	result := "G(com.specialedge.hangar) A(test-artifact) V(1.2.3) F(test-artifact-1.2.3.jar) T(jar)"

	if strings.Compare(a.ToString(), result) != 0 {
		t.Error("Artifact has not been generated from Request : " + a.ToString())
	}
}

func TestStorageIdentifierToArtifact(t *testing.T) {

	storageReq := storage.Identifier{
		Key:       "com/specialedge/hangar/test-artifact/1.2.3/test-artifact-1.2.3.pom.sha1",
		Separator: "/",
	}

	a := StorageIdentifierToArtifact(storageReq)
	result := "G(com.specialedge.hangar) A(test-artifact) V(1.2.3) F(test-artifact-1.2.3.pom.sha1) T(pom.sha1) C(sha1)"

	if strings.Compare(a.ToString(), result) != 0 {
		t.Error("Artifact has not been generated from Request : " + a.ToString())
	}

	storageReq = storage.Identifier{
		Key:       "\\com\\specialedge\\hangar\\test-artifact\\1.2.3\\test-artifact-1.2.3.jar",
		Separator: "\\",
	}

	a = StorageIdentifierToArtifact(storageReq)
	result = "G(com.specialedge.hangar) A(test-artifact) V(1.2.3) F(test-artifact-1.2.3.jar) T(jar)"

	if strings.Compare(a.ToString(), result) != 0 {
		t.Error("Artifact has not been generated from Request : " + a.ToString())
	}
}
