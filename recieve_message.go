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
		responseText = fmt.Sprintf("Available Commands:\nHELP - prints info\nLIST - lists tracked classes and their seats\nMore commands coming soon!")
	}
	if strings.Contains(text, "LIST") {

		//responseText = LookupUserSections(ctx, from, nil)
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
