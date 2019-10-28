package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
)

// WelcomeUserHandler sends the user their welcome text to MSUBot.
func WelcomeUserHandler(w http.ResponseWriter, r *http.Request) {
	//---------------------
	ctx := r.Context()
	client := http.DefaultClient
	defer r.Body.Close()
	log.WithContext(ctx).Infof("Context loaded. Starting execution.")

	//Set response headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")

	firebasePID := os.Getenv("FIREBASE_PROJECT")
	log.WithContext(ctx).Infof("Loaded firebase project ID: %v", firebasePID)
	if firebasePID == "" {
		log.WithContext(ctx).Error("Firebase Project ID is nil, I cannot continue.")
		panic("Firebase Project ID is nil")
	}

	fbClient, err := firestore.NewClient(ctx, firebasePID)
	defer fbClient.Close()
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Could not create new client for Firebase")
		w.WriteHeader(500)
		return
	}

	queryString := r.URL.Query()
	rawphNum := queryString["number"]
	if len(rawphNum) == 0 {
		log.WithContext(ctx).Error("Incorrect number of args")
		w.WriteHeader(422)
		return
	}
	phNum := strings.Join(rawphNum, "")

	userData, uid := FetchUserDataWithNumber(ctx, fbClient, phNum)
	if userData == nil {
		log.WithContext(ctx).Errorf("User doesn't exist in the database. Userdata: %v", userData)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	welcomeSent, ok := userData["welcomeSent"].(bool)
	if !ok {
		log.WithContext(ctx).Infof("welcomeSent: %v", welcomeSent)
	}
	if welcomeSent {
		log.WithContext(ctx).Infof("Already welcomed user")
		w.WriteHeader(200)
		return
	}
	messageText := fmt.Sprintf("Thanks for signing up for MSUBot! We'll text you from this number when a seat opens up. Go Cats!")
	_, err = SendText(client, userData["number"].(string), messageText)
	if err != nil {
		log.WithContext(ctx).Errorf("Could not send text to user!")
		w.WriteHeader(500)
		return
	}
	fbClient.Collection("users").Doc(uid).Set(ctx, map[string]interface{}{
		"welcomeSent": true,
	}, firestore.MergeAll)

	w.WriteHeader(200)

}
