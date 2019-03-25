package main

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// DatabaseCleanupHandler is awesome
func DatabaseCleanupHandler(w http.ResponseWriter, r *http.Request) {
	return
	// Load up a context and http client
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "Context loaded. Starting execution.")

	// // Make sure the request is from the appengine cron
	// if r.Header.Get("X-Appengine-Cron") == "" {
	// 	log.Warningf(ctx, "Request is not from the cron. Exiting")
	// 	w.WriteHeader(403)
	// 	return
	// }

	fbClient := GetFirebaseClient(ctx)
	if fbClient == nil {
		w.WriteHeader(500)
		return
	}
	defer fbClient.Close()

	// Get the list of sections we are actively tracking
	sectionsSnapshot := fbClient.Collection("tracked_sections").Documents(ctx)

	// Actually get all the data within these docs
	sectionDocuments, err := sectionsSnapshot.GetAll()
	if err != nil {
		log.Errorf(ctx, "Error getting tracked_sections! sec: %v", sectionDocuments)
		log.Errorf(ctx, "Error getting tracked_sections! Err: %v", err)
		w.WriteHeader(500)
		return
	}
	log.Debugf(ctx, "successfully got sectionDocuments we are tracking")

	fbBatch := fbClient.Batch()

	for _, doc := range sectionDocuments {
		data := doc.Data()

		data["crn"] = doc.Ref.ID
		fbBatch.Create(fbClient.Collection("sections_tracked").NewDoc(), data)
	}

	fbBatch.Commit(ctx)

}
