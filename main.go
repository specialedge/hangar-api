package main

import (
	"net/http"

	figure "github.com/common-nighthawk/go-figure"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/specialedge/hangar-api/api/healthcheck"
	"github.com/specialedge/hangar-api/api/java"
	"github.com/specialedge/hangar-api/index"
	"github.com/specialedge/hangar-api/storage"
)

func main() {

	// Create startup message to welcome the user.
	startUpMessage := figure.NewFigure("hangar-api", "smslant", true)
	startUpMessage.Print()
	log.WithFields(log.Fields{
		"module": "main",
		"action": "PrintStartUpMessage",
	}).Info("Running")

	// Create a router and admin & service endpoints
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", healthcheck.HandlerHealthcheck)

	// Initalise the repo endpoints
	initialiseJavaEndpoints(r)

  	// Serve on 8080 with CORS support.
	http.ListenAndServe(":8080", handlers.CORS()(r))
}

func initialiseJavaEndpoints(r *mux.Router) {
	// Create index to be used by the Java endpoint.
	ind := index.NewInMemory()

	// Create storage to be used by the Java endpoint.
	stor := storage.NewStorageLocal()

	// Java Endpoints
	javaEndpoints := java.Endpoints{
		ArtifactIndex:   ind,
		ArtifactStorage: stor,
	}

	// Add all the endpoints for the Java API
	javaEndpoints.AppendEndpoints(r)
	
	// Load all the current files in the directory as Java Artifacts.
	javaEndpoints.ReIndex()
}