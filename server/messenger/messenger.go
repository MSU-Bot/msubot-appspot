package messenger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/SpencerCornish/msubot-appspot/server/models"
	"github.com/plivo/plivo-go/xml"
)

const (
	notificationMessageTemplate = "%s%s - %s with CRN %s has %s open seats! Get to MyInfo and register before it's gone."
	sourceNumber                = "14068000110"
)

// PlivoRequest is the type sent to Plivo for texts
type plivoRequest struct {
	Src  string `json:"src"`
	Dst  string `json:"dst"`
	Text string `json:"text"`
}

// ReceiveMessage handles ingest of SMS messages from plivo
func ReceiveMessage(from, msgText string) (*xml.ResponseElement, error) {

	responseText := ""
	if strings.Contains(msgText, "HELP") {
		responseText = "Available Commands:\nHELP - prints this help message\nLIST - lists your tracked classes and their seats"
	} else if strings.Contains(msgText, "LIST") {
		responseText = "LIST is currently undergoing maintenance, please check the website for class status"
		// fbClient := serverutils.GetFirebaseClient(ctx)
		// _, uid := serverutils.FetchUserDataWithNumber(ctx, fbClient, from)

		// log.WithContext(ctx).Infof("Found user with UID: %s", uid)

		// if uid != "" {
		// 	trackedDocs, err := fbClient.Collection("sections_tracked").Where("users", "array-contains", uid).Documents(ctx).GetAll()
		// 	if err != nil {
		// 		log.WithContext(ctx).WithError(err).Errorf("Could not retrieve tracked docs for the user: %s", uid)
		// 		responseText = "An error occurred, and we couldn't retrieve the course list. Please try again later."
		// 	} else {
		// 		if len(trackedDocs) == 0 {
		// 			responseText = "You are not currently tracking any courses."
		// 		} else {
		// 			responseText = "Courses tracked:\n"

		// 			for _, doc := range trackedDocs {
		// 				docData := doc.Data()
		// 				responseText = fmt.Sprintf("%s%v%v - %v open seats\n", responseText, docData["departmentAbbr"], docData["courseNumber"], docData["openSeats"])
		// 			}
		// 		}
		// 	}

		// } else {
		// 	responseText = "We couldn't find your phone number in the database. Please try again later."
		// }
	} else {
		responseText = "Command not found. Reply HELP for available commands."
	}

	return &xml.ResponseElement{
		Contents: []interface{}{
			new(xml.MessageElement).
				SetType("sms").
				SetDst(from).
				SetSrc("14068000110").
				SetContents(responseText),
		},
	}, nil

}

func SendSeatMessages(ctx context.Context, client *http.Client, userNumbers []string, section models.Section) error {
	formattedMessage := fmt.Sprintf(notificationMessageTemplate, section.DeptAbbr, section.CourseNumber, section.CourseName, section.Crn, section.AvailableSeats)

	_, err := SendText(client, userNumbers, formattedMessage)

	return err
}

func getPlivoURL(authID string) string {
	return fmt.Sprintf("https://api.plivo.com/v1/Account/%s/Message/", authID)
}

// SendText sends a text message to the specified phone number
func SendText(client *http.Client, numbers []string, message string) (response *http.Response, err error) {
	authID := os.Getenv("PLIVO_AUTH_ID")
	authToken := os.Getenv("PLIVO_AUTH_TOKEN")
	if authID == "" || authToken == "" {
		panic("nil env")
	}

	formattedNumbers := ""
	for _, num := range numbers {
		if formattedNumbers == "" {
			formattedNumbers = num
		} else {
			formattedNumbers = fmt.Sprintf("%s<%s", formattedNumbers, num)
		}
	}

	url := getPlivoURL(authID)
	data := plivoRequest{Src: sourceNumber, Dst: formattedNumbers, Text: message}

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

// // TODO: Move to messager package
// func sendOpenSeatMessages(ctx context.Context, client *http.Client, ds dstore.DStore, users []string, section models.Section) error {
// 	var userNumbers string
// 	message := fmt.Sprintf("%v%v - %v with CRN %v has %v open seats! Get to MyInfo and register before it's gone!", section.DeptAbbr, section.CourseNumber, section.CourseName, section.Crn, section.AvailableSeats)
// 	for _, user := range users {
// 		userRecord, err := ds.GetUser(ctx, user)
// 		if err != nil {
// 			log.WithContext(ctx).Errorf("Unable to get user %s", user)
// 		}

// 		if userNumbers == "" {
// 			userNumbers = userRecord.PhoneNumber
// 		} else {
// 			userNumbers = fmt.Sprintf("%v<%v", userNumbers, userRecord.PhoneNumber)
// 		}
// 	}
// 	resp, err := serverutils.SendText(client, userNumbers, message)
// 	if err != nil {
// 		log.WithContext(ctx).WithError(err).Error("error sending texts")
// 		return err
// 	}
// 	log.WithContext(ctx).WithField("response", resp).Info("Texts sent")
// 	return nil
// }
