package mauth

import (
	"context"
	"errors"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// VerifyToken verifies tokens, and returns an authenticated token object upon successful verification
func VerifyToken(ctx echo.Context) (*auth.Token, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.WithContext(context.Background()).WithError(err).Error("Error getting authenticated firebase client")
		return nil, err
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		log.WithContext(context.Background()).WithError(err).Error("Error getting authenticated auth client")
		return nil, err
	}

	tok := ctx.Request().Header.Get("Authorization")
	tok = strings.Replace(tok, "Bearer ", "", -1)

	token, err := auth.VerifyIDTokenAndCheckRevoked(ctx.Request().Context(), tok)
	if err != nil {
		log.WithContext(ctx.Request().Context()).WithError(err).Warn("Invalid Token")
		return nil, err
	}

	return token, nil
}

// VerifyAppengineCron ensures the request is from the appengine cron
func VerifyAppengineCron(ctx echo.Context) error {
	if ctx.Request().Header.Get("X-Appengine-Cron") == "" {
		log.WithContext(ctx.Request().Context()).Warningf("Request is not from the cron. Exiting")
		ctx.Response().Writer.WriteHeader(403)
		return errors.New("request not from Appengine Cron")
	}
	return nil
}
