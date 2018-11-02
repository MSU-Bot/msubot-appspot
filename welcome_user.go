package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

// WelcomeUserHandler sends the user their welcome text to MSUBot.
func WelcomeUserHandler(w http.ResponseWriter, r *http.Request) {
	//---------------------
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	defer r.Body.Close()
	log.Infof(ctx, "Context loaded. Starting execution.")

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

	queryString := r.URL.Query()
	rawphNum := queryString["number"]
	if len(rawphNum) == 0 {
		log.Warningf(ctx, "Incorrect number of args")
		w.WriteHeader(422)
		return
	}
	phNum := strings.Join(rawphNum, "")

	userData, uid := FetchUserDataWithNumber(ctx, fbClient, phNum)
	if userData == nil {
		log.Errorf(ctx, "User doesn't exist in the database. Userdata: %v", userData)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	welcomeSent, ok := userData["welcomeSent"].(bool)
	if !ok {
		log.Infof(ctx, "welcomeSent: %v", welcomeSent)
	}
	if welcomeSent {
		log.Infof(ctx, "Already welcomed user")
		w.WriteHeader(200)
		return
	}
	messageText := fmt.Sprintf("Thanks for signing up for MSUBot! We'll text you from this number when a seat opens up. Go Cats!")
	_, err = SendText(client, userData["number"].(string), messageText)
	if err != nil {
		log.Errorf(ctx, "Could not send text to user!")
		w.WriteHeader(500)
		return
	}
	fbClient.Collection("users").Doc(uid).Set(ctx, map[string]interface{}{
		"welcomeSent": true,
	}, firestore.MergeAll)

	w.WriteHeader(200)

}
