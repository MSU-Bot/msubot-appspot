package checksections

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SpencerCornish/msubot-appspot/server/models"
	"github.com/SpencerCornish/msubot-appspot/server/serverutils"

	"cloud.google.com/go/firestore"

	log "github.com/sirupsen/logrus"
)

// HandleRequest runs often to check for open seats
func HandleRequest(w http.ResponseWriter, r *http.Request) {

	// Load up a context and http client
	ctx := r.Context()
	client := http.DefaultClient

	log.WithContext(ctx).Infof("Context loaded. Starting execution.")

	// Make sure the request is from the appengine cron
	// if r.Header.Get("X-Appengine-Cron") == "" {
	// 	log.Warningf(ctx, "Request is not from the cron. Exiting")
	// 	w.WriteHeader(403)
	// 	return
	// }

	fbClient := serverutils.GetFirebaseClient(ctx)
	if fbClient == nil {
		w.WriteHeader(500)
		return
	}
	defer fbClient.Close()

	// Get the list of sections we are actively tracking
	sectionsSnapshot := fbClient.Collection("sections_tracked").Documents(ctx)

	// Actually get all the data within these docs
	sectionDocuments, err := sectionsSnapshot.GetAll()
	if err != nil {
		log.WithContext(ctx).Errorf("Error getting sections_tracked! sec: %v", sectionDocuments)
		log.WithContext(ctx).Errorf("Error getting sections_tracked! Err: %v", err)
		w.WriteHeader(500)
		return
	}
	log.WithContext(ctx).Infof("successfully got sectionDocuments we are tracking")

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
		log.WithContext(ctx).Errorf("Writebatch failed: %v", err)
	}
	w.WriteHeader(200)

	ctx.Done()
	return
}

func sectionCheckWorker(ctx context.Context, jobs <-chan *firestore.DocumentSnapshot, returnChannel chan<- int, client *http.Client, fbClient *firestore.Client, fbBatch *firestore.WriteBatch) {
	for doc := range jobs {
		// The unique doc ID
		sectionUID := doc.Ref.ID
		// Get section data
		sectionData := doc.Data()
		if sectionData == nil {
			log.WithContext(ctx).Errorf("unexpected error with getting data for course")

			returnChannel <- 0
			continue
		}

		// Type conversions
		term, ok := sectionData["term"].(string)
		if !ok {
			log.WithContext(ctx).Errorf("type conv failed for courseNumber")

			returnChannel <- 0
			continue
		}
		departmentAbbr, ok := sectionData["departmentAbbr"].(string)
		if !ok {
			log.WithContext(ctx).Errorf("type conv failed for courseNumber")

			returnChannel <- 0
			continue
		}
		courseNumber, ok := sectionData["courseNumber"].(string)
		if !ok {
			log.WithContext(ctx).Errorf("type conv failed for courseNumber")

			returnChannel <- 0
			continue
		}
		crn := sectionData["crn"].(string)

		// Make a request to Atlas
		resp, err := serverutils.MakeAtlasSectionRequest(client, term, departmentAbbr, courseNumber)
		if err != nil {
			log.WithContext(ctx).Errorf("Making Atlas request failed for %v: %v", departmentAbbr+courseNumber, err)

			returnChannel <- 0
			continue
		}

		// Parse into a section struct
		newSectionData, err := serverutils.ParseSectionResponse(resp, crn)
		if err != nil {
			log.WithContext(ctx).Errorf("Parsing section failed: %v", err)

			returnChannel <- 0
			continue
		}
		// If we somehow get back more than one section, something super borked earlier
		if len(newSectionData) > 1 {
			log.WithContext(ctx).Errorf("Something went wrong with parsing the section response. Expected 1 section, recieved %v", len(newSectionData))

			returnChannel <- 0
			continue
		}

		// If we didn't get back any, warn us and move on.
		// This typically occurs when Banner is down
		if len(newSectionData) == 0 {
			log.WithContext(ctx).Infof("Couldn't find section from MSU: %v", crn)

			returnChannel <- 0
			continue
		}

		// Parse the new available seats to an int
		newSeatsAvailable, err := strconv.Atoi(newSectionData[0].AvailableSeats)
		if err != nil {
			log.WithContext(ctx).Errorf("couldn't parse newSeatsAvailable: %v", err)

			returnChannel <- 0
			continue
		}

		users, ok := sectionData["users"].([]interface{})
		if !ok {
			log.WithContext(ctx).Errorf("couldn't parse userslice")
			returnChannel <- 0
			continue
		}

		if len(users) < 1 {
			log.WithContext(ctx).Infof("CRN %s has %d users. Deleting CRN", crn, len(users))
			err := serverutils.MoveTrackedSection(ctx, fbClient, newSectionData[0].Crn, sectionUID, term)
			if err != nil {
				log.WithContext(ctx).Errorf("Failed to move the stale section data: %v", err)
			}
			returnChannel <- 0
			continue
		}
		log.WithContext(ctx).Infof("seats available for %v:%v", crn, newSeatsAvailable)
		// If there are seats available
		if newSeatsAvailable > 0 {
			users, ok := sectionData["users"].([]interface{})
			if !ok {
				log.WithContext(ctx).Errorf("couldn't parse userslice")

				returnChannel <- 0
				continue
			}
			log.WithContext(ctx).Infof("The CRN %s has %d open seats. Sending a message to %d users.", crn, newSeatsAvailable, len(users))
			sendOpenSeatMessages(ctx, client, fbClient, users, newSectionData[0])
			removeSectionFromUserData(ctx, fbClient, fbBatch, users, newSectionData[0].Crn)

			err := serverutils.MoveTrackedSection(ctx, fbClient, newSectionData[0].Crn, sectionUID, term)
			if err != nil {
				log.WithContext(ctx).Errorf("Failed to move the stale section data: %v", err)
			}

			returnChannel <- 0
			continue
		}

		// If we get here, we just need to update the stored section model so it's all clean and nice
		updateTrackedSection(ctx, fbBatch, fbClient.Collection("sections_tracked").Doc(sectionUID), newSectionData[0])

		returnChannel <- 0
	}
}

func updateTrackedSection(ctx context.Context, fbBatch *firestore.WriteBatch, fbRef *firestore.DocumentRef, section models.Section) {
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
			log.WithContext(ctx).Errorf("couldn't find user when trying to remove crn ref")
			continue
		}
		userDataMap := userData.Data()
		if userDataMap == nil {
			continue
		}
		untypedUserdata := userDataMap["sections"]
		if untypedUserdata == nil {
			log.WithContext(ctx).Infof("found user with no tracked sections on their userdata.")
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

func sendOpenSeatMessages(ctx context.Context, client *http.Client, fbClient *firestore.Client, users []interface{}, section models.Section) error {
	var userNumbers string
	message := fmt.Sprintf("%v%v - %v with CRN %v has %v open seats! Get to MyInfo and register before it's gone!", section.DeptAbbr, section.CourseNumber, section.CourseName, section.Crn, section.AvailableSeats)
	for _, user := range users {
		number, err := serverutils.LookupUserNumber(ctx, fbClient, user.(string))
		if err != nil {
			log.WithContext(ctx).Errorf("Unable to send a text to user %s", user.(string))
		}
		if userNumbers == "" {
			userNumbers = number
		} else {
			userNumbers = fmt.Sprintf("%v<%v", userNumbers, number)
		}
	}
	resp, err := serverutils.SendText(client, userNumbers, message)
	if err != nil {
		log.WithContext(ctx).Errorf("error sending text: %v", err)
		return err
	}
	log.WithContext(ctx).Infof("%v", resp)
	return nil
}
