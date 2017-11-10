package main

import (
	"net/http"

	"./healthcheck"
	"./java"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Initial admin & service endpoints
	r.HandleFunc("/healthcheck", healthcheck.HandlerHealthcheck)

	// Java Endpoints
	java.AppendJavaDownloadMetadataRouter(r)
	http.ListenAndServe(":8080", r)
}
