package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SpencerCornish/msubot-appspot/server/serverutils"
	log "github.com/sirupsen/logrus"
)

// HandleRequest scrapes
func HandleRequest(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	client := http.DefaultClient
	log.WithContext(ctx).Infof("Context loaded. Starting execution.")

	queryString := r.URL.Query()
	course := queryString["course"]
	dept := queryString["dept"]
	term := queryString["term"]

	if len(course) == 0 || len(dept) == 0 || len(term) == 0 {
		log.WithContext(ctx).Errorf("Malformed request to API")
		http.Error(w, "bad syntax. Missing params!", http.StatusBadRequest)
		return
	}
	log.WithContext(ctx).Debugf("term: %v", term)
	log.WithContext(ctx).Debugf("dept: %v", dept)
	log.WithContext(ctx).Debugf("course: %v", course)

	response, err := serverutils.MakeAtlasSectionRequest(client, term[0], dept[0], course[0])

	if err != nil {
		log.WithContext(ctx).Errorf("Request to myInfo failed with error: %v", err)
		errorStr := fmt.Sprintf("Request to myInfo failed with error: %v", err)
		http.Error(w, errorStr, http.StatusInternalServerError)
		return
	}

	start := time.Now()

	sections, err := serverutils.ParseSectionResponse(response, "")

	elapsed := time.Since(start)
	log.WithContext(ctx).Infof("Scrape time: %v", elapsed.String())

	if err != nil {
		log.WithError(err).WithContext(ctx).Errorf("Course Scrape Failed")
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
