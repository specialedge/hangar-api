package main

import (
	"net/http"

	"./api/healthcheck"
	"./api/java"
	"./index"
	"./storage"
	"github.com/gorilla/mux"
)

func main() {
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
