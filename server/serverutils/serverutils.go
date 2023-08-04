package serverutils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/PuerkitoBio/goquery"
	"github.com/SpencerCornish/msubot-appspot/server/models"
	log "github.com/sirupsen/logrus"
)

const (
	sectionRequestURL = "https://prodmyinfo.montana.edu/pls/bzagent/bzskcrse.PW_ListSchClassSimple"
)

// MakeAtlasSectionRequest makes a request to Atlas for section data in the term, department, and course
func MakeAtlasSectionRequest(client *http.Client, term, dept, course string) (*http.Response, error) {
	body := buildAtlasRequestBody(term, dept, course)

	req, err := http.NewRequest("POST", sectionRequestURL, body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func buildAtlasRequestBody(term, department, course string) io.Reader {
	body := fmt.Sprintf("sel_subj=dummy&bl_online=FALSE&sel_day=dummy&term=%v&sel_subj=%v&sel_inst=ANY&sel_online=&sel_crse=%v&begin_hh=0&begin_mi=0&end_hh=0&end_mi=0",
		term,
		department,
		course)

	return strings.NewReader(body)
}

// ParseSectionResponse turns the http.Response into a slice of sections
func ParseSectionResponse(response *http.Response, termString, crnToFind string) ([]models.Section, error) {
	sections := []models.Section{}
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

			newSection := models.Section{
				Term:           termString,
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
	return doc.Data()["number"].(string), nil
}

// GetFirebaseClient creates and returns a new firebase client, used to interact with the database
func GetFirebaseClient(ctx context.Context) *firestore.Client {
	firebasePID := os.Getenv("FIREBASE_PROJECT")
	log.WithContext(ctx).Infof("Loaded firebase project ID.")
	if firebasePID == "" {
		log.WithContext(ctx).Fatal("Firebase Project ID is nil, I cannot continue.")
	}

	fbClient, err := firestore.NewClient(ctx, firebasePID)
	if err != nil {
		log.WithContext(ctx).WithError(err).Fatal("Could not create new client for Firebase")
	}
	log.WithContext(ctx).Infof("successfully opened firestore client")

	if value := os.Getenv("FIRESTORE_EMULATOR_HOST"); value != "" {
		log.Warningf("Using Firestore Emulator: %s", value)
		// err := addTestingData(fbClient)
		// log.WithError(err).Info("added data")
	}

	return fbClient
}

// func addTestingData(fbClient *firestore.Client) error {
// 	wb := fbClient.Batch()

// 	// Globals
// 	wb.Create(fbClient.Collection("global").Doc("global"), map[string]interface{}{
// 		"coursesTracked": 999,
// 		"motd":           "",
// 		"textsSent":      -1,
// 		"users":          999,
// 	})

// 	// Depts
// 	wb.Create(fbClient.Collection("departments").Doc("CSCI"), map[string]interface{}{
// 		"name":        "Computerz",
// 		"updatedTime": time.Now(),
// 	})

// 	//Dept Classes
// 	wb.Create(fbClient.Collection("departments").Doc("CSCI").Collection("202270").Doc("999D"), map[string]interface{}{
// 		"title": "SPENCER",
// 	})

// 	//Users
// 	wb.Create(fbClient.Collection("users").NewDoc(), map[string]interface{}{
// 		"number":      "+14069999999",
// 		"welcomeSent": true,
// 	})

// 	//Tracked Sections
// 	wb.Create(fbClient.Collection("sections_tracked").NewDoc(), map[string]interface{}{
// 		"courseName":     "Multidisc Engineering Design",
// 		"courseNumber":   "310R",
// 		"creationTime":   time.Now(),
// 		"crn":            "33761",
// 		"department":     "Engineering",
// 		"departmentAbbr": "EGEN",
// 		"instructor":     "Rutherford, Spencer",
// 		"openSeats":      "0",
// 		"sectionNumber":  37,
// 		"term":           "202230",
// 		"totalSeats":     "10",
// 		"users": []string{
// 			"NTWUmukBvpRuEEXlRyUzKE5R73W2",

// 			"LerHamgx4gPUoqzaH4w26zVsEvu1",

// 			"7VE4YazkwlXpoQlSXHUF9RHmTB33",

// 			"6dBimgScTuYnn6vmCrDDEwyN5Xy1"},
// 	})
// 	_, err := wb.Commit(context.Background())
// 	return err
// }
