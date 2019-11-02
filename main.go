package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/SpencerCornish/msubot-appspot/server"
	// _ "net/http/pprof"
)

func main() {
	http.HandleFunc("/sections", server.ScrapeSectionHandler)
	http.HandleFunc("/welcomeuser", server.WelcomeUserHandler)
	http.HandleFunc("/checktrackedsections", server.CheckSectionsHandler)
	http.HandleFunc("/prunesections", server.PruneSectionsHandler)
	// http.HandleFunc("/DatabaseCleanup", server.DatabaseCleanupHandler)
	http.HandleFunc("/receivemessage", server.ReceiveMessageHandler)
	http.HandleFunc("/healthcheck", server.HealthcheckHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
