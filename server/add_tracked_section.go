package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type addSectionRequest struct {
	Token  string `json:"token"`
	UserID string `json:"userid"`
	CRNs   []int  `json:"crns"`
}

// AddTrackedSectionHandler tba
func AddTrackedSectionHandler(w http.ResponseWriter, r *http.Request) {

	//Set response headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST")

	// Load up a context
	ctx := r.Context()
	log.WithContext(ctx).Infof("Context loaded. Starting execution.")

	req, err := parseAndFormatRequest(r.Body)
	if err != nil {
		log.WithContext(ctx).WithError(err).Warning("Could not read request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := ValidateUserClaims(ctx, req.Token, req.UserID)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Could not validate user")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fbClient := GetFirebaseClient(ctx)

	fireclient
	if fbClient == nil {
		log.WithContext(ctx).Warning("Unable to get firebase client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func parseAndFormatRequest(body io.ReadCloser) (*addSectionRequest, error) {
	respStruct := new(addSectionRequest)
	requestBody, err := ioutil.ReadAll(body)
	if err != nil {
		return respStruct, err
	}

	err = json.Unmarshal(requestBody, &respStruct)
	if err != nil {
		return respStruct, err
	}

	return respStruct, nil
}
