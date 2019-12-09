package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

// Section is a section model
type Section struct {
	DeptAbbr string
	DeptName string

	CourseName   string
	CourseNumber string
	CourseType   string
	Credits      string

	Instructor string
	Time       string
	Location   string

	SectionNumber  string
	Crn            string
	TotalSeats     string
	TakenSeats     string
	AvailableSeats string
}

// MakeAtlasSectionRequest makes a request to Atlas for section data in the term, department, and course
func MakeAtlasSectionRequest(client *http.Client, term, dept, course string) (*http.Response, error) {
	body := fmt.Sprintf("sel_subj=dummy&bl_online=FALSE&sel_day=dummy&term=%v&sel_subj=%v&sel_inst=ANY&sel_online=&sel_crse=%v&begin_hh=0&begin_mi=0&end_hh=0&end_mi=0",
		term,
		dept,
		course)

	req, err := http.NewRequest("POST", "https://prodmyinfo.montana.edu/pls/bzagent/bzskcrse.PW_ListSchClassSimple", strings.NewReader(body))
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ParseSectionResponse turns the http.Response into a slice of sections
func ParseSectionResponse(response *http.Response, crnToFind string) ([]Section, error) {
	sections := []Section{}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}
	rows := doc.Find("TR")
	for i := range rows.Nodes {
		columnsFr := rows.Eq(i).Find("TD")
		columnsSr := rows.Eq(i + 1).Find("TD")

		if columnsFr.Length()+columnsSr.Length() == 15 {
			matcher := regexp.MustCompile("[A-Za-z0-9]+")

			matches := matcher.FindAllString(columnsFr.Eq(1).Text(), -1)
			if len(matches) != 3 {
				panic("regex didn't work. Did the data model change?")
			}

			newSection := Section{
				DeptAbbr:       matches[0],
				CourseNumber:   matches[1],
				SectionNumber:  matches[2],
				CourseName:     strings.TrimSpace(columnsFr.Eq(2).Text()),
				Crn:            strings.TrimSpace(columnsFr.Eq(3).Text()),
				TotalSeats:     strings.TrimSpace(columnsFr.Eq(4).Text()),
				TakenSeats:     strings.TrimSpace(columnsFr.Eq(5).Text()),
				AvailableSeats: strings.TrimSpace(columnsFr.Eq(6).Text()),
				Instructor:     strings.TrimSpace(columnsFr.Eq(7).Text()),
				DeptName:       strings.TrimSpace(columnsSr.Eq(0).Text()),
				CourseType:     strings.TrimSpace(columnsSr.Eq(1).Text()),
				Time:           strings.TrimSpace(columnsSr.Eq(2).Text()),
				Location:       strings.TrimSpace(columnsSr.Eq(3).Text()),
				Credits:        strings.TrimSpace(columnsSr.Eq(4).Text()),
			}
			// Fixes recitation credits being blank, rather than 0 like they should be
			if newSection.Credits == "" {
				newSection.Credits = "0"
			}

			// We're looking for a specific section in this context,
			// so check if this is it, return it or continue if it's not
			if crnToFind != "" {
				if newSection.Crn == crnToFind {
					sections = append(sections, newSection)
					return sections, nil
				}
				continue
			}
			sections = append(sections, newSection)
		}
	}
	doc = nil
	return sections, nil
}

////////////////////////////
// Phone Functions
////////////////////////////

// PlivoRequest is the type sent to Plivo for texts
type PlivoRequest struct {
	Src  string `json:"src"`
	Dst  string `json:"dst"`
	Text string `json:"text"`
}

// SendText sends a text message to the specified phone number
func SendText(client *http.Client, number, message string) (response *http.Response, err error) {
	authID := os.Getenv("PLIVO_AUTH_ID")
	authToken := os.Getenv("PLIVO_AUTH_TOKEN")
	if authID == "" || authToken == "" {
		panic("nil env")
	}
	// TODO: Create sms callback handler
	url := fmt.Sprintf("https://api.plivo.com/v1/Account/%v/Message/", authID)
	data := PlivoRequest{Src: "14068000110", Dst: number, Text: message}

	js, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(authID, authToken)
	request.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return resp, err
}

// FetchUserDataWithNumber check firebase to see if the user exists in our database
func FetchUserDataWithNumber(ctx context.Context, fbClient *firestore.Client, number string) (map[string]interface{}, string) {
	checkedNumber := fmt.Sprintf("+%v", strings.Trim(number, " "))

	docs := fbClient.Collection("users").Where("number", "==", checkedNumber).Documents(ctx)

	parsed, err := docs.GetAll()
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("DoesUserExist Error")
		panic(err)
	}
	if len(parsed) > 0 {
		userData := parsed[0].Data()
		uid := parsed[0].Ref.ID
		return userData, uid
	}
	return nil, ""
}

// LookupUserNumber looks up a user's phone number from their uid
func LookupUserNumber(ctx context.Context, fbClient *firestore.Client, uid string) (string, error) {
	doc, err := fbClient.Collection("users").Doc(uid).Get(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Tracked user not found. This should've been cleaned up")
		return "", err
	}

	numIntf := doc.Data()["number"]

	num, ok := numIntf.(string)
	if !ok {
		log.WithContext(ctx).Errorf("Unable to parse user's number as a string, they may be in the midst of resetting it")
		return "", err
	}

	return num, nil
}

// GetFirebaseClient creates and returns a new firebase client, used to interact with the database
func GetFirebaseClient(ctx context.Context) *firestore.Client {
	firebasePID := os.Getenv("FIREBASE_PROJECT")
	log.WithContext(ctx).Infof("Loaded firebase project ID.")
	if firebasePID == "" {
		log.WithContext(ctx).Errorf("Firebase Project ID is nil, I cannot continue.")
		panic("Firebase Project ID is nil")
	}

	fbClient, err := firestore.NewClient(ctx, firebasePID)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Could not create new client for Firebase")
		return nil
	}
	log.WithContext(ctx).Infof("successfully opened firestore client")

	return fbClient
}

// MoveTrackedSection moves old sections out of the prod area
func MoveTrackedSection(ctx context.Context, fbClient *firestore.Client, crn, uid, term string) error {

	// Look for an existing archive doc to add userdata to
	docArchiveIter := fbClient.Collection("sections_archive").Where("term", "==", term).Where("crn", "==", crn).Documents(ctx)
	archiveDocs, err := docArchiveIter.GetAll()

	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Could not get list of archive docs for uid %v: %v", uid, err)
		return err
	}

	// Get the document that we need to move
	docToMove, err := fbClient.Collection("sections_tracked").Doc(uid).Get(ctx)
	docToMoveData := docToMove.Data()

	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Could not get the new doc for uid %s : %v", uid, err)
		return err
	}

	//  if there is a doc, merge with it rather than making a new one
	if archiveDocs != nil {
		if len(archiveDocs) > 1 {
			log.WithContext(ctx).WithError(err).Errorf("Duplicate archiveDocs: %v", archiveDocs)
		}

		//  Get the data for the archive docs
		data := archiveDocs[0].Data()

		// get all the users
		users, ok := data["users"].([]interface{})
		if !ok {
			log.WithContext(ctx).WithError(err).Errorf("couldn't parse userslice")
			return nil
		}

		// get all the users
		usersToAdd, ok := docToMoveData["users"].([]interface{})
		if !ok {
			log.WithContext(ctx).WithError(err).Errorf("couldn't parse userslice")
			return nil
		}

		//  make a mega list
		allUsers := append(users, usersToAdd...)

		// Update that userlist
		_, err := archiveDocs[0].Ref.Set(ctx, map[string]interface{}{
			"users": allUsers,
		}, firestore.MergeAll)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("Error appending users to archive")
			return err
		}

	} else {

		// Add a new doc
		_, _, err := fbClient.Collection("sections_archive").Add(ctx, docToMoveData)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("Error creating a new archived doc")
			return err
		}

	}

	//  Finally delete the old one
	_, err = docToMove.Ref.Delete(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Error deleting old document")
		return err
	}

	return nil

}
