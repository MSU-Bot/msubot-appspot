package main

import (
	"net/http"
)

// PruneSectionsHandler is run daily to clean up expired course checkers from old semesters.
func PruneSectionsHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := appengine.NewContext(r)
	// client := urlfetch.Client(ctx)
	// log.Infof(ctx, "Context loaded. Starting execution.")

	// if r.Header.Get("X-Appengine-Cron") == "" {
	// 	log.Warningf(ctx, "Request is not from the cron! Exiting")
	// 	w.WriteHeader(403)
	// 	return
	// }

}
