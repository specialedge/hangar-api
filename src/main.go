package main

import (
	"net/http"

	"./api/healthcheck"
	"./api/java"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Initial admin & service endpoints
	r.HandleFunc("/healthcheck", healthcheck.HandlerHealthcheck)

	// Java Endpoints
	java.AppendJavaDownloadTopLevelMetadataRouter(r)
	java.AppendJavaDownloadArtifactRouter(r)

	http.ListenAndServe(":8080", r)
}
