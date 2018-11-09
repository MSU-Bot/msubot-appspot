package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

// CheckSectionsHandler runs often to check for open seats
func CheckSectionsHandler(w http.ResponseWriter, r *http.Request) {

	// Load up a context and http client
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	log.Infof(ctx, "Context loaded. Starting execution.")

	// Make sure the request is from the appengine cron
	if r.Header.Get("X-Appengine-Cron") == "" {
		log.Warningf(ctx, "Request is not from the cron. Exiting")
		w.WriteHeader(403)
		return
	}

	fbClient := GetFirebaseClient(ctx)
	if fbClient == nil {
		w.WriteHeader(500)
		return
	}
	defer fbClient.Close()

	// Get the list of sections we are actively tracking
	sectionsSnapshot := fbClient.Collection("tracked_sections").Documents(ctx)

	// Actually get all the data within these docs
	sectionDocuments, err := sectionsSnapshot.GetAll()
	if err != nil {
		log.Errorf(ctx, "Error getting tracked_sections! sec: %v", sectionDocuments)
		log.Errorf(ctx, "Error getting tracked_sections! Err: %v", err)
		w.WriteHeader(500)
		return
	}
	log.Debugf(ctx, "successfully got sectionDocuments we are tracking")

	fbBatch := fbClient.Batch()

	// This is the number of concurrent URLFetches that we will do.
	numWorkers := 10

	// A queue of all sections to check
	jobQueue := make(chan *firestore.DocumentSnapshot, len(sectionDocuments))

	// A return channel to let us know a job has completed
	requestCompleteChannel := make(chan int, len(sectionDocuments))

	// Start up some workers
	for r := 0; r < numWorkers; r++ {
		go sectionCheckWorker(ctx, jobQueue, requestCompleteChannel, client, fbClient, fbBatch)
	}

	// Add all sections to the queue
	for _, doc := range sectionDocuments {
		jobQueue <- doc
	}

	close(jobQueue)

	// Wait for the jobs to finish
	for i := 0; i < len(sectionDocuments); i++ {
		<-requestCompleteChannel
	}

	_, err = fbBatch.Commit(ctx)
	if err != nil {
		log.Criticalf(ctx, "Writebatch failed: %v", err)
	}
	w.WriteHeader(200)

	ctx.Done()
	return
}

func sectionCheckWorker(ctx context.Context, jobs <-chan *firestore.DocumentSnapshot, returnChannel chan<- int, client *http.Client, fbClient *firestore.Client, fbBatch *firestore.WriteBatch) {

	for doc := range jobs {

		// Get section data
		sectionData := doc.Data()
		if sectionData == nil {
			log.Errorf(ctx, "unexpected error with getting data for course")

			returnChannel <- 0
			continue
		}

		// Type conversions
		term, ok := sectionData["term"].(string)
		if !ok {
			log.Errorf(ctx, "type conv failed for courseNumber")

			returnChannel <- 0
			continue
		}
		departmentAbbr, ok := sectionData["departmentAbbr"].(string)
		if !ok {
			log.Errorf(ctx, "type conv failed for courseNumber")

			returnChannel <- 0
			continue
		}
		courseNumber, ok := sectionData["courseNumber"].(string)
		if !ok {
			log.Errorf(ctx, "type conv failed for courseNumber")

			returnChannel <- 0
			continue
		}
		crn := doc.Ref.ID

		// Make a request to Atlas
		resp, err := MakeAtlasSectionRequest(client, term, departmentAbbr, courseNumber)
		if err != nil {
			log.Errorf(ctx, "Making Atlas request failed for %v: %v", departmentAbbr+courseNumber, err)

			returnChannel <- 0
			continue
		}

		// Parse into a section struct
		newSectionData, err := ParseSectionResponse(resp, crn)
		if err != nil {
			log.Errorf(ctx, "Parsing section failed: %v", err)

			returnChannel <- 0
			continue
		}
		// If we somehow get back more than one section, something super borked earlier
		if len(newSectionData) > 1 {
			log.Errorf(ctx, "Something went wrong with parsing the section response. Expected 1 section, recieved %v", len(newSectionData))

			returnChannel <- 0
			continue
		}

		// If we didn't get back any, warn us and move on.
		// This typically occurs when Banner is down
		if len(newSectionData) == 0 {
			log.Warningf(ctx, "Couldn't find section from MSU")

			returnChannel <- 0
			continue
		}

		// Parse the new available seats to an int
		newSeatsAvailable, err := strconv.Atoi(newSectionData[0].AvailableSeats)
		if err != nil {
			log.Errorf(ctx, "couldn't parse newSeatsAvailable: %v", err)

			returnChannel <- 0
			continue
		}

		users, ok := sectionData["users"].([]interface{})
		if !ok {
			log.Errorf(ctx, "couldn't parse userslice")
			returnChannel <- 0
			continue
		}

		if len(users) < 1 {
			log.Infof(ctx, "CRN %s has %d users. Deleting CRN", crn, len(users))
			fbBatch.Delete(fbClient.Collection("tracked_sections").Doc(crn))
			returnChannel <- 0
			continue
		}

		// If there are seats available
		if newSeatsAvailable > 0 {
			users, ok := sectionData["users"].([]interface{})
			if !ok {
				log.Errorf(ctx, "couldn't parse userslice")

				returnChannel <- 0
				continue
			}
			log.Infof(ctx, "The CRN %s has %d open seats. Sending a message to %d users.", crn, newSeatsAvailable, len(users))
			sendOpenSeatMessages(ctx, client, fbClient, users, newSectionData[0])
			removeSectionFromUserData(ctx, fbClient, fbBatch, users, newSectionData[0].Crn)
			fbBatch.Delete(fbClient.Collection("tracked_sections").Doc(crn))

			returnChannel <- 0
			continue
		}

		// If we get here, we just need to update the stored section model so it's all clean and nice
		updateTrackedSection(ctx, fbBatch, fbClient.Collection("tracked_sections").Doc(newSectionData[0].Crn), newSectionData[0])

		returnChannel <- 0
	}
}

func updateTrackedSection(ctx context.Context, fbBatch *firestore.WriteBatch, fbRef *firestore.DocumentRef, section Section) {
	fbBatch.Set(fbRef, map[string]interface{}{
		"department": section.DeptName,
		"courseName": section.CourseName,
		"openSeats":  section.AvailableSeats,
		"totalSeats": section.TotalSeats,
		"instructor": section.Instructor,
	}, firestore.MergeAll)

}

func removeSectionFromUserData(ctx context.Context, fbClient *firestore.Client, fbBatch *firestore.WriteBatch, users []interface{}, crn string) {
	for _, user := range users {

		userData, err := fbClient.Collection("users").Doc(user.(string)).Get(ctx)
		if err != nil {
			log.Errorf(ctx, "couldn't find user when trying to remove crn ref")
			continue
		}
		userDataMap := userData.Data()
		if userDataMap == nil {
			continue
		}
		untypedUserdata := userDataMap["sections"]
		if untypedUserdata == nil {
			log.Infof(ctx, "found user with no tracked sections on their userdata.")
			continue
		}
		sectionSlice := untypedUserdata.([]interface{})
		for i, curCrn := range sectionSlice {
			if curCrn == crn {
				sectionSlice = append(sectionSlice[:i], sectionSlice[i+1:]...)
				break
			}
		}
		fbBatch.Set(userData.Ref, map[string]interface{}{
			"sections": sectionSlice,
		}, firestore.MergeAll)

	}
}

// Messaging
func sendOpenSeatMessages(ctx context.Context, client *http.Client, fbClient *firestore.Client, users []interface{}, section Section) error {
	var userNumbers string
	message := fmt.Sprintf("%v%v - %v with CRN %v has %v open seats! Get to MyInfo and register before it's gone!", section.DeptAbbr, section.CourseNumber, section.CourseName, section.Crn, section.AvailableSeats)
	for _, user := range users {
		number, err := LookupUserNumber(ctx, fbClient, user.(string))
		if err != nil {
			log.Errorf(ctx, "Unable to send a text to user %s", user.(string))
		}
		if userNumbers == "" {
			userNumbers = number
		} else {
			userNumbers = fmt.Sprintf("%v<%v", userNumbers, number)
		}
	}
	resp, err := SendText(client, userNumbers, message)
	if err != nil {
		log.Errorf(ctx, "error sending text: %v", err)
		return err
	}
	log.Debugf(ctx, "%v", resp)
	return nil
}
