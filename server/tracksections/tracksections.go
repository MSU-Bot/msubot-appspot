package tracksections

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"cloud.google.com/go/firestore"
	"github.com/SpencerCornish/msubot-appspot/server/mauth"
	"github.com/SpencerCornish/msubot-appspot/server/models"
	"github.com/SpencerCornish/msubot-appspot/server/serverutils"
	log "github.com/sirupsen/logrus"
)

const termRegex = `([0-9]){4}(30|50|70)`

type trackRequest struct {
	DepartmentAbbr string
	Course         string
	Crns           []string
	Term           string
}

// HandleRequest scrapes
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.WithContext(ctx).Infof("Context loaded. Starting execution.")

	// Only accept POST requests
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}

	authToken, err := mauth.VerifyToken(ctx, "TOKEN")
	if err != nil {
		http.Error(w, "Token is not valid", http.StatusForbidden)
		return
	}

	var request trackRequest
	err = decodeRequest(r.Body, &request)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// Get a firestore client
	fs := serverutils.GetFirebaseClient(ctx)
	writeBatch := fs.Batch()

	for _, crn := range request.Crns {
		// check for an existing tracked section
		docsForCrn := fs.Collection("sections_tracked").Where("crn", "==", crn).Where("term", "==", request.Term).Documents(ctx)
		docs, err := docsForCrn.GetAll()
		if err != nil {
			//TODO: Handle this shiz
		}

		if len(docs) != 0 {
			// if there are multiple docs, that's a big yikes. But we will still check for it and throw an error if it happens
			if len(docs) > 1 {
				log.WithContext(ctx).Error("Duplicate entries for the same class. Using the first one")
			}
			var trackedRecord models.TrackedSectionRecord
			err = docs[0].DataTo(&trackedRecord)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("Unable to parse existing section record")
				// TODO: Decide what we should do here...
			}

			for _, uid := range trackedRecord.Users {
				if uid == authToken.UID {
					log.WithContext(ctx).Warn("User is already tracking this crn, skipping")
					continue
				}
			}

			newUserSlice := append(trackedRecord.Users, authToken.UID)

			writeBatch.Set(docs[0].Ref, map[string]interface{}{
				"users": newUserSlice,
			}, firestore.MergeAll)

		} else {

			sectionData, err := getSectionMetadata(ctx, request, crn)
			if err != nil {
				http.Error(w, "Error getting course metadata", http.StatusInternalServerError)
				return
			}
			// Doc does not yet exist

		}

	}

}

func getSectionMetadata(ctx context.Context, request trackRequest, crn string) (models.Section, error) {
	client := http.DefaultClient

	resp, err := serverutils.MakeAtlasSectionRequest(client, request.Term, request.DepartmentAbbr, request.Course)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error getting new section data from ATLAS")
		return models.Section{}, err
	}
	sectionDatas, err := serverutils.ParseSectionResponse(resp, crn)
	if err != nil {
		return models.Section{}, err
	}

	return sectionDatas[0], nil

}

// decodeRequest decodes and validates the request
func decodeRequest(body io.ReadCloser, request *trackRequest) error {

	requestDecoder := json.NewDecoder(body)
	requestDecoder.DisallowUnknownFields()

	err := requestDecoder.Decode(&request)
	if err != nil {
		return err
	}

	isValidTerm, err := regexp.MatchString(termRegex, request.Term)
	if !isValidTerm || err != nil || len(request.Crns) == 0 {
		return err
	}

	return nil
}
