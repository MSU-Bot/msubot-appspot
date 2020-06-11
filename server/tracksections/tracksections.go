package tracksections

import (
	"encoding/json"
	"net/http"

	"github.com/SpencerCornish/msubot-appspot/server/mauth"
	log "github.com/sirupsen/logrus"
)

type trackRequest struct {
	Crns []int
	Term int
}

// HandleRequest scrapes
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.WithContext(ctx).Infof("Context loaded. Starting execution.")

	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}

	requestDecoder := json.NewDecoder(r.Body)
	requestDecoder.DisallowUnknownFields()

	authToken, err := mauth.VerifyToken(ctx, "TOKEN")
	if err != nil {
		http.Error(w, "Token is not valid", http.StatusForbidden)
		return
	}

}
