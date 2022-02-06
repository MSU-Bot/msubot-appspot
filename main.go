package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/SpencerCornish/msubot-appspot/server/apihandler"
	"github.com/SpencerCornish/msubot-appspot/server/dstore"
	"github.com/SpencerCornish/msubot-appspot/server/serverutils"
	log "github.com/sirupsen/logrus"
)

const (
	portEnvVariable = "PORT"
	defaultPort     = "8090"
)

func main() {

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	firebaseClient := serverutils.GetFirebaseClient(ctx)

	// datastore :=

	handler := apihandler.New(dstore.New(*firebaseClient))

	// dataStore :=

	log.Info("Defining http handlers...")
	endpoints.DefineServiceHandlers()
	log.Info("Defining http handlers... Done")

	port := os.Getenv(portEnvVariable)
	if port == "" {
		port = defaultPort
		log.Infof("Defaulting to port %s", port)
	}

	log.Infof("Starting HTTP server on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
