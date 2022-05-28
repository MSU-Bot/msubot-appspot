package usercrud

import (
	"context"

	"github.com/SpencerCornish/msubot-appspot/server/dstore"
)

func RemoveTrackedSection(ctx context.Context, userID, sectionID string, ds dstore.DStore) error {

	section, err := ds.GetTrackedSectionByID(ctx, sectionID)
	if err != nil {
		return err
	}

	for idx, uid := range section.Users {
		if uid == userID {
			section.Users = append(section.Users[:idx], section.Users[idx+1:]...)
			break
		}
	}

	return ds.UpdateTrackedSection(ctx, *section)
}

func AddTrackedSections(ctx context.Context, userID, term string, sectionIDs []string) error {

	// writeBatch := fs.Batch()
	// modifiedDocIDs := make([]string, len(request.Crns))
	// for i, crn := range request.Crns {
	// 	// check for an existing tracked section
	// 	docsForCrn := fs.Collection("sections_tracked").Where("crn", "==", crn).Where("term", "==", request.Term).Documents(ctx)
	// 	docs, err := docsForCrn.GetAll()
	// 	if err != nil {
	// 		log.WithContext(ctx).WithError(err).Error("Failed to fetch existing tracked sections")
	// 		http.Error(w, "Failed to fetch existing records", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	if len(docs) != 0 {
	// 		// if there are multiple docs, that's a big yikes. But we will still check for it and throw an error if it happens
	// 		if len(docs) > 1 {
	// 			log.WithContext(ctx).Error("Duplicate entries for the same class. Using the first one")
	// 		}

	// 		var trackedRecord models.TrackedSectionRecord
	// 		err = docs[0].DataTo(&trackedRecord)
	// 		if err != nil {
	// 			log.WithContext(ctx).WithError(err).Error("Unable to parse existing section record")
	// 			http.Error(w, "Internal data error", http.StatusInternalServerError)
	// 			return
	// 		}

	// 		for _, uid := range trackedRecord.Users {
	// 			if uid == authToken.UID {
	// 				log.WithContext(ctx).Warn("User is already tracking this crn, skipping")
	// 				continue
	// 			}
	// 		}

	// 		newUserSlice := append(trackedRecord.Users, authToken.UID)

	// 		writeBatch.Set(docs[0].Ref, map[string]interface{}{
	// 			"users": newUserSlice,
	// 		}, firestore.MergeAll)

	// 		modifiedDocIDs[i] = docs[0].Ref.ID

	// 	} else {
	// 		// This section has not been tracked yet. Go get it
	// 		sectionData, err := getSectionMetadata(ctx, request, crn)
	// 		if err != nil {
	// 			log.WithContext(ctx).WithError(err).WithFields(log.Fields{"request": request, "crn": crn}).Error("Failed to request data from ATLAS")
	// 			http.Error(w, "Error getting course metadata", http.StatusInternalServerError)
	// 			return
	// 		}
	// 		// if we didn't find a section, it's probably an invalid request. Log as an error for now, just to be safe
	// 		if sectionData == (models.Section{}) {
	// 			log.WithContext(ctx).WithFields(log.Fields{"request": request, "crn": crn}).Error("Didn't find a section. Invalid request")
	// 			http.Error(w, "Invalid Request", http.StatusBadRequest)
	// 			return
	// 		}
	// 		newDocRef := fs.Collection("sections_tracked").NewDoc()
	// 		modifiedDocIDs[i] = newDocRef.ID

	// 		writeBatch.Set(newDocRef, map[string]interface{}{
	// 			"term":           sectionData.Term,
	// 			"departmentAbbr": sectionData.DeptAbbr,
	// 			"department":     sectionData.DeptName,
	// 			"courseName":     sectionData.CourseName,
	// 			"crn":            sectionData.Crn,
	// 			"courseNumber":   sectionData.CourseNumber,
	// 			"openSeats":      sectionData.AvailableSeats,
	// 			"totalSeats":     sectionData.TotalSeats,
	// 			"sectionNumber":  sectionData.SectionNumber,
	// 			"instructor":     sectionData.Instructor,
	// 			"users":          []string{authToken.UID},
	// 			"creationTime":   firestore.ServerTimestamp,
	// 		})

	// 	}

	// }

	// _, err = writeBatch.Commit(ctx)
	// if err != nil {
	// 	log.WithContext(ctx).WithError(err).Error("Failed to commit writeBatch")
	// 	http.Error(w, "Failure tracking classes", http.StatusInternalServerError)
	// }
	// responseBody := trackResponse{DocIDs: modifiedDocIDs}
	// resp, err := json.Marshal(responseBody)
	// if err != nil {
	// 	log.WithContext(ctx).WithError(err).WithField("responseBody", responseBody).Error("Failed to marshal response body")
	// 	http.Error(w, "Failed to parse response", http.StatusInternalServerError)
	// 	return
	// }

	// w.Write(resp)
	return nil
}

func UpdateUserData() error {
	return nil
}
