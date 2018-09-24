package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/appengine" // Required external App Engine library
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

// ScrapeSectionHandler scrapes
func ScrapeSectionHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	log.Infof(ctx, "Context loaded. Starting execution.")

	queryString := r.URL.Query()
	course := queryString["course"]
	dept := queryString["dept"]
	term := queryString["term"]

	if len(course) == 0 || len(dept) == 0 || len(term) == 0 {
		log.Errorf(ctx, "Malformed request to API")
		http.Error(w, "bad syntax. Missing params!", http.StatusBadRequest)
		return
	}
	log.Debugf(ctx, "term: %v", term)
	log.Debugf(ctx, "dept: %v", dept)
	log.Debugf(ctx, "course: %v", course)

	response, err := MakeAtlasSectionRequest(client, term[0], dept[0], course[0])

	if err != nil {
		log.Errorf(ctx, "Request to myInfo failed with error: %v", err)
		errorStr := fmt.Sprintf("Request to myInfo failed with error: %v", err)
		http.Error(w, errorStr, http.StatusInternalServerError)
		return
	}

	start := time.Now()

	sections, err := ParseSectionResponse(response, "")

	elapsed := time.Since(start)
	log.Infof(ctx, "Scrape time: %v", elapsed.String())

	if err != nil {
		log.Criticalf(ctx, "Course Scrape Failed with error: %v", err)
		errorStr := fmt.Sprintf("Course Scrape Failed with error: %v", err)
		http.Error(w, errorStr, http.StatusInternalServerError)
		return
	}

	//Set response headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Content-Type", "application/json")

	js, err := json.Marshal(sections)
	if err != nil {
		return
	}

	w.Write(js)
	response.Body.Close()
}
