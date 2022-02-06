package apihandler

import (
	"github.com/SpencerCornish/msubot-appspot/server/checksections"
	"github.com/SpencerCornish/msubot-appspot/server/dstore"
	"github.com/SpencerCornish/msubot-appspot/server/gen/api"
	"github.com/SpencerCornish/msubot-appspot/server/mauth"
	"github.com/SpencerCornish/msubot-appspot/server/messenger"
	"github.com/SpencerCornish/msubot-appspot/server/pruner"
	"github.com/labstack/echo/v4"
)

type serverInterface struct {
	datastore dstore.DStore
}

func New(ds dstore.DStore) api.ServerInterface {
	return serverInterface{datastore: ds}
}

// Service API
func (s serverInterface) CheckTrackedSections(ctx echo.Context) error {
	if err := mauth.VerifyAppengineCron(ctx); err != nil {
		return err
	}

	return checksections.HandleRequest(ctx, s.datastore)
}

func (s serverInterface) PruneTrackedSections(ctx echo.Context) error {
	if err := mauth.VerifyAppengineCron(ctx); err != nil {
		return err
	}

	return pruner.HandleRequest(ctx, s.datastore)
}

func (s serverInterface) ReceiveSMS(ctx echo.Context) error {
	messenger.RecieveMessage(ctx)
	return nil
}

// Public API
func (s serverInterface) GetCoursesForDepartment(ctx echo.Context, params api.GetCoursesForDepartmentParams) error {
	panic("implement me")
}

func (s serverInterface) GetDepartments(ctx echo.Context) error {
	panic("implement me")
}

func (s serverInterface) GetMeta(ctx echo.Context) error {
	panic("implement me")
}

func (s serverInterface) GetSections(ctx echo.Context, params api.GetSectionsParams) error {
	panic("implement me")
}

// Authenticated API

func (s serverInterface) GetUserData(ctx echo.Context, userID string) error {
	// 	_, err := mauth.VerifyToken(ctx)

	panic("implement me")
}

func (s serverInterface) UpdateUserData(ctx echo.Context, userID string) error {
	// 	_, err := mauth.VerifyToken(ctx)

	panic("implement me")
}

func (s serverInterface) RemoveTrackedSectionForUser(ctx echo.Context, userID string, sectionID string) error {
	// 	_, err := mauth.VerifyToken(ctx)
	panic("implement me")
}

func (s serverInterface) GetTrackedSectionsForUser(ctx echo.Context, userID string) error {
	// 	_, err := mauth.VerifyToken(ctx)
	panic("implement me")
}

func (s serverInterface) AddTrackedSectionsForUser(ctx echo.Context, userID string) error {
	// 	_, err := mauth.VerifyToken(ctx)
	panic("implement me")
}
