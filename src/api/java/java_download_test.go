package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// This allows us to download the top level metadata - which won't have a version.
// Example Path : /java/com/spedge/hangar-artifact/maven-metadata.xml
func TestJavaDownloadMetadataHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/java/com/spedge/hangar-artifact/maven-metadata.xml", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(javaDownloadMetadataHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"metadata": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
