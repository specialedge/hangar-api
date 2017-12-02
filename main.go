package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/specialedge/hangar-api/api/healthcheck"
	"github.com/specialedge/hangar-api/api/java"
	"github.com/specialedge/hangar-api/index"
	"github.com/specialedge/hangar-api/storage"
)

func main() {

	// Create startup message to welcome the user.
	dat, err := ioutil.ReadFile("startup_message.txt")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))

	r := mux.NewRouter()

	// Create index to be used by all endpoints.
	ind := index.NewInMemory()

	// Create storage to be used by all endpoints.
	stor := storage.NewStorageLocal()

	// Initial admin & service endpoints
	r.HandleFunc("/healthcheck", healthcheck.HandlerHealthcheck)

	// Java Endpoints
	javaEndpoints := java.Endpoints{
		ArtifactIndex:   ind,
		ArtifactStorage: stor,
	}

	// Add all the endpoints for the Java API
	javaEndpoints.AppendEndpoints(r)

	http.ListenAndServe(":8080", r)
}
