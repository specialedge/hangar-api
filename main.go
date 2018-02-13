package main

import (
	"fmt"
	"net/http"

	figure "github.com/common-nighthawk/go-figure"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/specialedge/hangar-api/api/healthcheck"
	"github.com/specialedge/hangar-api/api/java"
	"github.com/spf13/viper"
)

func main() {

	//log.SetLevel(log.DebugLevel)

	// Create startup message to welcome the user.
	startUpMessage := figure.NewFigure("hangar-api", "smslant", true)
	startUpMessage.Print()

	// Print event to show the system is starting
	log.WithFields(log.Fields{
		"module": "main",
		"action": "PrintStartUpMessage",
	}).Info("Starting")

	// Read in the configuration file
	readConfiguration()

	// Create a router and admin & service endpoints
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", healthcheck.HandlerHealthcheck)

	// Initialise the repo endpoints
	if viper.IsSet("java") {

		java.InitialiseJavaEndpoints(r)

		log.WithFields(log.Fields{
			"module": "main",
			"action": "JavaEnabled",
		}).Info("Finished creating Java Endpoint")
	}

	// Serve on 8080 with CORS support.
	http.ListenAndServe(":8080", handlers.CORS()(r))
}

func readConfiguration() {

	// Without a configuration file, we don't want to start the system.
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()

	if err != nil {
		log.WithFields(log.Fields{
			"module": "main",
			"action": "readConfiguration",
		}).Info("Failed to read configuration file : %s", err)

		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	log.WithFields(log.Fields{
		"module": "main",
		"action": "readConfiguration",
	}).Info("Configuration parsed...")
}
