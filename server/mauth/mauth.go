package mauth

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	log "github.com/sirupsen/logrus"
)

// VerifyToken verifies tokens, and returns an authenticated token object upon successful verification
func VerifyToken(ctx context.Context, tokenString string) (*auth.Token, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error getting authenticated firebase client")
		return nil, err
	}

	auth, err := app.Auth(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error getting authenticated auth client")
		return nil, err
	}

	token, err := auth.VerifyIDTokenAndCheckRevoked(ctx, tokenString)
	if err != nil {
		log.WithContext(ctx).WithError(err).Warn("Invalid Token")
		return nil, err
	}

	return token, nil
}
