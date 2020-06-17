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

type trackResponse struct {
	DocIDs []string
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

	var request trackRequest
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
	modifiedDocIDs := make([]string, len(request.Crns))
	for i, crn := range request.Crns {
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

			modifiedDocIDs[i] = docs[0].Ref.ID

		} else {
			// This section has not been tracked yet. Go get it
			sectionData, err := getSectionMetadata(ctx, request, crn)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithFields(log.Fields{"request": request, "crn": crn}).Error("Failed to request data from ATLAS")
				http.Error(w, "Error getting course metadata", http.StatusInternalServerError)
				return
			}
			// if we didn't find a section, it's probably an invalid request. Log as an error for now, just to be safe
			if sectionData == (models.Section{}) {
				log.WithContext(ctx).WithFields(log.Fields{"request": request, "crn": crn}).Error("Didn't find a section. Invalid request?")
				http.Error(w, "Invalid Request", http.StatusBadRequest)
				return
			}
			newDocRef := fs.Collection("sections_tracked").NewDoc()
			modifiedDocIDs[i] = newDocRef.ID

			writeBatch.Set(newDocRef, map[string]interface{}{
				"term":           sectionData.Term,
				"departmentAbbr": sectionData.DeptAbbr,
				"department":     sectionData.DeptName,
				"courseName":     sectionData.CourseName,
				"crn":            sectionData.Crn,
				"courseNumber":   sectionData.CourseNumber,
				"openSeats":      sectionData.AvailableSeats,
				"totalSeats":     sectionData.TotalSeats,
				"sectionNumber":  sectionData.SectionNumber,
				"instructor":     sectionData.Instructor,
				"users":          []string{authToken.UID},
				"creationTime":   firestore.ServerTimestamp,
			})

		}

	}

	_, err = writeBatch.Commit(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to commit writeBatch")
		http.Error(w, "Failure tracking classes", http.StatusInternalServerError)
	}
	responseBody := trackResponse{DocIDs: modifiedDocIDs}
	resp, err := json.Marshal(responseBody)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("responseBody", responseBody).Error("Failed to marshal response body")
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func getSectionMetadata(ctx context.Context, request trackRequest, crn string) (models.Section, error) {
	client := http.DefaultClient

	resp, err := serverutils.MakeAtlasSectionRequest(client, request.Term, request.DepartmentAbbr, request.Course)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error getting new section data from ATLAS")
		return models.Section{}, err
	}
	sectionDatas, err := serverutils.ParseSectionResponse(resp, request.Term, crn)
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
