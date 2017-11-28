package main

import (
	"net/http"

	"./api/healthcheck"
	"./api/java"
	"./index"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Create index to be used by all endpoints.
	ind := index.NewInMemory()

	// Initial admin & service endpoints
	r.HandleFunc("/healthcheck", healthcheck.HandlerHealthcheck)

	// Java Endpoints
	javaEndpoints := java.JavaEndpoints{
		Ind: ind,
	}

	javaEndpoints.AppendEndpoints(r)

	http.ListenAndServe(":8080", r)
}
