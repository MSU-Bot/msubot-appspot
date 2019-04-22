package main

import (
	"fmt"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Load up a context and http client
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	log.Infof(ctx, "Context loaded. Starting execution.")

	beforeReq := time.Now()
	_, err := client.Get("https://prodmyinfo.montana.edu/pls/bzagent/bzskcrse.PW_SelSchClass")
	totalReqTime := time.Since(beforeReq)
	if err != nil {
		log.Errorf(ctx, "Atlas Appears to be down?")
		w.WriteHeader(500)
		return
	}
	responseBody := fmt.Sprintf("ATLAS Ping: %s", totalReqTime.String())

	log.Infof(ctx, responseBody)
	w.Write([]byte(responseBody))
}
