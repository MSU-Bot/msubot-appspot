package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"cloud.google.com/go/firestore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func CheckSectionsHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	log.Infof(ctx, "Context loaded. Starting execution.")
	// TODO: Check for cron header

	firebasePID := os.Getenv("FIREBASE_PROJECT")
	log.Debugf(ctx, "Loaded firebase project ID: %v", firebasePID)
	if firebasePID == "" {
		log.Criticalf(ctx, "Firebase Project ID is nil, I cannot continue.")
		panic("Firebase Project ID is nil")
	}

	fbClient, err := firestore.NewClient(ctx, firebasePID)
	defer fbClient.Close()
	if err != nil {
		log.Errorf(ctx, "Could not create new client for Firebase %v", err)
		w.WriteHeader(500)
		return
	}
	log.Debugf(ctx, "successfully opened firestore client")

	sectionsSnapshot := fbClient.Collection("tracked_sections").Documents(ctx)

	sectionDocuments, err := sectionsSnapshot.GetAll()
	if err != nil {
		log.Errorf(ctx, "Error getting tracked_sections!")
		w.WriteHeader(500)
		return
	}
	log.Debugf(ctx, "successfully got a slice of sectionDocuments we are tracking")

	// For fun tracking data
	var notifiedSections, notifiedUsers int

	fbBatch := fbClient.Batch()
	var waitGroup sync.WaitGroup

	for _, doc := range sectionDocuments {
		waitGroup.Add(1)
		go sectionCheckWorker(ctx, &waitGroup, doc, client, fbClient, fbBatch)
	}
	waitGroup.Wait()
	_, err = fbBatch.Commit(ctx)
	if err != nil {
		log.Criticalf(ctx, "Writebatch failed: %v", err)
	}
	log.Infof(ctx, "Tracked Courses:%v Total Notified: Users:%v Sections:%v", len(sectionDocuments), notifiedUsers, notifiedSections)
	w.WriteHeader(200)

	ctx.Done()
	return
}

func sectionCheckWorker(ctx context.Context, wg *sync.WaitGroup, doc *firestore.DocumentSnapshot, client *http.Client, fbClient *firestore.Client, fbBatch *firestore.WriteBatch) {
	defer wg.Done()

	// Get section data
	sectionData := doc.Data()
	if sectionData == nil {
		log.Errorf(ctx, "Big, unexpected error with getting data for course")
		return
	}

	// Type conversions
	term, ok := sectionData["term"].(string)
	if !ok {
		log.Errorf(ctx, "type conv failed for courseNumber")
		panic("foo")
	}
	departmentAbbr, ok := sectionData["departmentAbbr"].(string)
	if !ok {
		log.Errorf(ctx, "type conv failed for courseNumber")
		panic("foo")
	}
	courseNumber, ok := sectionData["courseNumber"].(string)
	if !ok {
		log.Errorf(ctx, "type conv failed for courseNumber")
		panic("foo")
	}
	crn := doc.Ref.ID

	// Make a request to Atlas
	resp, err := MakeAtlasSectionRequest(client, term, departmentAbbr, courseNumber)
	if err != nil {
		log.Errorf(ctx, "Making Atlas request failed for %v: %v", departmentAbbr+courseNumber, err)
		return
	}

	// Parse into a section struct
	newSectionData, err := ParseSectionResponse(resp, crn)
	if err != nil {
		log.Errorf(ctx, "Parsing section failed: %v", err)

	}
	// If we somehow get back more than one section, something super borked earlier
	if len(newSectionData) > 1 {
		log.Errorf(ctx, "Something went wrong with parsing the section response. Expected 1 section, recieved %v", len(newSectionData))
		return
	}

	// If we didn't get back any, notify the users that stuff's broke, and move on
	if len(newSectionData) == 0 {
		log.Warningf(ctx, "Couldn't find section. Proceeding with notify and delete. This should be pretty rare")
		users, ok := sectionData["users"].([]interface{})
		if !ok {
			log.Errorf(ctx, "couldn't parse userslice")
		}
		sendDeletedSectionMessages(ctx, client, fbClient, users, newSectionData[0])
		fbBatch.Delete(fbClient.Collection("tracked_sections").Doc(crn))
		return
	}
	// Parse the new available seats to an int
	newSeatsAvailable, err := strconv.Atoi(newSectionData[0].AvailableSeats)
	if err != nil {
		log.Errorf(ctx, "couldn't parse newSeatsAvailable: %v", err)
		return
	}

	// If there are seats available
	if newSeatsAvailable > 0 {
		users, ok := sectionData["users"].([]interface{})
		if !ok {
			log.Errorf(ctx, "couldn't parse userslice")
			return
		}
		sendOpenSeatMessages(ctx, client, fbClient, users, newSectionData[0])
		removeSectionFromUserData(ctx, fbClient, fbBatch, users, newSectionData[0].Crn)
		fbBatch.Delete(fbClient.Collection("tracked_sections").Doc(crn))
		return
	}

	// If we get here, we just need to update the stored section model so it's all clean and nice
	updateTrackedSection(ctx, fbBatch, fbClient.Collection("tracked_sections").Doc(newSectionData[0].Crn), newSectionData[0])
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
			log.Warningf(ctx, "found legacy user with no tracked sections.")
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

		}
		if userNumbers == "" {
			userNumbers = number
		} else {
			userNumbers = fmt.Sprintf("%v<%v", userNumbers, number)
		}
	}
	log.Infof(ctx, userNumbers)
	resp, err := SendText(client, userNumbers, message)
	if err != nil {
		log.Errorf(ctx, "error sending text: %v", err)
		return err
	}
	log.Debugf(ctx, "%v", resp)
	return nil
}

func sendDeletedSectionMessages(ctx context.Context, client *http.Client, fbClient *firestore.Client, users []interface{}, section Section) error {
	var userNumbers string
	message := fmt.Sprintf("%v%v - %v with CRN %v appears to have been closed by MSU. Please check MyInfo to confirm: %v ", section.DeptAbbr, section.CourseNumber, section.CourseName, section.Crn, "https://goo.gl/58VYz5")
	for _, user := range users {
		number, err := LookupUserNumber(ctx, fbClient, user.(string))
		if err != nil {

		}
		userNumbers = fmt.Sprintf("%v<%v", userNumbers, number)
	}
	resp, err := SendText(client, userNumbers, message)
	if err != nil {
		log.Errorf(ctx, "error sending text: %v", err)
		return err
	}
	log.Debugf(ctx, "%v", resp)

	return nil
}
