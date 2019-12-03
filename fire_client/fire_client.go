package fireclient

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

var firebaseProjectID = os.Getenv("FIREBASE_PROJECT")
var errNotInitialized = fmt.Errorf("Not Initialized")

type fireclient struct {
	App         *firebase.App
	Auth        *auth.Client
	Firestore   *firestore.Client
	initialized bool
}

type firebaseImpl interface {
	Initialize() error
	GetApp() (*firebase.App, error)
	GetAuth() (*firebase.App, error)
	GetFireStore() (*firebase.App, error)
}

func GetImpl() *firebaseImpl {
	return firebaseImpl{}
}

func (f fireclient) Initialize() error {
	if f.initialized {
		return fmt.Errorf("Already Initialized")
	}

	ctx := context.Background()

	firebaseApp, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return err
	}

	firebaseAuth, err := firebaseApp.Auth(ctx)
	if err != nil {
		return err
	}

	if firebaseProjectID == "" {
		log.WithContext(ctx).Errorf("Firebase Project ID is nil, I cannot continue.")
		panic("Firebase Project ID is nil")
	}

	fbClient, err := firestore.NewClient(ctx, firebaseProjectID)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Could not create new client for Firebase")
		return nil
	}

	f = fireclient{
		App:         firebaseApp,
		Auth:        firebaseAuth,
		Firestore:   fbClient,
		initialized: true,
	}
	return nil
}

func (f fireclient) GetApp() (*firebase.App, error) {
	if !f.initialized {
		return nil, errNotInitialized
	}
	return f.App, nil
}

func (f fireclient) GetAuth() (*auth.Client, error) {
	if !f.initialized {
		return nil, errNotInitialized
	}
	return f.Auth, nil
}

func (f fireclient) GetFireStore() (*firestore.Client, error) {
	if !f.initialized {
		return nil, errNotInitialized
	}
	return f.Firestore, nil
}
