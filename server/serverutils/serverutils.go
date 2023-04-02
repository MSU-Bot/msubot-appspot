package serverutils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/MSU-Bot/msubot-appspot/server/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/plivo/plivo-go"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"
)

const (
	sectionRequestURL    = "https://prodmyinfo.montana.edu/pls/bzagent/bzskcrse.PW_ListSchClassSimple"
	departmentRequestURL = "https://prodmyinfo.montana.edu/pls/bzagent/bzskcrse.PW_SelSchClass"
	sourceNumber         = "14068000110"
)

func GetFirebaseApp(ctx context.Context) (*firebase.App, error) {
	return firebase.NewApp(ctx, &firebase.Config{})
}

// GetFirebaseClient creates and returns a new firebase client, used to interact with the database
func GetFirebaseClient(ctx context.Context) *firestore.Client {
	fbApp, err := GetFirebaseApp(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Could not create new appclient for Firebase")
		return nil
	}

	fbClient, err := fbApp.Firestore(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Could not create new client for Firebase")
		return nil
	}
	log.WithContext(ctx).Infof("successfully opened firestore client")

	return fbClient
}

func MakeAtlasDepartmentRequest(client *http.Client) (*http.Response, error) {
	resp, err := http.Get(departmentRequestURL)
	return resp, err
}

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
	body := fmt.Sprintf("sel_subj=dummy&bl_online=FALSE&sel_day=dummy&sel_online=dummy&term=%v&sel_subj=%v&sel_inst=0&sel_online=&sel_crse=%v&begin_hh=0&begin_mi=0&end_hh=0&end_mi=0",
		term,
		department,
		course)

	return strings.NewReader(body)
}

func ParseDepartmentResponse(response *http.Response) ([]*models.Department, error) {
	departments := []*models.Department{}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	rawDepts := doc.Find("#selsubj").Children()
	rawDepts.Each(func(i int, s *goquery.Selection) {
		abbr, name, found := strings.Cut(s.Text(), "-")
		if !found {
			log.WithField("dept text", s.Text()).Warn("Failed to split department in string")
			return
		}

		departments = append(departments, &models.Department{
			Id:   strings.TrimSpace(abbr),
			Name: strings.TrimSpace(name),
		})
	})

	log.WithField("Processed Departments", departments[0].Name).Info(`Finished processing department list`)
	return departments, nil

}

// ParseSectionResponse turns the http.Response into a slice of sections
func ParseSectionResponse(response *http.Response, crnToFind string) ([]models.Section, error) {
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
// Notification Functions
////////////////////////////

func SendEmail(userdata []*auth.UserRecord, section models.Section) error {
	m := mail.NewV3Mail()
	from := mail.NewEmail("MSUBot", "noreply@unwent.com")
	m.SetFrom(from)
	m.SetTemplateID("d-2be4c913792e48a7ac9860f4216967e3")

	p := mail.NewPersonalization()
	tos := make([]*mail.Email, len(userdata))
	for i, user := range userdata {
		tos[i] = mail.NewEmail(user.DisplayName, user.Email)
	}
	p.AddTos(from)
	p.AddBCCs(tos...)

	p.SetDynamicTemplateData("open_seats", section.AvailableSeats)
	p.SetDynamicTemplateData("course_info", fmt.Sprintf("%s-%s: %s CRN: %s", section.DeptAbbr, section.CourseNumber, section.CourseName, section.Crn))
	p.SetDynamicTemplateData("time", time.Now().Format(time.UnixDate))

	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		return err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return nil
}

// SendText sends a text message to the specified phone number
func SendText(client *http.Client, number, message string) error {

	authID := os.Getenv("PLIVO_AUTH_ID")
	authToken := os.Getenv("PLIVO_AUTH_TOKEN")
	if authID == "" || authToken == "" {
		panic("nil env")
	}

	plivoClient, err := plivo.NewClient(authID, authToken, &plivo.ClientOptions{HttpClient: client})
	if err != nil {
		return err
	}

	response, err := plivoClient.Messages.Create(
		plivo.MessageCreateParams{
			Src:  sourceNumber,
			Dst:  number,
			Text: message,
		},
	)
	if err != nil {
		return err
	}
	log.Infof("Response: %#v\n", response)
	return nil
}

// FetchUserDataWithNumber check firebase to see if the user exists in our database
func FetchUserDataWithNumber(ctx context.Context, number string) (map[string]interface{}, string) {
	checkedNumber := fmt.Sprintf("+%v", strings.Trim(number, " "))

	fbClient := GetFirebaseClient(ctx)
	defer fbClient.Close()

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

func GetUserdata(ctx context.Context, useruids []interface{}) ([]*auth.UserRecord, error) {
	fbApp, err := GetFirebaseApp(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Could not create new appclient for Firebase")
		return nil, err
	}

	authClient, err := fbApp.Auth(ctx)
	if err != nil {
		return nil, err
	}

	firebaseUsers := make([]auth.UserIdentifier, len(useruids))
	for i, uid := range useruids {
		firebaseUsers[i] = auth.UIDIdentifier{UID: uid.(string)}
	}

	result, err := authClient.GetUsers(ctx, firebaseUsers)
	if err != nil {
		return nil, err
	}
	return result.Users, nil
}

// LookupUserNumber looks up a user's phone number from their uid
func LookupUserNumber(ctx context.Context, uid string) (string, error) {
	fbClient := GetFirebaseClient(ctx)
	defer fbClient.Close()

	doc, err := fbClient.Collection("users").Doc(uid).Get(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Tracked user not found. This should've been cleaned up")
		return "", err
	}
	return doc.Data()["number"].(string), nil
}

// MoveTrackedSection moves old sections out of the prod area
func MoveTrackedSection(ctx context.Context, crn, uid, term string) error {
	fbClient := GetFirebaseClient(ctx)
	defer fbClient.Close()

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
