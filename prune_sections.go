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
	return

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
	if now.Month() == 12 {
		// If our current month is December, we should remove fall (year) and before
		term = fmt.Sprintf("%d%d", year, 70)
	} else if now.Month() > 8 {
		// If our current month is after August, we should remove summer (year) and before
		term = fmt.Sprintf("%d%d", year, 50)

	} else if now.Month() > 4 {
		// If our current month is after April, remove spring (year) and before
		term = fmt.Sprintf("%d%d", year, 30)

	}
	log.Infof(ctx, "!! removing all where term <= %s", term)
	// Get the list of sections we are actively tracking
	sectionsSnapshot := fbClient.Collection("sections_tracked").Where("term", "<=", term).Documents(ctx)

	expiredSectionDocs, err := sectionsSnapshot.GetAll()
	if err != nil {
		log.Errorf(ctx, "Unable to retrieve expired sections!")
	}
	log.Infof(ctx, "Number of expired courses: %d", len(expiredSectionDocs))
	for i, doc := range expiredSectionDocs {
		data := doc.Data()
		log.Infof(ctx, "%d - TERM:%s  FULL: %v", i, data["term"], data)
	}
}
