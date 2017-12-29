package python

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

func TestAppendEndpoints(t *testing.T) {

	// Create Endpoints object
	pythonEndpoints := endpoints{
		ArtifactIndex:   nil,
		ArtifactStorage: nil,
	}

	// Override the MUX handler functions that work for this test
	artifactFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	versionsFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	r := mux.NewRouter()
	pythonEndpoints.AppendEndpoints(r, artifactFunc, versionsFunc)

	// Should be picked up by artifactFunction
	runRequest("/python/packages/source/h/hangar-api/hangar-api-1.2.3.tar.gz", 200, r, t)
	runRequest("/python/packages/source/h/hangar-api/hangar-api-1.2.3.zip", 200, r, t)
	runRequest("/python/hangar-api/", 200, r, t)

	runRequest("/python/packages/source/h/hangar-api/hangar-api-1.2.3", 404, r, t)
	runRequest("/python/packages/source/h/hangar-api/hangar-api-1.2.3.tar", 404, r, t)
}

func runRequest(path string, codeExpected int, r *mux.Router, t *testing.T) {
	req, _ := http.NewRequest("GET", path, nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	if res.Code != codeExpected {
		t.Error("Expected "+strconv.Itoa(codeExpected)+" but got ", res.Code, " for ", path)
	}
}
