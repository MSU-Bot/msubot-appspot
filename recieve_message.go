package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/plivo/plivo-go/xml"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// ReceiveMessageHandler handles ingest of SMS messages from plivo
func ReceiveMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "Context loaded. Starting execution.")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
	}
	from := strings.Join(r.PostForm["From"], "")
	text := strings.Join(r.PostForm["Text"], "")
	text = strings.ToUpper(text)

	responseText := ""
	if strings.Contains(text, "HELP") {
		responseText = fmt.Sprintf("Available Commands:\nHELP - prints this help message\nLIST - lists your tracked classes and their seats")
	} else if strings.Contains(text, "LIST") {
		fbClient := GetFirebaseClient(ctx)

		_, uid := FetchUserDataWithNumber(ctx, fbClient, from)

		log.Infof(ctx, "Found user with UID: %s", uid)

		if uid != "" {
			trackedDocs, err := fbClient.Collection("sections_tracked").Where("users", "array-contains", uid).Documents(ctx).GetAll()
			if err != nil {
				log.Errorf(ctx, "Could not retrieve tracked docs for the user: %s Err: %v", uid, err)
				responseText = "An error occurred, and we couldn't retrieve the course list. Please try again later."
			} else {
				if len(trackedDocs) == 0 {
					responseText = "You are not currently tracking any courses."
				} else {
					responseText = "Courses tracked:\n"

					for _, doc := range trackedDocs {
						docData := doc.Data()
						responseText = fmt.Sprintf("%s%v%v - %v open seats\n", responseText, docData["departmentAbbr"], docData["courseNumber"], docData["openSeats"])
					}
				}
			}

		} else {
			responseText = "We couldn't find your phone number in the database. Please try again later."
		}
	} else {
		responseText = "Command not found. Reply HELP for available commands."
	}

	responseBody := xml.ResponseElement{
		Contents: []interface{}{
			new(xml.MessageElement).
				SetType("sms").
				SetDst(from).
				SetSrc("14068000110").
				SetContents(responseText),
		},
	}

	w.Header().Add("Content-Type", "text/xml")
	w.Write([]byte(responseBody.String()))
	r.Body.Close()
}
