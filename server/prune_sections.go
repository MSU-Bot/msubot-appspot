package server

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// PruneSectionsHandler is run daily to clean up expired course checkers from old semesters.
func PruneSectionsHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log.WithContext(ctx).Infof("Context loaded. Starting execution.")

	// if r.Header.Get("X-Appengine-Cron") == "" {
	// 	log.Warningf(ctx, "Request is not from the cron! Exiting")
	// 	w.WriteHeader(403)
	// 	return
	// }

	fbClient := GetFirebaseClient(ctx)
	if fbClient == nil {
		w.WriteHeader(500)
		return
	}

	now := time.Now()
	term := "00"
	year := now.Year()
	if now.Month() > 10 {
		// If our current month is after October or greater, remove fall (year) and before
		term = fmt.Sprintf("%d%d", year, 70)

	} else if now.Month() > 8 {
		// If our current month is September or greater, we should remove summer (year) and before
		term = fmt.Sprintf("%d%d", year, 50)

	} else if now.Month() > 3 {
		// If our current month is April or greater, remove spring (year) and before
		term = fmt.Sprintf("%d%d", year, 30)
	}

	log.WithContext(ctx).Infof("Removing all trackedsections where term <= %s", term)
	// Get the list of sections we are actively tracking
	sectionsSnapshot := fbClient.Collection("sections_tracked").Where("term", "<=", term).Documents(ctx)

	expiredSectionDocs, err := sectionsSnapshot.GetAll()
	if err != nil {
		log.WithContext(ctx).Errorf("Unable to retrieve expired sections!")
		w.WriteHeader(500)
		return
	}

	log.WithContext(ctx).Infof("Number of expired courses: %d", len(expiredSectionDocs))
	for _, doc := range expiredSectionDocs {
		data := doc.Data()

		// Type conversions
		term, ok := data["term"].(string)
		if !ok {
			log.WithContext(ctx).Errorf("type conv failed for courseNumber. Could not prune %s", doc.Ref.ID)
			continue
		}
		// Type conversions
		crn, ok := data["crn"].(string)
		if !ok {
			log.WithContext(ctx).Errorf("type conv failed for crn. Could not prune %s", doc.Ref.ID)
			continue
		}

		log.WithContext(ctx).Infof("Moving Expired doc with uid %s because term was %v ", doc.Ref.ID, data["term"])

		err := MoveTrackedSection(ctx, fbClient, crn, doc.Ref.ID, term)
		if err != nil {
			log.WithContext(ctx).Errorf("Unable to move doc with UID: %s", doc.Ref.ID)
		}
	}
}
