package java

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

func TestAppendEndpoints(t *testing.T) {

	// Create Endpoints object
	javaEndpoints := Endpoints{
		ArtifactIndex:   nil,
		ArtifactStorage: nil,
	}

	// Override the MUX handler functions that work for this test rather than
	// heading down a path of randomness and tomfoolery.
	artifactFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	checksumFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}

	r := mux.NewRouter()
	javaEndpoints.AppendEndpoints(r, artifactFunc, checksumFunc)

	// Should be picked up by artifactFunction
	runRequest("/java/com/specialedge/hangar-api/1.2.3/hangar-api-1.2.3.jar", 200, r, t)
	runRequest("/java/com/specialedge/hangar-api/1.2.3/hangar-api-1.2.3.pom", 200, r, t)

	// Should be picked up by artifactChecksumFunction
	runRequest("/java/com/specialedge/hangar-api/1.2.3/hangar-api-1.2.3.jar.sha1", 201, r, t)
	runRequest("/java/com/specialedge/hangar-api/1.2.3/hangar-api-1.2.3.jar.md5", 201, r, t)

	runRequest("/java/com/", 404, r, t)
	runRequest("/java/com/specialedge", 404, r, t)
	runRequest("/java/com/specialedge/hangar-api", 404, r, t)
	runRequest("/java/com/specialedge/hangar-api/1.2.3", 404, r, t)
	runRequest("/java/com/specialedge/hangar-api/1.2.3/hangar-api-1.2.3", 404, r, t)
	runRequest("/java/com/specialedge/hangar-api/1.2.3/hangar-api-1.2.3.jar.fish", 404, r, t)
	runRequest("/java/com/specialedge/hangar-api/1.2.3/hangar-api-1.2.3.jar.fish", 404, r, t)
	runRequest("/java/com/specialedge/hangar-api/1.2.3/hangar-api-1.2.3.jar.sha", 404, r, t)
	runRequest("/java/com/specialedge/hangar-api/1.2.3/hangar-api-1.2.3.jar.sha1.fish", 404, r, t)
}

func runRequest(path string, codeExpected int, r *mux.Router, t *testing.T) {
	req, _ := http.NewRequest("GET", path, nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	if res.Code != codeExpected {
		t.Error("Expected "+strconv.Itoa(codeExpected)+"but got ", res.Code)
	}
}
