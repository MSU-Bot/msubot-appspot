package main

import (
	"fmt"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// PruneSectionsHandler is run daily to clean up expired course checkers from old semesters.
func PruneSectionsHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	log.Infof(ctx, "Context loaded. Starting execution.")

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
	defer fbClient.Close()

	now := time.Now()
	term := "00"
	year := now.Year()

	if now.Month() > 8 {
		// If our current month is September or greater, we should remove summer (year) and before
		term = fmt.Sprintf("%d%d", year, 50)

	} else if now.Month() > 4 {
		// If our current month is after May or greater, remove spring (year) and before
		term = fmt.Sprintf("%d%d", year, 30)
	} else if now.Month() > 0 {
		// If our current month is after January or greater, remove fall (year-1) and before
		term = fmt.Sprintf("%d%d", year-1, 70)

	}
	log.Infof(ctx, "Removing all trackedsections where term <= %s", term)
	// Get the list of sections we are actively tracking
	sectionsSnapshot := fbClient.Collection("sections_tracked").Where("term", "<=", term).Documents(ctx)

	expiredSectionDocs, err := sectionsSnapshot.GetAll()
	if err != nil {
		log.Errorf(ctx, "Unable to retrieve expired sections!")
		w.WriteHeader(500)
		return
	}

	log.Infof(ctx, "Number of expired courses: %d", len(expiredSectionDocs))
	for _, doc := range expiredSectionDocs {
		data := doc.Data()

		// Type conversions
		term, ok := data["term"].(string)
		if !ok {
			log.Errorf(ctx, "type conv failed for courseNumber. Could not prune %s", doc.Ref.ID)
			continue
		}
		// Type conversions
		crn, ok := data["crn"].(string)
		if !ok {
			log.Errorf(ctx, "type conv failed for crn. Could not prune %s", doc.Ref.ID)
			continue
		}

		log.Debugf(ctx, "Moving Expired doc with uid %s because term was %v ", doc.Ref.ID, data["term"])

		err := MoveTrackedSection(ctx, fbClient, crn, doc.Ref.ID, term)
		if err != nil {
			log.Errorf(ctx, "Unable to move doc with UID: %s", doc.Ref.ID)
		}
	}
}
