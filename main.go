package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/MSU-Bot/msubot-appspot/server"
	"github.com/MSU-Bot/msubot-appspot/server/checksections"
	"github.com/MSU-Bot/msubot-appspot/server/healthcheck"
	"github.com/MSU-Bot/msubot-appspot/server/messenger"
	"github.com/MSU-Bot/msubot-appspot/server/pruner"
	"github.com/MSU-Bot/msubot-appspot/server/scraper"
)

func main() {
	http.HandleFunc("/sections", scraper.HandleRequest)
	http.HandleFunc("/welcomeuser", server.WelcomeUserHandler)
	http.HandleFunc("/checktrackedsections", checksections.HandleRequest)
	http.HandleFunc("/prunesections", pruner.HandleRequest)
	http.HandleFunc("/receivemessage", messenger.RecieveMessage)
	http.HandleFunc("/healthcheck", healthcheck.CheckHealth)
	http.HandleFunc("/updatedepartments", scraper.HandleDepartmentRequest)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
