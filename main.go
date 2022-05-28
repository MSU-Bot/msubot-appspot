package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/SpencerCornish/msubot-appspot/server/api"
	"github.com/SpencerCornish/msubot-appspot/server/dstore"
	"github.com/SpencerCornish/msubot-appspot/server/serverutils"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

const (
	portEnvVariable = "PORT"
	defaultPort     = "8090"
)

func main() {
	port := os.Getenv(portEnvVariable)
	if port == "" {
		port = defaultPort
		log.Infof("Defaulting to port %s", port)
	}

	firebaseClient := serverutils.GetFirebaseClient(context.Background())

	// Get the API Spec (for validation)
	swagger, err := api.GetSwagger()
	if err != nil {
		log.WithError(err).Fatal("Failed to load swagger spec")
	}
	swagger.Servers = nil

	msubotAPI := api.New(dstore.New(*firebaseClient))

	ec := echo.New()
	ec.Use(echoMiddleware.Logger())

	api.RegisterHandlers(ec, msubotAPI)

	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.WithError(err).Fatal("Could not parse port as int")
	}

	ec.Logger.Fatal(ec.Start(fmt.Sprintf("0.0.0.0:%d", portInt)))
}

// func RegisterHandlers(router codegen.EchoRouter, si ServerInterface) {
// 	wrapper := ServerInterfaceWrapper{
// 		Handler: si,
// 	}
// 	router.GET("/pets", wrapper.FindPets)
// 	router.POST("/pets", wrapper.AddPet)
// 	router.DELETE("/pets/:id", wrapper.DeletePet)
// 	router.GET("/pets/:id", wrapper.FindPetById)
// }
