package messenger

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/plivo/plivo-go/xml"

	log "github.com/sirupsen/logrus"
)

// RecieveMessage handles ingest of SMS messages from plivo
func RecieveMessage(ctx echo.Context) {
	r := ctx.Request()

	if err := r.ParseForm(); err != nil {
		ctx.Response().Writer.WriteHeader(http.StatusBadRequest)
		log.Error("Could not parse form")
	}
	from := strings.Join(r.PostForm["From"], "")
	text := strings.Join(r.PostForm["Text"], "")
	text = strings.ToUpper(text)

	responseText := ""
	if strings.Contains(text, "HELP") {
		responseText = "Available Commands:\nHELP - prints this help message\nLIST - lists your tracked classes and their seats"
	} else if strings.Contains(text, "LIST") {
		responseText = "LIST is currently undergoing maintenence, please check the website for class status"
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

	responseBody := xml.ResponseElement{
		Contents: []interface{}{
			new(xml.MessageElement).
				SetType("sms").
				SetDst(from).
				SetSrc("14068000110").
				SetContents(responseText),
		},
	}

	w := ctx.Response().Writer
	w.Header().Add("Content-Type", "text/xml")
	w.Write([]byte(responseBody.String()))
	r.Body.Close()
}
