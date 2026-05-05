package cleanup

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// MigrateDatabase is awesome
func MigrateDatabase(w http.ResponseWriter, r *http.Request) {

	// Load up a context and http client
	ctx := r.Context()
	log.WithContext(ctx).Infof("Context loaded. Starting execution.")

	log.WithContext(ctx).Info("No Migrations to run, exiting :)")
	return
	// fbClient := serverutils.GetFirebaseClient(ctx)
	// if fbClient == nil {
	// 	w.WriteHeader(500)
	// 	return
	// }
	// defer fbClient.Close()

	// // Get the list of sections we are actively tracking
	// sectionsSnapshot := fbClient.Collection("tracked_sections").Documents(ctx)

	// // Actually get all the data within these docs
	// sectionDocuments, err := sectionsSnapshot.GetAll()
	// if err != nil {
	// 	log.WithContext(ctx).Errorf("Error getting tracked_sections! sec: %v", sectionDocuments)
	// 	log.WithContext(ctx).Errorf("Error getting tracked_sections! Err: %v", err)
	// 	w.WriteHeader(500)
	// 	return
	// }
	// log.WithContext(ctx).Debugf("successfully got sectionDocuments we are tracking")

	// fbBatch := fbClient.Batch()

	// for _, doc := range sectionDocuments {
	// 	data := doc.Data()

	// 	data["crn"] = doc.Ref.ID
	// 	fbBatch.Create(fbClient.Collection("sections_tracked").NewDoc(), data)
	// }

	// fbBatch.Commit(ctx)

}
