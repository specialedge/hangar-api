package java

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// This allows us to download the top level metadata - which won't have a version.
// Example Path : /com/spedge/hangar-artifact/maven-metadata.xml
func TestJavaDownloadMetadataHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/com/spedge/hangar-artifact/maven-metadata.xml", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	AppendJavaDownloadMetadataRouter(r)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	r.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"metadata": {"group": "com/spedge", "artifact": "hangar-artifact"}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
