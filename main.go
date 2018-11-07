package main

import (
	"net/http"
	// _ "net/http/pprof"

	"google.golang.org/appengine" // Required external App Engine library
)

func main() {
	http.HandleFunc("/sections", ScrapeSectionHandler)
	http.HandleFunc("/welcomeuser", WelcomeUserHandler)
	http.HandleFunc("/checktrackedsections", CheckSectionsHandler)
	http.HandleFunc("/prunesections", PruneSectionsHandler)
	http.HandleFunc("/DatabaseCleanup", DatabaseCleanupHandler)
	http.HandleFunc("/receivemessage", ReceiveMessageHandler)
	http.HandleFunc("/healthcheck", HealthcheckHandler)
	appengine.Main() // Starts the server to receive requests
}
