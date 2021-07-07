package removetracked

import (
	"encoding/json"
	"io"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/SpencerCornish/msubot-appspot/server/mauth"
	"github.com/SpencerCornish/msubot-appspot/server/models"
	"github.com/SpencerCornish/msubot-appspot/server/serverutils"
	log "github.com/sirupsen/logrus"
)

const termRegex = `([0-9]){4}(30|50|70)`

type removeTrackedRequest struct {
	TrackingIDs []string
}

type removeTrackedResponse struct {
	Success bool
}

// HandleRequest scrapes
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.WithContext(ctx).Infof("Context loaded. Starting execution.")

	// Only accept POST requests
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}

	authToken, err := mauth.VerifyToken(ctx, r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Token is not valid", http.StatusForbidden)
		return
	}

	var request removeTrackedRequest
	err = decodeRequest(r.Body, &request)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Get a firestore client
	fs := serverutils.GetFirebaseClient(ctx)
	defer fs.Close()

	writeBatch := fs.Batch()
	modifiedDocIDs := make([]string, len(request.TrackingIDs))
	for i, ID := range request.TrackingIDs {
		// check for an existing tracked section
		docRef := fs.Collection("sections_tracked").Doc(ID)
		document, err := docRef.Get(ctx)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to fetch existing tracked sections")
			http.Error(w, "Failed to fetch existing records", http.StatusInternalServerError)
			return
		}
		if !document.Exists() {
			log.WithContext(ctx).WithField("DocID", ID).Warn("Course with id not found, skipping")
		}

		var trackedRecord models.TrackedSectionRecord
		err = document.DataTo(&trackedRecord)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Unable to parse existing section record")
			http.Error(w, "Internal data error", http.StatusInternalServerError)
			return
		}

		for idx, uid := range trackedRecord.Users {
			if uid == authToken.UID {
				trackedRecord.Users = append(trackedRecord.Users[:idx], trackedRecord.Users[i+1:]...)
				break
			}
		}

		writeBatch.Set(docRef, map[string]interface{}{
			"users": trackedRecord.Users,
		}, firestore.MergeAll)

		modifiedDocIDs[i] = docRef.ID

	}

	_, err = writeBatch.Commit(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to commit writeBatch")
		http.Error(w, "Failure to commit", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
}

// decodeRequest decodes and validates the request
func decodeRequest(body io.ReadCloser, request *removeTrackedRequest) error {

	requestDecoder := json.NewDecoder(body)
	requestDecoder.DisallowUnknownFields()

	err := requestDecoder.Decode(&request)
	if err != nil {
		return err
	}

	return nil
}
