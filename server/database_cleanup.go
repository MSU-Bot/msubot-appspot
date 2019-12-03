package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// DatabaseCleanupHandler is awesome
func DatabaseCleanupHandler(w http.ResponseWriter, r *http.Request) {
	return
	// Load up a context and http client
	ctx := r.Context()
	log.WithContext(ctx).Infof("Context loaded. Starting execution.")

	// // Make sure the request is from the appengine cron
	// if r.Header.Get("X-Appengine-Cron") == "" {
	// 	log.WithContext(ctx).Warningf("Request is not from the cron. Exiting")
	// 	w.WriteHeader(403)
	// 	return
	// }

	fbClient := GetFirebaseClient(ctx)
	if fbClient == nil {
		w.WriteHeader(500)
		return
	}

	// Get the list of sections we are actively tracking
	sectionsSnapshot := fbClient.Collection("tracked_sections").Documents(ctx)

	// Actually get all the data within these docs
	sectionDocuments, err := sectionsSnapshot.GetAll()
	if err != nil {
		log.WithContext(ctx).Errorf("Error getting tracked_sections! sec: %v", sectionDocuments)
		log.WithContext(ctx).Errorf("Error getting tracked_sections! Err: %v", err)
		w.WriteHeader(500)
		return
	}
	log.WithContext(ctx).Debugf("successfully got sectionDocuments we are tracking")

	fbBatch := fbClient.Batch()

	for _, doc := range sectionDocuments {
		data := doc.Data()

		data["crn"] = doc.Ref.ID
		fbBatch.Create(fbClient.Collection("sections_tracked").NewDoc(), data)
	}

	fbBatch.Commit(ctx)

}
