package server

import (
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// AddUserToSectionHandler tba
func AddUserToSectionHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	log.Infof(ctx, "Context loaded. Starting execution.")

	requestBody, err := r.GetBody()
	if err != nil || requestBody == nil {
		log.Errorf(ctx, "Could not get request body: %v", err)
		w.WriteHeader(500)
		return
	}
	defer r.Body.Close()

	queryString := r.URL.Query()

	course := queryString["course"]
	dept := queryString["dept"]
	term := queryString["term"]

	//Set response headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")

	firebasePID := os.Getenv("FIREBASE_PROJECT")
	log.Debugf(ctx, "Loaded firebase project ID: %v", firebasePID)
	if firebasePID == "" {
		log.Criticalf(ctx, "Firebase Project ID is nil, I cannot continue.")
		panic("Firebase Project ID is nil")
	}

	fbClient, err := firestore.NewClient(ctx, firebasePID)
	defer fbClient.Close()
	if err != nil {
		log.Errorf(ctx, "Could not create new client for Firebase %v", err)
		w.WriteHeader(500)
		return
	}

}
