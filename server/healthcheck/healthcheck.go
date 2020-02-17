package healthcheck

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// CheckHealth Can be called to check the state of myInfo
func CheckHealth(w http.ResponseWriter, r *http.Request) {
	// Load up a context and http client
	// ctx := r.Context()
	// log.Printf(string(ctx))
	log.Printf("Context loaded. Starting execution.")

	beforeReq := time.Now()
	_, err := http.Get("https://prodmyinfo.montana.edu/pls/bzagent/bzskcrse.PW_SelSchClass")
	totalReqTime := time.Since(beforeReq)
	if err != nil {
		log.Printf("Atlas Appears to be down?")
		w.WriteHeader(500)
		return
	}
	responseBody := fmt.Sprintf("ATLAS Ping: %s", totalReqTime.String())

	log.Printf(responseBody)
	w.Write([]byte(responseBody))
}
